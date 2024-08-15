FROM golang:1.22.2-alpine AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . /app
RUN CGO_ENABLED=0 go build -o calculator cmd/calculator/main.go

FROM alpine:latest
COPY --from=build /app/calculator /usr/bin/calculator
CMD ["calculator"]
