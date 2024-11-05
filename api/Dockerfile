FROM golang:1.23.2-alpine

COPY . .

WORKDIR /api

RUN go build -o api .

EXPOSE 8080

CMD ["./api"]