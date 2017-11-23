package main

import (
	"github.com/nats-io/go-nats"
	"log"
	"os"
)

// need to get this from ENV, because GitHub public project will expose this. Oops.
const passphrase string = "fbac-FJfxeMQCzXBPqrIY8Hhk"

type person struct {
	Id          int64
	Name        string
	Valid       bool
	Jwt         string
	AccessToken string
}

type accessToken struct {
	Value string
}

func main() {
	//NATS_HOST = nats://localhost:4222
	nc, _ := nats.Connect(os.Getenv("NATS_HOST"))
	ec, _ := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	defer ec.Close()

	var at accessToken
	ec.QueueSubscribe("auth.generateaccesstoken", "job_workers", func(msg *nats.Msg) {
		log.Printf("Authenticating: %s\n", msg.Data)

		at.Value = "FBAC-123456"
		ec.Publish(msg.Reply, at)
	})

	select {}
}
