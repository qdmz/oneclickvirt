# Backend Dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN apk add --no-cache build-base git && go build -o oneclickvirt .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/oneclickvirt /app/
COPY config/config.yaml /app/config.yaml
EXPOSE 8080
ENTRYPOINT ["/app/oneclickvirt"]
