FROM golang:1.24.1-alpine

RUN apk add --no-cache postgresql-client

RUN go env -w GOPROXY=https://goproxy.io,direct

RUN mkdir -p /home/app

WORKDIR /home/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

COPY ./scripts/create_database.sh /home/app/scripts/create_database.sh

RUN go build -o Pleiades ./cmd/*

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN chmod +x /home/app/scripts/create_database.sh

EXPOSE 8080

RUN echo '#!/bin/sh' > /home/app/entrypoint.sh && \
    echo 'ls scripts' > /home/app/entrypoint.sh && \
    echo '/home/app/scripts/create_database.sh' >> /home/app/entrypoint.sh && \
    echo 'exec /home/app/Pleiades "$@"' >> /home/app/entrypoint.sh

ENTRYPOINT ["sh", "/home/app/entrypoint.sh"]