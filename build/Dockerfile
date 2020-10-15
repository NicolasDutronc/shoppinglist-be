FROM golang:1.14-alpine as builder

WORKDIR /app/src

COPY . .

ENV CGO_ENABLED 0

RUN go build -o app ./cmd/main.go && chmod +x app

FROM scratch

COPY --from=builder /app/src/app /app
COPY --from=builder /app/src/config.yml /config.yml
COPY --from=builder /app/src/public/index.html /public/index.html

CMD [ "/app" ]