FROM golang:1.25.7-alpine3.23@sha256:f6751d823c26342f9506c03797d2527668d095b0a15f1862cddb4d927a7a4ced AS builder

WORKDIR /app

RUN apk --no-cache add ca-certificates \
    && addgroup -S app \
    && adduser --uid 19998 --shell /bin/false -S app -G app \
    && cat /etc/passwd | grep app > /etc/passwd_app \
    && cat /etc/group | grep app > /etc/group_app \
    && mkdir -p /home/app /app_runtime \
    && chown -R app:app /home/app /app_runtime

ENV CGO_ENABLED=0
ENV GOOS=linux

ARG SERVICE_NAME

COPY go.mod go.sum ./
COPY . .

RUN if [ -d "./vendor" ]; then \
        go build -mod=vendor -o /app_runtime/app -ldflags="-s -w" -trimpath -v ./cmd/$SERVICE_NAME; \
    else \
        go mod download; \
        go build -o /app_runtime/app -ldflags="-s -w" -trimpath -v ./cmd/$SERVICE_NAME; \
    fi

FROM scratch

COPY --from=builder /etc/passwd_app /etc/passwd
COPY --from=builder /etc/group_app  /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder --chown=app:app /home/app /home/app
COPY --from=builder --chown=app:app /app_runtime/ /app/

WORKDIR /app

ENV HOME=/home/app

USER app

EXPOSE 7000

ENTRYPOINT ["./app"]
