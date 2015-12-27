FROM golang:latest
RUN mkdir /app
ENV GOPATH /app
ADD . /app/
WORKDIR /app
RUN go get github.com/saromanov/notesapp/configs
RUN go get github.com/saromanov/notesapp/db
RUN go get github.com/saromanov/notesapp/service
RUN go get github.com/saromanov/notesapp/logging
EXPOSE 8081
RUN go build -o main .
CMD ["/app/main", "config.hcl"]
