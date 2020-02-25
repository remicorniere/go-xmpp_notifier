package main

import (
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
)

const (
	defaultServerPort = "5222"
	serverDomain      = iota
	correspondent
	login
	pass
	serverPort
)

func main() {
	var port string
	if os.Args[serverPort] == "" {
		port = defaultServerPort
	} else {
		port = os.Args[serverPort]
	}
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: os.Args[serverDomain] + ":" + port,
		},
		Jid:          os.Args[login],
		Credential:   xmpp.Password(os.Args[pass]),
		StreamLogger: os.Stdout,
		Insecure:     false,
	}

	fmt.Println("connecting to : " + os.Args[serverDomain][1:] + ":" + port)
	fmt.Println("sending messages to user :" + os.Args[correspondent][1:])
	router := xmpp.NewRouter()

	client, err := xmpp.NewClient(config, router, errorHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	m := stanza.Message{Attrs: stanza.Attrs{To: os.Args[correspondent], Type: stanza.MessageTypeChat}, Body: "JUST TESTING LUL"}
	err = client.Send(m)
	if err != nil {
		panic(err)
	}

}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
