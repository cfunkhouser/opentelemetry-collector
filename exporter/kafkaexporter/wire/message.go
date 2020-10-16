// Package wire contains types and functions to interact with the Kafka wire
// protocol.
package wire

// Message encapsulates Kafka's message payload.
type Message struct {
	Value []byte
}
