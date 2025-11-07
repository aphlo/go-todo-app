# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS base
WORKDIR /app

COPY go.mod ./
RUN go mod download
COPY . .

FROM base AS dev
RUN go install github.com/air-verse/air@v1.52.3
CMD ["air", "-c", ".air.toml"]

FROM base AS build
RUN CGO_ENABLED=0 GOOS=linux go build -o /todo ./cmd/server

FROM gcr.io/distroless/base-debian12 AS runtime
WORKDIR /app
COPY --from=build /todo /app/todo
ENV PORT=8080
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/todo"]
