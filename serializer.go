package watchdog

import "encoding/json"

type Marshaller interface {
	Unmarshal([]byte, any) error
	Marshal(any) ([]byte, error)
}

type JSONMarshaller struct{}

func (m *JSONMarshaller) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (m *JSONMarshaller) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}
