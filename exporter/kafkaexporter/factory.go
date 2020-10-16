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

package kafkaexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/exporter/kafkaexporter/config"
	"go.opentelemetry.io/collector/exporter/kafkaexporter/trace"
)

// FactoryOption applies changes to kafkaExporterFactory.
type FactoryOption func(factory *kafkaExporterFactory)

// WithTraceMarshaller adds trace Marshallers to the exporter.
func WithTraceMarshaller(encodingMarshaller map[string]trace.Marshaller) FactoryOption {
	return func(factory *kafkaExporterFactory) {
		for encoding, marshaller := range encodingMarshaller {
			factory.traceMarshallers[encoding] = marshaller
		}
	}
}

// NewFactory creates Kafka exporter factory.
func NewFactory(options ...FactoryOption) component.ExporterFactory {
	f := &kafkaExporterFactory{
		traceMarshallers: trace.DefaultMarshallers(),
	}
	for _, o := range options {
		o(f)
	}
	return exporterhelper.NewFactory(
		config.ExporterTypeName,
		config.Default,
		exporterhelper.WithTraces(f.createTraceExporter))
}

type kafkaExporterFactory struct {
	traceMarshallers map[string]trace.Marshaller
}

func (f *kafkaExporterFactory) createTraceExporter(
	_ context.Context,
	params component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.TraceExporter, error) {
	oCfg := cfg.(*config.Config)
	exp, err := newExporter(*oCfg, params, f.traceMarshallers)
	if err != nil {
		return nil, err
	}
	return exporterhelper.NewTraceExporter(
		cfg,
		exp.traceDataPusher,
		// Disable exporterhelper Timeout, because we cannot pass a Context to the Producer,
		// and will rely on the sarama Producer Timeout logic.
		exporterhelper.WithTimeout(exporterhelper.TimeoutSettings{Timeout: 0}),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(exp.Close))
}
