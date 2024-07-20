FROM golang:1-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o llmt github.com/blwsh/llmt/cmd/analyze

FROM golang:1-alpine
COPY --from=builder /app/llmt /llmt
CMD ["/llmt"]
