package msgbroker

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/my/app/server/config"
)

// Init .
func Init(env *config.EnvConfig) mqtt.Client {
	opts := createClientOptions(env)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Panic(err)
	}
	return client
}

func createClientOptions(env *config.EnvConfig) *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s", env.MQTTHost))
	opts.SetUsername(env.MQTTUser)
	opts.SetPassword(env.MQTTPass)
	opts.SetClientID("pub")
	return opts
}
