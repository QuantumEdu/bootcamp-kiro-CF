FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOARCH=arm64 GOOS=linux go build -o /handler cmd/lambda/main.go

FROM public.ecr.aws/lambda/provided:al2-arm64
COPY --from=builder /handler /var/runtime/bootstrap
COPY templates/ /var/task/templates/
WORKDIR /var/task
CMD ["bootstrap"]
