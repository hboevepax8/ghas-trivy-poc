FROM golang:1.20 as build_api
ENV CGO_ENABLED 0
COPY . /app
WORKDIR /app/cmd/s3
RUN go build -ldflags "-X main.build=production"

FROM alpine:3.10
ARG BUILD_DATE
ARG BUILD_REF
RUN addgroup -g 1000 -S api && \
    adduser -u 1000 -h /app -G api -S api
COPY --from=build_api --chown=api:api /app/cmd/s3 /app/s3
WORKDIR /app
USER api
CMD ["./s3"]
