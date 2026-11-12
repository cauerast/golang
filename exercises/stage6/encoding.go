package stage6

import "context"

// Message is a JSON envelope. Note the struct tags: they control JSON names,
// and omitempty drops the field when it holds its zero value.
type Message struct {
	Type string `json:"type"`
	Body string `json:"body,omitempty"`
	Size int64  `json:"size"`
}

// EncodeMessage marshals m to JSON.
// Example: Message{Type: "file", Size: 10} -> {"type":"file","size":10}
//
// TODO: json.Marshal. Only EXPORTED (capitalized) fields are encoded: why?
func EncodeMessage(m Message) ([]byte, error) {
	return nil, nil
}

// DecodeMessage unmarshals data into a Message. Unknown JSON fields must be
// rejected with an error (strict decoding).
//
// TODO: json.NewDecoder(bytes.NewReader(data)) + DisallowUnknownFields().
// json.Unmarshal alone silently ignores unknown fields.
func DecodeMessage(data []byte) (Message, error) {
	return Message{}, nil
}

// CollectUntilDone receives values from in and returns them in order. It stops
// and returns what it has collected when in is closed OR ctx is cancelled,
// whichever happens first.
//
// TODO: for/select over ctx.Done() and in. A bare `for v := range in` would
// ignore cancellation and could block forever.
func CollectUntilDone(ctx context.Context, in <-chan int) []int {
	return nil
}
