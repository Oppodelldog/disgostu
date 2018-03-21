package capturing

import (
	"sync"
)

var store sync.Map

type TcpDataCapture struct {
	RecordName    string
	TargetAddress string
	TimeOffset    int64
	To            string
	Data          []byte
}

type StoreCapture struct {
	TcpDataCaptures []TcpDataCapture
}

func AddCapture(capture TcpDataCapture) {
	if storeCapture, ok := GetCapture(capture.RecordName); ok {
		storeCapture.TcpDataCaptures = append(storeCapture.TcpDataCaptures, capture)
		store.Store(capture.RecordName, storeCapture)
	} else {
		storeCapture = StoreCapture{
			TcpDataCaptures: []TcpDataCapture{capture},
		}
		store.Store(capture.RecordName, storeCapture)
	}
}

func GetCapture(recordName string) (StoreCapture, bool) {
	if sth, ok := store.Load(recordName); ok {
		return sth.(StoreCapture), true
	}
	return StoreCapture{}, false
}
