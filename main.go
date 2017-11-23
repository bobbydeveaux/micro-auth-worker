package main

import (
	"encoding/json"
	"github.com/dvsekhvalnov/jose2go"
	"github.com/nats-io/go-nats"
	"log"
	"os"
	"time"
)

// need to get this from ENV, because GitHub public project will expose this. Oops.
const jwtPassphrase string = "fbac-FJfxeMQCzXBPqrIY8Hhk"
const atPassphrase string = "fbac-VviNJvCUqDK1v9BtMUop"

type person struct {
	Id          int64
	Name        string
	Valid       bool
	Jwt         string
	AccessToken accessToken
}

type accessToken struct {
	Value  string
	Expiry int
}

type jwtAuth struct {
	User   person
	Expiry int
}

type jwtToken struct {
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

		at.Value = "FBAC." + secureToken
		ec.Publish(msg.Reply, at)
	})

	var authJwt jwtAuth
	var jwt jwtToken
	var p person
	ec.QueueSubscribe("auth.jwt", "job_workers", func(msg *nats.Msg) {
		log.Printf("Generating JWT: %s\n", msg.Data)

		// decode access token
		// validate access token
		// validate user
		// validate user acccess
		// generate JWT

		err := json.Unmarshal(msg.Data, &p)
		if err != nil {
			log.Println(err.Error())
		}

		authJwt.User = p
		authJwt.Expiry = int(time.Now().Unix()) + int(time.Second*600)

		payload, err := json.Marshal(authJwt)
		if err != nil {
			log.Println("error:", err)
		}
		strPayload := string(payload[:])
		log.Printf("payload is %v, ", strPayload)

		secureToken, err := jose.Encrypt(strPayload, jose.PBES2_HS256_A128KW, jose.A256GCM, jwtPassphrase)
		if err != nil {
			log.Println("error:", err)
		}

		jwt.Value = secureToken
		ec.Publish(msg.Reply, jwt)
	})

	select {}
}
