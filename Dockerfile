FROM golang:1.14.4-buster as builder
WORKDIR /go/src/sshls
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o sshls .

FROM alpine:2.7
WORKDIR /root
COPY --from=builder /go/src/sshls/sshls .

ENTRYPOINT [ "./sshls" ]
