FROM golang:latest

WORKDIR /app

COPY . .

EXPOSE 1234

CMD ["go", "run", "app/main.go"]
