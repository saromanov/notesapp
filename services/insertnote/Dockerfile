FROM golang:latest
RUN mkdir /app
ENV GOPATH /app
ADD . /app/
WORKDIR /app
RUN go get ./...
EXPOSE 8081
CMD ["/app/main", "config.hcl"]