FROM golang:1.23.2

RUN mkdir /app
ADD . /app

WORKDIR /app

COPY go.mod go.sum app.env ./

RUN go build -o app ./cmd/app/main.go

EXPOSE 8080

CMD [ "/app/app" ]