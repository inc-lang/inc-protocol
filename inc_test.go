package inccore

import (
	"testing"

	"google.golang.org/protobuf/proto"
)

func TestVariableOneofRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		v    *Variable
	}{
		{"int", &Variable{Kind: &Variable_IntVal{IntVal: 42}}},
		{"float", &Variable{Kind: &Variable_FloatVal{FloatVal: 3.14}}},
		{"string", &Variable{Kind: &Variable_StringVal{StringVal: "hello"}}},
		{"bool", &Variable{Kind: &Variable_BoolVal{BoolVal: true}}},
		{"bytes", &Variable{Kind: &Variable_BinaryVal{BinaryVal: []byte{0xDE, 0xAD}}}},
		{"null", &Variable{Kind: &Variable_NullVal{NullVal: &NullValue{}}}},
		{"list", &Variable{Kind: &Variable_ListVal{ListVal: &ListValue{
			Values: []*Variable{
				{Kind: &Variable_IntVal{IntVal: 1}},
				{Kind: &Variable_StringVal{StringVal: "two"}},
			},
		}}}},
		{"map", &Variable{Kind: &Variable_MapVal{MapVal: &MapValue{
			Pairs: map[string]*Variable{
				"key": {Kind: &Variable_IntVal{IntVal: 99}},
			},
		}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := proto.Marshal(tt.v)
			if err != nil {
				t.Fatalf("marshal: %v", err)
			}

			got := &Variable{}
			if err := proto.Unmarshal(data, got); err != nil {
				t.Fatalf("unmarshal: %v", err)
			}

			if !proto.Equal(tt.v, got) {
				t.Errorf("round-trip mismatch: got %v, want %v", got, tt.v)
			}
		})
	}
}

func TestExecutionRequestSerialization(t *testing.T) {
	req := &ExecutionRequest{
		RequestId: "test-123",
		Code:      "print(x + y)",
		Variables: map[string]*Variable{
			"x": {Kind: &Variable_IntVal{IntVal: 10}},
			"y": {Kind: &Variable_IntVal{IntVal: 20}},
		},
		Config: map[string]string{
			"memory_limit": "128M",
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &ExecutionRequest{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !proto.Equal(req, got) {
		t.Errorf("round-trip mismatch")
	}

	if got.RequestId != "test-123" {
		t.Errorf("request_id: got %q, want %q", got.RequestId, "test-123")
	}
	if got.Code != "print(x + y)" {
		t.Errorf("code: got %q, want %q", got.Code, "print(x + y)")
	}
	if len(got.Variables) != 2 {
		t.Errorf("variables count: got %d, want 2", len(got.Variables))
	}
	if got.Config["memory_limit"] != "128M" {
		t.Errorf("config memory_limit: got %q, want %q", got.Config["memory_limit"], "128M")
	}
}

func TestExecutionResponseFields(t *testing.T) {
	resp := &ExecutionResponse{
		RequestId: "test-456",
		Result:    &Variable{Kind: &Variable_StringVal{StringVal: "done"}},
		Stdout:    "hello world\n",
		Stderr:    "",
		ExitCode:  0,
	}

	data, err := proto.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &ExecutionResponse{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.RequestId != "test-456" {
		t.Errorf("request_id: got %q, want %q", got.RequestId, "test-456")
	}
	if got.Stdout != "hello world\n" {
		t.Errorf("stdout: got %q, want %q", got.Stdout, "hello world\n")
	}
	if got.ExitCode != 0 {
		t.Errorf("exit_code: got %d, want 0", got.ExitCode)
	}
	if !proto.Equal(resp, got) {
		t.Errorf("round-trip mismatch")
	}
}

func TestMapValueAndListValue(t *testing.T) {
	nested := &Variable{Kind: &Variable_MapVal{MapVal: &MapValue{
		Pairs: map[string]*Variable{
			"numbers": {Kind: &Variable_ListVal{ListVal: &ListValue{
				Values: []*Variable{
					{Kind: &Variable_IntVal{IntVal: 1}},
					{Kind: &Variable_IntVal{IntVal: 2}},
					{Kind: &Variable_IntVal{IntVal: 3}},
				},
			}}},
			"name": {Kind: &Variable_StringVal{StringVal: "test"}},
			"nested_map": {Kind: &Variable_MapVal{MapVal: &MapValue{
				Pairs: map[string]*Variable{
					"inner": {Kind: &Variable_BoolVal{BoolVal: true}},
				},
			}}},
		},
	}}}

	data, err := proto.Marshal(nested)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got := &Variable{}
	if err := proto.Unmarshal(data, got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if !proto.Equal(nested, got) {
		t.Errorf("round-trip mismatch for nested map/list structure")
	}

	mapVal := got.GetMapVal()
	if mapVal == nil {
		t.Fatal("expected MapValue, got nil")
	}
	if len(mapVal.Pairs) != 3 {
		t.Errorf("pairs count: got %d, want 3", len(mapVal.Pairs))
	}

	numbers := mapVal.Pairs["numbers"].GetListVal()
	if numbers == nil {
		t.Fatal("expected ListValue for 'numbers', got nil")
	}
	if len(numbers.Values) != 3 {
		t.Errorf("list length: got %d, want 3", len(numbers.Values))
	}
}
