FROM golang:1.26.4-alpine

WORKDIR /app

COPY go.mod ./
COPY main.go ./
COPY frontend ./frontend

RUN go build -o url-shortener .

ENV PORT=8080
EXPOSE 8080

CMD ["./url-shortener"]
