#The build statge
FROM golang:1.22 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/*.go

#The run statge
FROM scratch

WORKDIR /app

#Copy CA certificates

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/api . 

EXPOSE 8080

CMD ["./api"]