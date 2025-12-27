FROM golang:1.25-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags '-s -w' -o /goapp .

FROM gcr.io/distroless/static AS runner

COPY --from=builder /goapp /app

EXPOSE 8080

ENTRYPOINT ["/app"]

