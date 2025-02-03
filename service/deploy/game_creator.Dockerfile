FROM golang:1.23.5 AS builder

WORKDIR /build

# build app
COPY cmd ./cmd
COPY internal ./internal
COPY go.* ./

RUN CGO_ENABLED=0 go build -o /build/app ./cmd/game_creator/main.go

FROM alpine:3.20.3

RUN adduser --disabled-password -u 1000 user

USER user

WORKDIR /app

COPY --chown=user:user --from=builder /build/app app

CMD [ "/app/app" ]