FROM golang:1.22

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN go build -o main main.go

EXPOSE 3030
CMD ["./main"]
