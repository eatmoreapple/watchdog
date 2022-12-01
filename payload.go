package watchdog

import "encoding/json"

type Payload struct {
	// TaskID is the unique identifier of the task.
	TaskID     string
	RoutingKey string
	// Payload is the payload of the task.
	Data []byte
}

func (p Payload) MarshalBinary() (data []byte, err error) {
	return json.Marshal(p)
}
