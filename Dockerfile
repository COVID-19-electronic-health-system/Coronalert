FROM golang:1.13.5
LABEL Carter Klein <carterklein13@gmail.com>
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]