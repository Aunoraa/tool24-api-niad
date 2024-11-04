FROM golang:latest

WORKDIR /app

# Copy source code vào container
COPY . .

# Build ứng dụng
RUN go build -o api.exe .

# Expose port mà ứng dụng lắng nghe
EXPOSE 8080

# Chạy ứng dụng bằng file api.exe
CMD ["./api.exe"]