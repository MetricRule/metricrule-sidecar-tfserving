package tfmetric

import (
	"reflect"
	"testing"

	"go.opentelemetry.io/otel/metric"
	"google.golang.org/protobuf/encoding/prototext"

	configpb "github.com/metricrule-sidecar-tfserving/api/proto/metricconfigpb"
)

func TestInputCounterInstrumentSpec(t *testing.T) {
	configTextProto := `
		input_metrics {
			name: "simple"
			simple_counter: {}
		}`
	var config configpb.SidecarConfig
	prototext.Unmarshal([]byte(configTextProto), &config)

	specs := GetInstrumentSpecs(&config)

	gotInputLen := len(specs[InputContext])
	wantInputLen := 1
	if gotInputLen != wantInputLen {
		t.Errorf("Unexpected length of input context specs, got %v, wanted %v", gotInputLen, wantInputLen)
	}

	gotOutputLen := len(specs[OutputContext])
	wantOutputLen := 0
	if gotOutputLen != wantOutputLen {
		t.Errorf("Unexpected length of output context specs, got %v, wanted %v", gotOutputLen, wantOutputLen)
	}

	if gotInputLen < 1 {
		return
	}

	gotSpecInstance := specs[InputContext][0]
	wantInstrumentKind := metric.CounterInstrumentKind
	wantMetricKind := reflect.Int64
	wantMetricName := "simple"
	if gotSpecInstance.InstrumentKind != wantInstrumentKind {
		t.Errorf("Unexpected instrument kind in spec, got %v, wanted %v", gotSpecInstance.InstrumentKind, wantInstrumentKind)
	}
	if gotSpecInstance.MetricValueKind != wantMetricKind {
		t.Errorf("Unexpected metric kind in spec, got %v, wanted %v", gotSpecInstance.MetricValueKind, wantMetricKind)
	}
	if gotSpecInstance.Name != wantMetricName {
		t.Errorf("Unexpected metric name in spec, got %v, wanted %v", gotSpecInstance.Name, wantMetricName)
	}
}

func TestInputCounterMetrics(t *testing.T) {
	configTextProto := `
		input_metrics {
			simple_counter: {}
		}`
	var config configpb.SidecarConfig
	prototext.Unmarshal([]byte(configTextProto), &config)

	metrics := GetMetricInstances(&config, "{}", InputContext)

	gotLen := len(metrics)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Unexpected length of metrics, got %v, wanted %v", gotLen, wantLen)
	}

	if gotLen == 0 {
		return
	}

	counter := 0
	for spec, instance := range metrics {
		if counter >= wantLen {
			t.Errorf("Exceeded expected iteration length: %v", wantLen)
		}

		gotInstrumentKind := spec.InstrumentKind
		wantInstrumentKind := metric.CounterInstrumentKind
		if gotInstrumentKind != wantInstrumentKind {
			t.Errorf("Unexpected instrument kind, got %v, wanted %v", gotInstrumentKind, wantInstrumentKind)
		}

		gotMetricKind := spec.MetricValueKind
		wantMetricKind := reflect.Int64
		if gotMetricKind != wantMetricKind {
			t.Errorf("Unexpected metric kind, got %v, wanted %v", gotMetricKind, wantMetricKind)
		}

		gotValue := instance.MetricValue
		wantValue := int64(1)
		if gotValue != wantValue {
			t.Errorf("Unexpected metric value, got %v, wanted %v", gotValue, wantValue)
		}

		gotLabelsLen := len(instance.Labels)
		wantLabelsLen := 0
		if gotLabelsLen != wantLabelsLen {
			t.Errorf("Unexpected labels length, got %v, wanted %v", gotLabelsLen, wantLabelsLen)
		}
	}
}

