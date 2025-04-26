FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GDOS=linux go build -o /api-gateway ./cmd/api-gateway

FROM alpine:3.20
WORKDIR /root/
COPY --from=builder /api-gateway api-gateway
EXPOSE 8181
ENTRYPOINT ["./api-gateway"]

#FROM golang:1.23-alpine
 #
 #WORKDIR /app
 #
 #ADD go.mod .
 #ADD go.sum .
 #RUN go mod download
 #
 #COPY . .
 #
 #RUN CGO_ENABLED=0 GDOS=linux go build -o main .
 #EXPOSE 8080
 #CMD ["./main"]
