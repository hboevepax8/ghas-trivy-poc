package main

import (
    "context"
    "errors"
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "runtime"
    "syscall"
    "time"

    "github.com/ardanlabs/conf/v3"
    "go.uber.org/automaxprocs/maxprocs"
    "go.uber.org/zap"

    "github.com/hboevepax8/ghas-trivy-poc/cmd/api/handlers"
    "github.com/hboevepax8/ghas-trivy-poc/internal/sys/logger"
)

var build = "development"

func main() {
    log, err := logger.New("ghas-trivy-poc")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer log.Sync()

    if err := run(log); err != nil {
        log.Errorw("startup", "ERROR", err)
        log.Sync()
        os.Exit(1)
    }
}

func run(log *zap.SugaredLogger) error {

    // =========================================================================
    // GOMAXPROCS

    opt := maxprocs.Logger(log.Infof)
    if _, err := maxprocs.Set(opt); err != nil {
        return fmt.Errorf("maxprocs: %w", err)
    }
    log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

    // =========================================================================
    // Configuration

    cfg := struct {
        conf.Version
        Web struct {
            ReadTimeout     time.Duration `conf:"default:5s"`
            WriteTimeout    time.Duration `conf:"default:10s"`
            IdleTimeout     time.Duration `conf:"default:120s"`
            ShutdownTimeout time.Duration `conf:"default:20s"`
            APIHost         string        `conf:"default:0.0.0.0:4000"`
        }
    }{
        Version: conf.Version{
            Build: build,
            Desc:  "copyright ...",
        },
    }

    const prefix = "SECTEST"
    help, err := conf.Parse(prefix, &cfg)
    if err != nil {
        if errors.Is(err, conf.ErrHelpWanted) {
            fmt.Println(help)
            return nil
        }
        return fmt.Errorf("parsing config: %w", err)
    }

    // =========================================================================
    // Start API

    log.Infow("startup", "status", "initializing V1 API support")

    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

    mux := handlers.ApiMux(handlers.APIMuxConfig{Log: log})

    api := http.Server{
        Addr:         cfg.Web.APIHost,
        Handler:      mux,
        ReadTimeout:  cfg.Web.ReadTimeout,
        WriteTimeout: cfg.Web.WriteTimeout,
        IdleTimeout:  cfg.Web.IdleTimeout,
        ErrorLog:     zap.NewStdLog(log.Desugar()),
    }

    serverErrors := make(chan error, 1)

    go func() {
        log.Infow("startup", "status", "api router started", "host", api.Addr)
        serverErrors <- api.ListenAndServe()
    }()

    // =========================================================================
    // Shutdown

    select {
    case err := <-serverErrors:
        return fmt.Errorf("server error: %w", err)

    case sig := <-shutdown:
        log.Infow("shutdown", "status", "shutdown started", "signal", sig)
        defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

        ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
        defer cancel()

        if err := api.Shutdown(ctx); err != nil {
            api.Close()
            return fmt.Errorf("could not stop server gracefully: %w", err)
        }
    }

    return nil
}
