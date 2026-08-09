package xsse

import (
	"encoding/json"
	"sync"
	"time"
)

var (
	// NewHeartbeatEvent is the default heartbeat event.
	NewHeartbeatEvent = defaultHeartbeatEvent

	// HeartbeatInterval is the default heartbeat interval.
	HeartbeatInterval time.Duration = time.Second * 115

	// heartBeatEvent is the default heartbeat event.
	heartBeatEvent SSEEvent
	heartBeatOnce  sync.Once
)

var (
	// default marshal callback
	DefaultMarshalCallback MarshalCallback = json.Marshal
	// default unmarshal callback
	DefaultUnmarshalCallback UnmarshalCallback = func(bytes []byte) (data any, err error) {
		return string(bytes), nil
	}
)

// GetHeartBeatEvent returns the default heartbeat event.
func GetHeartBeatEvent() SSEEvent {
	heartBeatOnce.Do(func() {
		heartBeatEvent = NewHeartbeatEvent("")
	})
	return heartBeatEvent
}
