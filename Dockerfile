FROM golang:1.23.3 AS go

WORKDIR /app

COPY . ./

RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o ./app .

FROM alpine:latest
WORKDIR /app

COPY --from=go /app/app /app

# EXPOSE 8080

CMD ["./app"]
