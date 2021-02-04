FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/eviofragoso/go_shortener
COPY main.go  .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go_shortener .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/eviofragoso/go_shortener/go_shortener .
RUN echo '{}' > database.json
CMD ["./go_shortener"]