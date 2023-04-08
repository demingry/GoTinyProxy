FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o proxy-server .

EXPOSE 18080

CMD ["./proxy-server"]
