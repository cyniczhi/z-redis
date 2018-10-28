FROM golang:latest
ENV SRC_DIR /go/src/github.com/cyniczhi/z-redis/
ENV PATH $PATH:$SRC_DIR
RUN mkdir -p $SRC_DIR
ADD . $SRC_DIR
WORKDIR $SRC_DIR
RUN go build -o server.out main.go
EXPOSE 9999
CMD ["/bin/sh", "-c", "./server.out"]