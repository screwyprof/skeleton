FROM golang:alpine as builder

RUN apk add --update --no-cache curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.9.1/migrate.linux-amd64.tar.gz | tar xvz && mv migrate.linux-amd64 /go/bin/migrate

WORKDIR /migrate
COPY . ./

ENTRYPOINT ["/migrate/migrate.sh"]
CMD ["--help"]