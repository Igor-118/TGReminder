FROM golang:1.26 AS bobrkurwa

WORKDIR /app

COPY . .

RUN go mod download && go build -o app ./cmd/bot

FROM golang:1.26

WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

COPY --from=bobrkurwa /app/app .
COPY --from=bobrkurwa /app/entrypoint.sh .
COPY --from=bobrkurwa /app/migrations /app/migrations

# RUN apk --no-cache add ca-certificates
CMD ["./entrypoint.sh"]