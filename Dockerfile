FROM golang:1.5

ADD . /go/src/write-better

RUN go get github.com/dansackett/go-text-processors
RUN go install write-better
RUN echo "export GOPATH=/var/www" >> /etc/profile

ENTRYPOINT /go/bin/write-better
EXPOSE 8080
