// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package config handles configuration for the kafkaexporter.
package config

import (
	"time"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	// ExporterTypeName of this exporter.
	ExporterTypeName = "kafka"

	// DefaultTopic name for exporting to Kafka.
	DefaultTopic = "otlp"

	// DefaultEncoding used to export to Kafka.
	DefaultEncoding = "otlp_proto"

	// DefaultBroker for Kafka exports.
	DefaultBroker = "localhost:9092"

	// The following defaults are copied from sarama.NewConfig()

	// DefaultMetadataRetryMax number of retries submitting to the broker.
	DefaultMetadataRetryMax = 3

	// DefaultMetadataRetryBackoff time between retries submitting to the broker.
	DefaultMetadataRetryBackoff = time.Millisecond * 250

	// DefaultMetadataFull is to send the full metadata.
	DefaultMetadataFull = true
)

// Config defines configuration for Kafka exporter.
type Config struct {
	configmodels.ExporterSettings  `mapstructure:",squash"`
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`

	// The list of kafka brokers (default localhost:9092)
	Brokers []string `mapstructure:"brokers"`
	// Kafka protocol version
	ProtocolVersion string `mapstructure:"protocol_version"`
	// The name of the kafka topic to export to (default "otlp")
	Topic string `mapstructure:"topic"`
	// Encoding of the messages (default "otlp_proto")
	Encoding string `mapstructure:"encoding"`

	// Metadata is the namespace for metadata management properties used by the
	// Client, and shared by the Producer/Consumer.
	Metadata Metadata `mapstructure:"metadata"`

	// Authentication defines used authentication mechanism.
	Authentication Authentication `mapstructure:"auth"`
}

// Metadata defines configuration for retrieving metadata from the broker.
type Metadata struct {
	// Whether to maintain a full set of metadata for all topics, or just
	// the minimal set that has been necessary so far. The full set is simpler
	// and usually more convenient, but can take up a substantial amount of
	// memory if you have many topics and partitions. Defaults to true.
	Full bool `mapstructure:"full"`

	// Retry configuration for metadata.
	// This configuration is useful to avoid race conditions when broker
	// is starting at the same time as collector.
	Retry MetadataRetry `mapstructure:"retry"`
}

// MetadataRetry defines retry configuration for Metadata.
type MetadataRetry struct {
	// The total number of times to retry a metadata request when the
	// cluster is in the middle of a leader election or at startup (default 3).
	Max int `mapstructure:"max"`
	// How long to wait for leader election to occur before retrying
	// (default 250ms). Similar to the JVM's `retry.backoff.ms`.
	Backoff time.Duration `mapstructure:"backoff"`
}

// Default exporter configuration.
func Default() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: ExporterTypeName,
			NameVal: ExporterTypeName,
		},
		TimeoutSettings: exporterhelper.CreateDefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.CreateDefaultRetrySettings(),
		QueueSettings:   exporterhelper.CreateDefaultQueueSettings(),
		Brokers:         []string{DefaultBroker},
		Topic:           DefaultTopic,
		Encoding:        DefaultEncoding,
		Metadata: Metadata{
			Full: DefaultMetadataFull,
			Retry: MetadataRetry{
				Max:     DefaultMetadataRetryMax,
				Backoff: DefaultMetadataRetryBackoff,
			},
		},
	}
}
