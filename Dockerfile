FROM golang

LABEL Carter Klein <carterklein13@gmail.com>

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build

EXPOSE 8080:8080

CMD ["./app"]