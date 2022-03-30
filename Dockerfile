FROM golang:1.18.0-alpine3.15 as builder

ARG BITLY_OAUTH_TOKEN=${BITLY_OAUTH_TOKEN:-""}
ENV BITLY_OAUTH_TOKEN=${BITLY_OAUTH_TOKEN}


WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o service ./cmd/server/main.go


# generate clean, final image for end users
FROM gcr.io/distroless/base-debian10

COPY --from=builder /build/service /app/
COPY --from=builder /build/configs /app/configs/


# executable
ENTRYPOINT [ "/app/service", "--config", "/app/configs/prod" ]
