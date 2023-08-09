package handlers

import (
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
	"go.uber.org/zap"
)

type APIMuxConfig struct {
	Log *zap.SugaredLogger
}

func ApiMux(cfg APIMuxConfig) http.Handler {
	m := httptreemux.NewContextMux()

	return m
}
