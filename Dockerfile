FROM golang:1.15 AS base

RUN mkdir /app
ADD . /go/src/github.com/kevineaton/simple-auth
WORKDIR /go/src/github.com/kevineaton/simple-auth

RUN go build -mod=vendor .

FROM busybox:glibc
WORKDIR /go/src/github.com/kevineaton/simple-auth
COPY --from=base /etc/ssl/certs /etc/ssl/certs
COPY --from=base /go/src/github.com/kevineaton/simple-auth/simple-auth .

CMD ./simple-auth