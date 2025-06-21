FROM golang:1.24.4-alpine

WORKDIR /usr/src/app

COPY . .

RUN go build -o server

CMD [ "./server" ]