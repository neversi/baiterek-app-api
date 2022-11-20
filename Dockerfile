
FROM golang
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod init
RUN go build
EXPOSE 3000:3000
CMD ["./app"]