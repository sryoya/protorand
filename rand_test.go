package protorand

import (
	"testing"

	"google.golang.org/protobuf/proto"

	testpb "github.com/wattch/protorand/testdata"
)

func TestEmbedValues(t *testing.T) {
	p := New()
	p.Seed(0)

	input := &testpb.TestMessage{}
	res, err := p.Gen(input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// check if the input is not mutated.
	if !proto.Equal(input, &testpb.TestMessage{}) {
		t.Errorf("The input was unexpectedly mutated.")
	}

	got := res.(*testpb.TestMessage)

	// assert all the fields got some value
	if got.SomeInt32 == 0 {
		t.Errorf("Field SomeInt32 is not set")
	}
	if got.SomeSint32 == 0 {
		t.Errorf("Field SomeSint32 is not set")
	}
	if got.SomeUint32 == 0 {
		t.Errorf("Field SomeUint32 is not set")
	}
	if got.SomeFloat32 == 0 {
		t.Errorf("Field SomeFloat32 is not set")
	}
	if got.SomeFixed32 == 0 {
		t.Errorf("Field SomeFixed32 is not set")
	}
	if got.SomeSfixed32 == 0 {
		t.Errorf("Field SomeSfixed32 is not set")
	}
	if got.SomeInt64 == 0 {
		t.Errorf("Field SomeInt64 is not set")
	}
	if got.SomeSint64 == 0 {
		t.Errorf("Field SomeSint64 is not set")
	}
	if got.SomeUint64 == 0 {
		t.Errorf("Field SomeUint64 is not set")
	}
	if got.SomeFloat64 == 0 {
		t.Errorf("Field SomeFloat64 is not set")
	}
	if got.SomeFixed64 == 0 {
		t.Errorf("Field SomeFixed64 is not set")
	}
	if got.SomeSfixed64 == 0 {
		t.Errorf("Field SomeSfixed64 is not set")
	}
	if got.SomeStr == "" {
		t.Errorf("Field SomeStr is not set")
	}
	if got.SomeMsg == nil {
		t.Errorf("Field SomeMsg is not set")
	}
	if len(got.SomeSlice) == 0 {
		t.Errorf("Field SomeSlice is not set")
	}
	if len(got.SomeMsgs) == 0 {
		t.Errorf("Field SomeMsgs is not set")
	}
	if len(got.SomeMap) == 0 {
		t.Errorf("Field SomeMap is not set")
	}
	if got.SomeEnum > 3 { // undeclared enum value
		t.Errorf("Field SomeEnum is not set")
	}
	if got.SomeEnum2 > 1 { // undeclared enum value
		t.Errorf("Field SomeEnum2 is not set")
	}
	if got.SomeOneOf == nil {
		t.Errorf("Field SomeOneOf is not set")
	}
}
