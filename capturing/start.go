package capturing

import (
	"github.com/sirupsen/logrus"
	"net"
	"context"
)

const numberOfConcurrentConnections = 1

type CapturingConfig struct {
	RecordName   string
	ProxyConfigs []ProxyConfig
}

type ProxyConfig struct {
	RecordName    string
	TargetAddress string
	ProxyAddress  string
}

type ProxyConnection struct {
	sourceConnection *net.TCPConn
	targetAddress    string
	RecordName       string
}

func handleProxyConnection(proxyConnection ProxyConnection) {

	logrus.Debug("New proxy connection to target server")
	targetConnection, err := connectToTargetServer(proxyConnection)
	if err != nil {
		panic(err)
	}
	defer targetConnection.Close()

	clientCaptureChan, serverCaptureChan := RunCaptureWorker()
	clientShutDownChan := make(chan bool)
	serverShutDownChan := make(chan bool)

	clientClosingChan := NewCapturingProxy(proxyConnection.RecordName, "client->server", proxyConnection.targetAddress, proxyConnection.sourceConnection, targetConnection, clientCaptureChan, clientShutDownChan)
	serverClosingChan := NewCapturingProxy(proxyConnection.RecordName, "client<-server", proxyConnection.targetAddress, targetConnection, proxyConnection.sourceConnection, serverCaptureChan, serverShutDownChan)

	select {
	case <-clientClosingChan:
		close(serverShutDownChan)
	case <-serverClosingChan:
		// nothing to to, client connection is closed when this function returns
		close(clientShutDownChan)
	}

	logrus.Debug("proxy shutdown")
}

func connectToTargetServer(proxyConnection ProxyConnection) (*net.TCPConn, error) {
	rAddr, err := net.ResolveTCPAddr("tcp", proxyConnection.targetAddress)
	if err != nil {
		return nil, err
	}
	targetConnection, err := net.DialTCP("tcp", nil, rAddr)
	if err != nil {
		return nil, err
	}

	return targetConnection, nil
}

func handleConn(in <-chan ProxyConnection, out chan<- ProxyConnection) {
	for proxyConnection := range in {
		logrus.Debug("handle client connection")
		handleProxyConnection(proxyConnection)
		out <- proxyConnection
	}
	logrus.Debug("connection handler finished")
}

func closeConn(in <-chan ProxyConnection) {
	for proxyConnection := range in {
		proxyConnection.sourceConnection.Close()
		logrus.Debug("Closed client connection")
	}
	logrus.Debug("connection closer finished")
}

func StartProxy(proxyConfig ProxyConfig, ctx context.Context) {
	logrus.Debugf("Staring CapturingProxy: %v\nProxying: %v\n\n", proxyConfig.ProxyAddress, proxyConfig.TargetAddress)
	addr, err := net.ResolveTCPAddr("tcp", proxyConfig.ProxyAddress)
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	pending, complete := make(chan ProxyConnection), make(chan ProxyConnection)

	for i := 0; i < numberOfConcurrentConnections; i++ {
		go handleConn(pending, complete)
	}

	go closeConn(complete)

	go func() {
		isRunning := true
		for isRunning {
			select {
			case <-ctx.Done():
				isRunning = false
				break
			default:


				conn, err := listener.AcceptTCP()
				if err != nil {
					panic(err)
				}

				pending <- ProxyConnection{
					sourceConnection: conn,
					targetAddress:    proxyConfig.TargetAddress,
					RecordName:       proxyConfig.RecordName,
				}
			}
		}
		logrus.Debug("shutting down proxy listener")
		listener.Close()
	}()
}
