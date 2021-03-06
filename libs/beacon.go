package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net"
	"net/http"
	"strconv"
	"time"
)

type Message struct {
	Addresses []string `json:"addresses"`
	Port      int      `json:"port"`
}

func createBeaconMessage() (string, error) {
	var err error
	Print("Creating message for beacon")
	config, err := GetConfig(true)
	message := Message{Port: config.Server.Port}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, a := range addresses {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Print("Found address %s", ipnet.IP.String())
				message.Addresses = append(message.Addresses, ipnet.IP.String())
			}
		}
	}
	data, err := json.Marshal(message)
	if err != nil {
		Print("Error converting message to JSon")
		return "", err
	}
	return string(data), nil
}

func publishMessage(message string) error {
	Print("Registering server at beacon")
	beacon := fmt.Sprintf("%s://%s:%s@%s/%s", config.Beacon.Protocol, config.Beacon.UserName,
		config.Beacon.Password, config.Beacon.Address, config.Beacon.Queue)
	Print("Dialing to beacon at %s", beacon)
	conn, err := amqp.Dial(beacon)
	if err != nil {
		Print("Error connecting to beacon at %s\n%v", beacon, err)
		return err
	}
	defer conn.Close()

	Print("Openning a channel to beacon")
	ch, err := conn.Channel()
	if err != nil {
		Print("Error creating a new channel to the beacon\n%v", err)
		return err
	}
	defer ch.Close()

	Print("Declaring an exchange to beacon")
	err = ch.ExchangeDeclare(config.Beacon.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		Print("Error declaring exchange beacon\n%v", err)
		return err
	}

	Print("Publishing message")
	err = ch.Publish(config.Beacon.Exchange, config.Beacon.RoutingKey, false, false,
		amqp.Publishing{ContentType: "text/plain", Body: []byte(message), Expiration: strconv.Itoa(config.Beacon.Expiration * 1000)})
	if err != nil {
		Print("Error publishing message to beacon\n%v", err)
		return err
	}
	Print("Message sent")

	return nil
}

func consumeBeaconMessage(debug bool) (*Message, error) {
	var serverInfo *Message

	PrintDebug(debug, "Looking servers at beacon")
	beacon := fmt.Sprintf("%s://%s:%s@%s/%s", config.Beacon.Protocol, config.Beacon.UserName,
		config.Beacon.Password, config.Beacon.Address, config.Beacon.Queue)
	PrintDebug(debug, "Dialing to beacon at %s", beacon)
	conn, err := amqp.Dial(beacon)
	if err != nil {
		PrintDebug(debug, "Error connecting to beacon at %s\n%v", beacon, err)
		return nil, err
	}
	defer conn.Close()

	PrintDebug(debug, "Openning a channel to beacon")
	ch, err := conn.Channel()
	if err != nil {
		PrintDebug(debug, "Error creating a new channel to the beacon\n%v", err)
		return nil, err
	}
	defer ch.Close()

	PrintDebug(debug, "Declaring an exchange to beacon")
	err = ch.ExchangeDeclare(config.Beacon.Exchange, "direct", true, false, false, false, nil)
	if err != nil {
		PrintDebug(debug, "Error declaring exchange beacon\n%v", err)
		return nil, err
	}

	PrintDebug(debug, "Declaring a queue to beacon")
	q, err := ch.QueueDeclare(config.Beacon.Queue, true, false, false, false, nil)
	if err != nil {
		PrintDebug(debug, "Error declaring queue to beacon\n%v", err)
		return nil, err
	}

	PrintDebug(debug, "Binding to queue")
	err = ch.QueueBind(q.Name, config.Beacon.RoutingKey, config.Beacon.Exchange, false, nil)
	if err != nil {
		PrintDebug(debug, "Error binding to queue\n%v", err)
		return nil, err
	}

	messages, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		PrintDebug(debug, "Error consuming message\n%v", err)
		return nil, err
	}

	data := <-messages
	err = json.Unmarshal(data.Body, &serverInfo)
	if err != nil {
		PrintDebug(debug, "Error unmarshalling beacon info\n%v", err)
		return nil, err
	}

	PrintDebug(debug, "Server info - Addresses: %v, Port %d", serverInfo.Addresses, serverInfo.Port)
	return serverInfo, nil
}

func RegisterServer() error {
	message, err := createBeaconMessage()
	Print("Message for beacon:\n %s", message)
	if err != nil {
		Print("Error creating beacon message: %v", err)
		return err
	}
	for {
		publishMessage(message)
		time.Sleep(time.Duration(config.Beacon.Expiration) * time.Second)
	}
	return nil
}

func checkServerConnection(debug bool, address string, port int) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	url := fmt.Sprintf("http://%s:%d", address, port)
	PrintDebug(debug, "Trying connection to %s", url)
	body, err := client.Get(url + "/static/")
	if err != nil {
		return "", err
	}
	PrintDebug(debug, "Got status code %d", body.StatusCode)
	if body.StatusCode >= 300 {
		return "", errors.New("Connection failed")
	}
	return url, nil
}

func GetServerAddress(debug bool) (string, error) {
	var uri string
	server, err := consumeBeaconMessage(debug)
	if err != nil {
		PrintDebug(debug, "Error reading message from beacon: %v", err)
		return "", err
	}

	for _, address := range server.Addresses {
		uri, err = checkServerConnection(debug, address, server.Port)
		if err == nil {
			return uri, nil
		}
	}
	return "", errors.New("Server not found")

}
