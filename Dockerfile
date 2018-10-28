FROM golang:latest
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o server.out main.go
PORT 9999
CMD ["/app/server.out"]