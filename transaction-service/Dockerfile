FROM golang:1.23-alpine

WORKDIR /app

ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GDOS=linux go build -o main .
EXPOSE 8088
CMD ["./main"]
