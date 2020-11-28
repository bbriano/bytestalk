FROM golang:1.15

WORKDIR /app

COPY . .

RUN go build server.go

ENV PORT=8080

EXPOSE 8080

CMD ["./server"]
