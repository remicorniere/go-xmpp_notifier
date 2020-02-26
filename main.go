package main

import (
	"fmt"
	"gosrc.io/xmpp"
	"gosrc.io/xmpp/stanza"
	"log"
	"os"
	"strconv"
)

const (
	defaultServerPort = "5222"
	serverDomain      = iota
	correspondent
	login
	pass
	serverPort
	message
	correspondent_is_room
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

	router := xmpp.NewRouter()

	client, err := xmpp.NewClient(config, router, errorHandler)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	err = client.Connect()
	if err != nil {
		panic(err)
	}

	// Send presence to connect to chat room, if specified, and set the correspondentJid
	var correspondentJid string
	isCorrespRoom, err := strconv.ParseBool(os.Args[correspondent_is_room])
	if err == nil && isCorrespRoom {
		correspondentJid = os.Args[correspondent] + "@" + os.Args[serverDomain]
		// Sending room presence
		p := stanza.NewPresence(
			stanza.Attrs{To: correspondentJid},
		)
		err = client.Send(p)
		if err != nil {
			panic("failed to send presence to enter chat room :" + err.Error())
		}
	} else {
		correspondentJid = os.Args[correspondent]
	}

	m := stanza.Message{Attrs: stanza.Attrs{To: correspondentJid, Type: stanza.MessageTypeChat}, Body: os.Args[message]}
	err = client.Send(m)
	if err != nil {
		panic(err)
	}

	// After sending the action message, let's disconnect from the chat room if we were connected to one.
	if isCorrespRoom {
		pOut := stanza.NewPresence(
			stanza.Attrs{To: correspondentJid, Type: stanza.PresenceTypeUnavailable},
		)
		client.Send(pOut)
	}

	client.Disconnect()
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}
