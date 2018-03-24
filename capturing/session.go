package capturing

import (
	"context"
)

func NewSession(config CapturingConfig) *Session {
	return &Session{
		capturingConfig: config,
	}
}

type Session struct {
	capturingConfig CapturingConfig
	tcpDataCapture  TcpDataCapture
	cancelProxies   context.CancelFunc
}

func (s *Session) Start() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	s.cancelProxies = cancelFunc

	for _, proxyConfig := range s.capturingConfig.ProxyConfigs {
		proxyConfig.RecordName = s.capturingConfig.RecordName
		StartProxy(proxyConfig, ctx)
	}
}

func (s *Session) Stop() {
	s.cancelProxies()
}

func (s *Session) GetCapturedData() (StoreCapture, bool) {
	return GetCapture(s.capturingConfig.RecordName)
}
