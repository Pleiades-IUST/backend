FROM golang:1.24.1-alpine

RUN apk add --no-cache postgresql-client

RUN go env -w GOPROXY=https://goproxy.io,direct

RUN mkdir -p /home/app

WORKDIR /home/app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o Pleiades ./cmd/*

RUN chmod +x scripts/create_database.sh

EXPOSE 8080

RUN echo '#!/bin/sh' > /home/app/entrypoint.sh && \
    echo 'sh /home/app/scripts/create_database.sh' >> /home/app/entrypoint.sh && \
    echo 'exec /home/app/Pleiades "$@"' >> /home/app/entrypoint.sh

RUN chmod +x entrypoint.sh

ENTRYPOINT ["sh", "/home/app/entrypoint.sh"]