func TestInputCounterWithLabels(t *testing.T) {
	configTextProto := `
		input_metrics {
			simple_counter: {}
			labels: {
				label_key: { string_value: "Application" }
				label_value: { string_value: "MetricRule" }
			}
		}`
	var config configpb.SidecarConfig
	prototext.Unmarshal([]byte(configTextProto), &config)

	metrics := GetMetricInstances(&config, "{}", InputContext)

	gotLen := len(metrics)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Unexpected length of metrics, got %v, wanted %v", gotLen, wantLen)
	}

	if gotLen == 0 {
		return
	}

	counter := 0
	for spec, instance := range metrics {
		if counter >= wantLen {
			t.Errorf("Exceeded expected iteration length: %v", wantLen)
		}

		gotInstrumentKind := spec.InstrumentKind
		wantInstrumentKind := metric.CounterInstrumentKind
		if gotInstrumentKind != wantInstrumentKind {
			t.Errorf("Unexpected metric kind, got %v, wanted %v", gotInstrumentKind, wantInstrumentKind)
		}

		gotMetricKind := spec.MetricValueKind
		wantMetricKind := reflect.Int64
		if gotMetricKind != wantMetricKind {
			t.Errorf("Unexpected metric kind, got %v, wanted %v", gotMetricKind, wantMetricKind)
		}

		gotValue := instance.MetricValue
		wantValue := int64(1)
		if gotValue != wantValue {
			t.Errorf("Unexpected metric value, got %v, wanted %v", gotValue, wantValue)
		}

		gotLabelsLen := len(instance.Labels)
		wantLabelsLen := 1
		if gotLabelsLen != wantLabelsLen {
			t.Errorf("Unexpected labels length, got %v, wanted %v", gotLabelsLen, wantLabelsLen)
		}

		if gotLabelsLen == 0 {
			return
		}

		gotLabel := instance.Labels[0]
		wantLabelKey := "Application"
		wantLabelValue := "MetricRule"
		if string(gotLabel.Key) != wantLabelKey {
			t.Errorf("Unexpected label key, got %v, wanted %v", gotLabel.Key, wantLabelKey)
		}
		if gotLabel.Value.AsString() != wantLabelValue {
			t.Errorf("Unexpected label key, got %v, wanted %v", gotLabel.Value.AsString(), wantLabelValue)
		}
	}
}

func TestOutputValues(t *testing.T) {
	configTextProto := `
		output_metrics {
			value {
				value {
					parsed_value {
						field_path: {
							paths: "prediction"
						}
						parsed_type: FLOAT
					}
				}
			}
		}`
	var config configpb.SidecarConfig
	prototext.Unmarshal([]byte(configTextProto), &config)

	metrics := GetMetricInstances(&config, "{ \"prediction\": 0.495 }", OutputContext)

	gotLen := len(metrics)
	wantLen := 1
	if gotLen != wantLen {
		t.Errorf("Unexpected length of metrics, got %v, wanted %v", gotLen, wantLen)
	}

	if gotLen == 0 {
		return
	}

	counter := 0
	for spec, instance := range metrics {
		if counter >= wantLen {
			t.Errorf("Exceeded expected iteration length: %v", wantLen)
		}

		gotInstrumentKind := spec.InstrumentKind
		wantInstrumentKind := metric.ValueRecorderInstrumentKind
		if gotInstrumentKind != wantInstrumentKind {
			t.Errorf("Unexpected metric kind, got %v, wanted %v", gotInstrumentKind, wantInstrumentKind)
		}

		gotMetricKind := spec.MetricValueKind
		wantMetricKind := reflect.Float64
		if gotMetricKind != wantMetricKind {
			t.Errorf("Unexpected metric kind, got %v, wanted %v", gotMetricKind, wantMetricKind)
		}

		gotValue := instance.MetricValue
		wantValue := 0.495
		if gotValue != wantValue {
			t.Errorf("Unexpected metric value, got %v, wanted %v", gotValue, wantValue)
		}

		gotLabelsLen := len(instance.Labels)
		wantLabelsLen := 0
		if gotLabelsLen != wantLabelsLen {
			t.Errorf("Unexpected labels length, got %v, wanted %v", gotLabelsLen, wantLabelsLen)
		}
	}
}
