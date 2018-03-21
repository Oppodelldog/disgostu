package testapps

import (
	"testing"
	"github.com/Oppodelldog/disgostu/testapps/client"
	"github.com/Oppodelldog/disgostu/testapps/server"
	"time"
	"github.com/Oppodelldog/disgostu/capturing"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"fmt"
)

func TestWhatever(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	capturing.StartCapture(getCapturingConfig())


	go server.RunServerBehavior("0.0.0.0:19001")
	time.Sleep(200 * time.Millisecond)


	client.RunClientBehavior("localhost:9001")
	go server.RunServerBehavior("0.0.0.0:19001")
	time.Sleep(200 * time.Millisecond)
	client.RunClientBehavior("localhost:9001")

	time.Sleep(100*time.Millisecond)

	data, err := capturing.GetCapture("CapturingSession-001")
	assert.True(t, err)
	fmt.Printf("%+v", data)
}

func getCapturingConfig() capturing.CapturingConfig {
	return capturing.CapturingConfig{
		RecordName: "CapturingSession-001",
		ProxyConfigs: []capturing.ProxyConfig{
			{
				TargetAddress: "localhost:19001",
				ProxyAddress:  "0.0.0.0:9001",
			},
			{
				TargetAddress: "localhost:19002",
				ProxyAddress:  "0.0.0.0:9002",
			},
		},
	}
}
