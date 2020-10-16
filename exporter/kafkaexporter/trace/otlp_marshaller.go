// Copyright 2020 The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/exporter/kafkaexporter/wire"
	otlptrace "go.opentelemetry.io/collector/internal/data/opentelemetry-proto-gen/collector/trace/v1"
)

type OTLPProtoMarshaller struct {
}

var _ Marshaller = (*OTLPProtoMarshaller)(nil)

func (*OTLPProtoMarshaller) Encoding() string {
	return "otlp_proto"
}

func (*OTLPProtoMarshaller) Marshal(traces pdata.Traces) ([]wire.Message, error) {
	request := otlptrace.ExportTraceServiceRequest{
		ResourceSpans: pdata.TracesToOtlp(traces),
	}
	bts, err := request.Marshal()
	if err != nil {
		return nil, err
	}
	return []wire.Message{{Value: bts}}, nil
}
