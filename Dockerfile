FROM golang:alpine

WORKDIR /app

COPY . .

RUN go build -o proxy-server .

EXPOSE 18080

RUN wget https://bin.equinox.io/c/4VmDzA7iaHb/ngrok-stable-linux-amd64.zip -O ngrok.zip \
    && unzip ngrok.zip \
    && chmod +x ngrok \
    && rm ngrok.zip

CMD ["sh", "-c", "./proxy-server & ./ngrok tcp 18080 --log stdout --log-level debug --region us --authtoken $NGROK_AUTH_TOKEN"]