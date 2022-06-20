package main

import (
	"consumer/models"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"os"
	"github.com/rabbitmq/amqp091-go"
)

var rabbitHost = os.Getenv("RABBIT_HOST")
var rabbitPort = os.Getenv("RABBIT_PORT")
var rabbitUsername = os.Getenv("RABBIT_USER")
var rabbitPassword = os.Getenv("RABBIT_PASS")

var psqlUser = os.Getenv("PSQL_USER")
var psqlPass = os.Getenv("PSQL_PASS")
var psqlPort = os.Getenv("PSQL_PORT")
var psqlName = os.Getenv("PSQL_NAME")

func main() {
	fmt.Println("Consumer: I'm sentient.")

	db := models.DBInit(rabbitHost, psqlUser, psqlPass, psqlName, psqlPort)

	connLink := fmt.Sprintf("amqp://%v:%v@%v:%v/", rabbitUsername, rabbitPassword, rabbitHost, rabbitPort)

	conn, err := amqp091.Dial(connLink)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to RMQ", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("%s: %s", "Failed to create channel", err)
	}

	queue, err := channel.QueueDeclare(
		"publisher", //name
		false,       //durable
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to connect to queue", err)
	}

	defer conn.Close()
	defer channel.Close()

	messages, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("%s: %s", "Failed to register consumer", err)
	}

	forever := make(chan bool)

	go func() {
		for message := range messages {
			log.Printf("Received a message: %s", message.Body)

			msg := parseMessage(message.Body)
			log.Printf("Parsed message: %s", msg)
			switch msg["tag"] {
			case "user":
				log.Println("Identified user tag.")
				user := models.User{
					Name:  msg["name"],
					Email: msg["email"],
				}
				log.Printf("Creating user: %v", user)
				db.Create(&user)
			case "msg":
				uid, err := strconv.Atoi(msg["uid"])
				if err != nil {
					log.Fatal("Failed to parse user ID from paylaod.")
				}
				db.Create(&models.Message{Content: msg["content"], UserID: uint(uid)})
			}

			message.Ack(false)
		}
	}()

	fmt.Println("Consumer is running...")
	<-forever
}

func parseMessage(msg []byte) map[string]string {
	var message map[string]string
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Fatalf("%s: %s", "Failed to parse JSON message.", err)
	}

	return message
}