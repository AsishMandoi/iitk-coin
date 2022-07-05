FROM golang:1.16.6-alpine3.14 as build
RUN apk add build-base
WORKDIR /app/src
COPY . .
RUN go mod download
RUN go build -o iitk-coin-server

###
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/bin
COPY --from=build /app/src/iitk-coin-server .
COPY --from=build /app/src/.env .

ENV PORT=80
EXPOSE 80

CMD ["./iitk-coin-server"]
