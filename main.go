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
	var correspondentJid *stanza.Jid
	isCorrespRoom, err := strconv.ParseBool(os.Args[correspondent_is_room])
	if err == nil && isCorrespRoom {
		correspondentJid, err = stanza.NewJid(os.Args[correspondent] + "@" + os.Args[serverDomain] + "/github_bot")
		if err != nil {
			panic(err)
		}
		// Sending room presence
		joinMUC(client, correspondentJid)
	} else {
		correspondentJid, err = stanza.NewJid(os.Args[correspondent])
		if err != nil {
			panic(err)
		}
	}

	m := stanza.Message{Attrs: stanza.Attrs{To: correspondentJid.Bare(), Type: stanza.MessageTypeChat}, Body: os.Args[message]}
	err = client.Send(m)
	if err != nil {
		panic(err)
	}

	// After sending the action message, let's disconnect from the chat room if we were connected to one.
	if isCorrespRoom {
		leaveMUC(client, correspondentJid)
	}

	client.Disconnect()
}

func errorHandler(err error) {
	fmt.Println(err.Error())
}

func joinMUC(c xmpp.Sender, toJID *stanza.Jid) error {
	return c.Send(stanza.Presence{Attrs: stanza.Attrs{To: toJID.Full()},
		Extensions: []stanza.PresExtension{
			stanza.MucPresence{
				History: stanza.History{MaxStanzas: stanza.NewNullableInt(0)},
			}},
	})
}

func leaveMUC(c xmpp.Sender, muc *stanza.Jid) {
	c.Send(stanza.Presence{Attrs: stanza.Attrs{
		To:   muc.Full(),
		Type: stanza.PresenceTypeUnavailable,
	}})

}
