package xsse

type ServerSentEvents interface {
	GetEvent() (evt *Event, ok bool)
}
