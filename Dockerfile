
FROM golang:1.19-alpine
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod tidy
RUN go build -o app main.go
EXPOSE 3000:3000
EXPOSE 25:25
CMD ["./app"]