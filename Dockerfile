FROM golang:1.18.0-alpine3.15 as builder

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o service ./cmd/server/main.go


# generate clean, final image for end users
FROM gcr.io/distroless/base-debian10

COPY --from=builder /build/service /app/
COPY --from=builder /build/configs /app/configs/

LABEL org.opencontainers.image.source="https://github.com/dobriychelpozitivniy/go-coding-challenge-jr"

# executable
ENTRYPOINT [ "/app/service", "--config", "/app/configs/prod" ]
