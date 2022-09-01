FROM golang:1.19.0-bullseye

COPY . /app
WORKDIR /app

RUN go build .

CMD [ "./booty-mover" ]