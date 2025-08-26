FROM golang:1.21-alpine AS builder
WORKDIR /app
# Cài git (bắt buộc cho go mod download với module từ git)
RUN apk add --no-cache git
# copy go.mod trước để tận dụng cache
COPY go.mod go.sum ./
RUN go mod download

# copy toàn bộ source
COPY . .

# dọn dẹp dependencies để tránh lỗi
RUN go mod tidy

# build binary
RUN go build -o /bin/health-memory ./cmd

# final image
FROM gcr.io/distroless/base-debian12
COPY --from=builder /bin/health-memory /bin/health-memory
EXPOSE 8080
CMD ["/bin/health-memory"]
