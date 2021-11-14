package protorand

import (
	"testing"

	testpb "github.com/sryoya/protorand/testdata"
)

func TestEmbedValues(t *testing.T) {
	p := New()
	p.Seed(0)

	target := &testpb.TestMessage{}

	err := p.EmbedValues(target)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// assert all the fields got some value
	if target.SomeInt32 == 0 {
		t.Errorf("Field SomeInt32 is not set")
	}
	if target.SomeSint32 == 0 {
		t.Errorf("Field SomeSint32 is not set")
	}
	if target.SomeUint32 == 0 {
		t.Errorf("Field SomeUint32 is not set")
	}
	if target.SomeFloat32 == 0 {
		t.Errorf("Field SomeFloat32 is not set")
	}
	if target.SomeFixed32 == 0 {
		t.Errorf("Field SomeFixed32 is not set")
	}
	if target.SomeSfixed32 == 0 {
		t.Errorf("Field SomeSfixed32 is not set")
	}
	if target.SomeInt64 == 0 {
		t.Errorf("Field SomeInt64 is not set")
	}
	if target.SomeSint64 == 0 {
		t.Errorf("Field SomeSint64 is not set")
	}
	if target.SomeUint64 == 0 {
		t.Errorf("Field SomeUint64 is not set")
	}
	if target.SomeFloat64 == 0 {
		t.Errorf("Field SomeFloat64 is not set")
	}
	if target.SomeFixed64 == 0 {
		t.Errorf("Field SomeFixed64 is not set")
	}
	if target.SomeSfixed64 == 0 {
		t.Errorf("Field SomeSfixed64 is not set")
	}
	if target.SomeStr == "" {
		t.Errorf("Field SomeStr is not set")
	}
	if target.SomeMsg == nil {
		t.Errorf("Field SomeMsg is not set")
	}
	if len(target.SomeSlice) == 0 {
		t.Errorf("Field SomeSlice is not set")
	}
	if len(target.SomeMsgs) == 0 {
		t.Errorf("Field SomeMsgs is not set")
	}
	if len(target.SomeMap) == 0 {
		t.Errorf("Field SomeMap is not set")
	}
	if target.SomeEnum > 3 { // undeclared enum value
		t.Errorf("Field SomeEnum is not set")
	}
	if target.SomeEnum2 > 1 { // undeclared enum value
		t.Errorf("Field SomeEnum2 is not set")
	}
	if target.SomeOneOf == nil {
		t.Errorf("Field SomeOneOf is not set")
	}
}
