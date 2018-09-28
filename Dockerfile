FROM golang:1.10

ADD . /go/src/github.com/khisakuni/schubert-services/
WORKDIR /go/src/github.com/khisakuni/schubert-services

RUN go get -u github.com/golang/dep/...
RUN dep ensure

RUN go build -o ./cmd/migrate/migrate ./cmd/migrate
# RUN ["./cmd/migrate/migrate", "up"]
 
RUN go build -o ./user/user ./user
CMD ./user/user
