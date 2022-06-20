FROM golang:1-alpine AS build

RUN apk add --no-cache git

WORKDIR /src

RUN go mod init consumer
RUN go get gorm.io/gorm
RUN go get gorm.io/driver/postgres
RUN go get github.com/rabbitmq/amqp091-go

COPY . /src

RUN go build .

FROM alpine as runtime

COPY --from=build /src/consumer /app/consumer

CMD [ "/app/consumer" ]