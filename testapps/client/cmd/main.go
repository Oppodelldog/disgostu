package main

import (
	"time"
	"github.com/Oppodelldog/disgostu/testapps/client"
)

func main() {

	for {
		time.Sleep(time.Second * 2)
		client.RunClientBehavior("localhost:9999")
	}
}
