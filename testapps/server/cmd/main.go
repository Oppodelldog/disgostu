package main

import (
	"time"
	"github.com/Oppodelldog/disgostu/testapps/server"
)

func main() {

	for {
		time.Sleep(time.Second * 2)
		server.RunServerBehavior("0.0.0.0:19001")
	}
}
