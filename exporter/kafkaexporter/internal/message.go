// Package internal contains types and functions internal to the kafkaexporter
// implementation, and should not be used by outside code.
package internal

// Message encapsulates Kafka's message payload.
type Message struct {
	Value []byte
}
