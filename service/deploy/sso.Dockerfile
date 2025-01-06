FROM golang:1.22.2 AS builder

WORKDIR /sso

# build app
COPY cmd ./cmd
COPY internal ./internal
COPY go.* ./

RUN CGO_ENABLED=0 go build -o /build/app ./cmd/sso/main.go

FROM alpine:3.20.3

RUN adduser --disabled-password -u 1000 user

USER user

WORKDIR /app

COPY --chown=user:user --from=builder /build/app app
COPY --chown=user:user config/docker_config.yml config.yml

CMD [ "/app/app", "--config", "/app/config.yml" ]