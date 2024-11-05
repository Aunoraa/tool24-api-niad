FROM golang:1.23-alpine

WORKDIR /api

COPY . .

RUN go build -o api .

EXPOSE 8080

CMD ["./api"]