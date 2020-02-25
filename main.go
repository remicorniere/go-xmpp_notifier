package main

import (
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
)

const (
	server_ip = iota
	serverDomain
	correspondent
	login
)

func main() {
	config := xmpp.Config{
		TransportConfiguration: xmpp.TransportConfiguration{
			Address: os.Args[serverDomain] + "5222",
		},
		Jid:          os.Args[login],
		Credential:   xmpp.Password("test"),
		StreamLogger: os.Stdout,
		Insecure:     true,
		// TLSConfig: tls.Config{InsecureSkipVerify: true},
	}

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
