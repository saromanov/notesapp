FROM golang:latest
RUN mkdir /app
ENV GOPATH /app
ADD . /app/
WORKDIR /app
RUN go get github.com/tools/godep
RUN $GOPATH/bin/godep go build -o main .
CMD ["/app/main", "config.hcl"]