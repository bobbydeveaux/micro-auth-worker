package main

import (
	"github.com/dvsekhvalnov/jose2go"
	"github.com/nats-io/go-nats"
	"log"
	"os"
)

// need to get this from ENV, because GitHub public project will expose this. Oops.
const jwtPassphrase string = "fbac-FJfxeMQCzXBPqrIY8Hhk"
const atPassphrase string = "fbac-VviNJvCUqDK1v9BtMUop"

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
		log.Printf("Generating Access Token: %s\n", msg.Data)

		payload := msg.Data
		strPayload := string(payload[:])
		log.Printf("payload is %v, ", strPayload)

		secureToken, err := jose.Encrypt(strPayload, jose.PBES2_HS256_A128KW, jose.A256GCM, atPassphrase)
		if err != nil {
			log.Println("error:", err)
		}

		at.Value = secureToken
		ec.Publish(msg.Reply, at)
	})

	select {}
}
