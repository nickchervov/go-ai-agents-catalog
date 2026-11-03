FROM golang:1.26.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /ai-catalog ./cmd/api/main.go

FROM alpine:3.20

COPY --from=builder /ai-catalog /ai-catalog

USER 65532:65532

ENTRYPOINT [ "/ai-catalog" ]