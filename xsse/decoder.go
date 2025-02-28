package xsse

import (
	"bufio"
	"bytes"
	"io"
	"strings"

	"golang.org/x/sync/errgroup"
)

type decoder struct {
	unmarshalCallback UnmarshalCallback
}

func Decode(r io.Reader, opts ...DecoderOption) (<-chan Event, error) {
	var dec decoder = decoder{
		unmarshalCallback: DefaultUnmarshalCallback,
	}
	for i := range opts {
		opts[i](&dec)
	}
	return dec.decode(r)
}

func (d *decoder) dispatchEvent(events chan Event, event Event, data []byte) {
	dataLength := len(data)
	if dataLength > 0 {
		//If the data buffer's last character is a U+000A LINE FEED (LF) character, then remove the last character from the data buffer.
		data = data[:dataLength-1]
		dataLength--
	}
	if dataLength == 0 && event.Event == "" {
		return
	}
	if event.Event == "" {
		event.Event = "message"
	}
	if d.unmarshalCallback != nil {
		event.Data, event.Err = d.unmarshalCallback(data)
	}
	events <- event
}

func (d *decoder) decode(r io.Reader) (events chan Event, err error) {
	group := errgroup.Group{}
	events = make(chan Event)
	group.Go(func() error {
		defer func() {
			close(events)
		}()

		var dataBuffer *bytes.Buffer = new(bytes.Buffer)
		var currentEvent Event
		reader := bufio.NewReader(r) // 使用 bufio.NewReader 来读取数据
		for {
			// 读取每个事件
			line, err := reader.ReadString('\n') // 读取直到换行符，假设每个事件是一个独立的行
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				return err
			}

			line = strings.TrimSpace(line)
			if len(line) == 0 {
				// If the line is empty (a blank line). Dispatch the event.
				d.dispatchEvent(events, currentEvent, dataBuffer.Bytes())
				// reset current event and data buffer
				currentEvent = Event{}
				dataBuffer.Reset()
				continue
			}
			if line[0] == byte(':') {
				// If the line starts with a U+003A COLON character (:), ignore the line.
				continue
			}

			var field, value string
			colonIndex := strings.IndexRune(line, ':')
			if colonIndex != -1 {
				// If the line contains a U+003A COLON character character (:)
				// Collect the characters on the line before the first U+003A COLON character (:),
				// and let field be that string.
				field = line[:colonIndex]
				// Collect the characters on the line after the first U+003A COLON character (:),
				// and let value be that string.
				value = line[colonIndex+1:]
				// If value starts with a single U+0020 SPACE character, remove it from value.
				if len(value) > 0 && value[0] == ' ' {
					value = value[1:]
				}
			} else {
				// Otherwise, the string is not empty but does not contain a U+003A COLON character character (:)
				// Use the whole line as the field name, and the empty string as the field value.
				field = line
				value = ""
			}
			// The steps to process the field given a field name and a field value depend on the field name,
			// as given in the following list. Field names must be compared literally,
			// with no case folding performed.
			switch string(field) {
			case "event":
				// Set the event name buffer to field value.
				currentEvent.Event = string(value)
			case "id":
				// Set the event stream's last event ID to the field value.
				currentEvent.Id = string(value)
			case "retry":
				// If the field value consists of only characters in the range U+0030 DIGIT ZERO (0) to U+0039 DIGIT NINE (9),
				// then interpret the field value as an integer in base ten, and set the event stream's reconnection time to that integer.
				// Otherwise, ignore the field.
				currentEvent.Id = string(value)
			case "data":
				// Append the field value to the data buffer,
				dataBuffer.Write([]byte(value))
				// then append a single U+000A LINE FEED (LF) character to the data buffer.
				dataBuffer.WriteString("\n")
			default:
				//Otherwise. The field is ignored.
				continue
			}

		}
		// Once the end of the file is reached, the user agent must dispatch the event one final time.
		d.dispatchEvent(events, currentEvent, dataBuffer.Bytes())
		return nil
	})
	return events, nil
}
