FROM golang:1.8
WORKDIR /go/src/app
COPY . .
RUN go-wrapper download
CMD ["go", "run", "tunnel.go"]