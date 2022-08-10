# STEP 1 Build executable binary
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

ENV USER=appuser
ENV UID=1000

# https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

WORKDIR /app
COPY . .
RUN GOOS=linux GOARCH=amd64 go build -v -ldflags="-w -s" -o /app/backend ./cmd/api/main.go

# STEP 2 Build small image
FROM alpine
WORKDIR /app

COPY --from=alpine /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

COPY --from=builder /app/backend /app/backend

USER appuser:appuser

EXPOSE 5000

ENTRYPOINT ["/app/backend"]