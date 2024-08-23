FROM golang:1.22.5-alpine AS builder

ENV GO111MODULE=on

WORKDIR /app

# Copia el código fuente al contenedor
COPY . /app/repo_remote

WORKDIR /app/repo_remote

# Descarga las dependencias y compila la aplicación
RUN go mod tidy
RUN GOOS=linux go build -o /app/myapp

FROM alpine:latest

RUN apk --no-cache add tzdata
ENV TZ=America/Santiago

RUN addgroup -S app -g 1000 && adduser -S -g app app --uid 1000

COPY --from=builder --chown=app:app /app/myapp /app/myapp
RUN apk add --no-cache git postgresql-client
COPY /script/wait-for-postgres.sh /app/wait-for-postgres.sh
RUN chmod +x /app/wait-for-postgres.sh

USER app
WORKDIR /app

EXPOSE 8080
CMD ["/bin/sh", "/app/wait-for-postgres.sh", "/app/myapp"]
