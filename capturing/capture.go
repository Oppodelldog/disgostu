package capturing

import (
	"sync"
	"time"
	"github.com/sirupsen/logrus"
)

func RunCaptureWorker() (chan TcpDataCapture, chan TcpDataCapture) {
	clientCaptureChan := make(chan TcpDataCapture)
	serverCaptureChan := make(chan TcpDataCapture)

	go func(clientCaptureChan chan TcpDataCapture, serverCaptureChan <-chan TcpDataCapture) {
		alphaTime := time.Now().UnixNano()
		for captureData := range merge(clientCaptureChan, serverCaptureChan) {
			captureData.TimeOffset = time.Now().UnixNano() - alphaTime
			logrus.Debugf("%+v", captureData)
			AddCapture(captureData)
		}
		logrus.Debug("capture worker shutdown")

	}(clientCaptureChan, serverCaptureChan)

	return clientCaptureChan, serverCaptureChan
}

func merge(cs ...<-chan TcpDataCapture) <-chan TcpDataCapture {
	var wg sync.WaitGroup
	out := make(chan TcpDataCapture)

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan TcpDataCapture) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
