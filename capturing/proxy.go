package capturing

import (
	"io"
	"net"
	"github.com/sirupsen/logrus"
)

func NewCapturingProxy(recordName string, connectionName string, targetAddress string, readConnection *net.TCPConn, writeConnection *net.TCPConn, captureChan chan TcpDataCapture, shutdownChannel chan bool) chan bool {
	closingChannel := make(chan bool)
	go func() {
		isRunning := true
		for isRunning {

			select {
			case _, ok := <-shutdownChannel:
				if !ok {
					isRunning = false
					break
				}
			default:

				data := make([]byte, 256)
				n, err := readConnection.Read(data)
				if err != nil {
					if err == io.EOF {
						isRunning = false
						break
					} else {
						logrus.Debugf("error while reading %s: %v", connectionName, err)
						isRunning = false
						break
					}
				}
				if _, err := writeConnection.Write(data[:n]); err != nil {
					logrus.Debugf("error while writing %s: %v", connectionName, err)
					isRunning = false
					break
				}
				captureChan <- TcpDataCapture{RecordName: recordName, Data: data[:n], To: connectionName, TargetAddress: targetAddress}
			}
		}
		logrus.Debugf("%s worker shutdown", connectionName)
		closingChannel <- true
	}()

	return closingChannel
}
