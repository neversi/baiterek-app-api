
FROM golang
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o app main.go
EXPOSE 3000:3000
CMD ["./app"]