FROM golang:1.24.4-alpine

WORKDIR /usr/src/app

EXPOSE 8080

ENV port=8080

COPY . .

RUN go build -o server

CMD [ "./server" ]