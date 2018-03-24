package testapps

import (
	"testing"
	"github.com/Oppodelldog/disgostu/testapps/client"
	"github.com/Oppodelldog/disgostu/testapps/server"
	"time"
	"github.com/Oppodelldog/disgostu/capturing"
	"github.com/sirupsen/logrus"
	"github.com/Oppodelldog/disgostu/disgostutest"
)

func TestWhatever(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	// in normal test cases the server would of course be available
	// here I just simulate two servers responding two bytes
	go server.RunServerBehavior("0.0.0.0:19001")
	go server.RunServerBehavior("0.0.0.0:19002")
	time.Sleep(200 * time.Millisecond)

	// Start a new capturing session
	disgostutest.StartRecording(t,getCapturingConfig())

	// Now run the code under test, like server simply send two bytes
	client.RunClientBehavior("localhost:9001")
	client.RunClientBehavior("localhost:9002")
	time.Sleep(100 * time.Millisecond)

	disgostutest.Assert()
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

/**
	TcpDataCaptures: []capturing.TcpDataCapture{
			{RecordName: "CapturingSession - 001", TargetAddress: "localhost:19001", TimeOffset: 23233, To: "client->server", Data: []byte{74, 79}},
			{RecordName: "CapturingSession - 001", TargetAddress: "localhost:19001", TimeOffset: 1800065281, To: "client<-server", Data: []byte{72, 79}},
			{RecordName: "CapturingSession - 001", TargetAddress: "localhost:19001", TimeOffset: 30738, To: "client->server", Data: []byte{74, 79}},
			{RecordName: "CapturingSession - 001", TargetAddress: "localhost:19001", TimeOffset: 1800181702, To: "client<-server", Data: []byte{72, 79}},
		},
 */
