FROM golang:1.22.5-bullseye
# FROM golang:1.22.5-alpine

WORKDIR /ginapp
COPY . .

RUN go build -o main main.go

CMD ["./main"]
