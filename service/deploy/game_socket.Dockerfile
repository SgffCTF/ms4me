FROM golang:1.24.3 AS builder

WORKDIR /build

# build app
COPY go.* ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /build/app ./cmd/game_socket/main.go

FROM alpine:3.20.3

RUN apk update && apk add --no-cache curl

RUN adduser --disabled-password -u 1000 user
USER user

WORKDIR /app
COPY --chown=user:user --from=builder /build/app app

CMD [ "/app/app" ]