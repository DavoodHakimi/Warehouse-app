FROM golang:1.26-alpine AS builder

WORKDIR /app

# RUN go env -w GOPROXY=https://mirror-go.runflare.com,direct
COPY go.mod go.sum ./
RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o Warehouse-backend ./cmd/api

FROM gcr.io/distroless/static-debian12

WORKDIR /app
COPY --from=builder /app/Warehouse-backend /Warehouse-backend
COPY --from=builder /app/.env.docker /.env

CMD ["/Warehouse-backend"]