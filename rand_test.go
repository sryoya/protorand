package protorand

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	testpb "github.com/sryoya/protorand/testdata"
)

func init() {
	// inject random generated values to be fixed
	randomInt32 = func() int32 { return 10 }
	randomUint32 = func() uint32 { return 11 }
	randomFloat32 = func() float32 { return 10.1 }
	randomInt64 = func() int64 { return 20 }
	randomUint64 = func() uint64 { return 21 }
	randomFloat64 = func() float64 { return 20.22 }
	randomString = func(int) string { return "Gopher" }
	randomBool = func() bool { return true }
	randIntForEnum = func(n int) int { return n }
	randIndexForEnum = func(n int) int { return n }
}

func TestEmbedValues(t *testing.T) {
	target := &testpb.TestMessage{}

	expected := &testpb.TestMessage{
		SomeInt32:    10,
		SomeSint32:   10,
		SomeUint32:   11,
		SomeFloat32:  10.1,
		SomeFixed32:  11,
		SomeSfixed32: 10,
		SomeInt64:    20,
		SomeSint64:   20,
		SomeUint64:   21,
		SomeFloat64:  10.1,
		SomeFixed64:  21,
		SomeSfixed64: 20,
		SomeStr:      "Gopher",
		SomeBool:     true,
		SomeMsg: &testpb.ChildMessage{
			SomeInt: 10,
		},
		SomeSlice: []string{"Gopher"},
		SomeMsgs: []*testpb.ChildMessage{
			{SomeInt: 10},
		},
		SomeMap: map[int32]*testpb.ChildMessage{
			10: {
				SomeInt: 10,
			},
		},
		SomeEnum:  testpb.SomeEnum_SOME_ENUM_VALUE_2,
		SomeEnum2: 0,
		SomeOneOf: &testpb.TestMessage_OneOfStr{
			OneOfStr: "Gopher",
		},
	}

	err := EmbedValues(target)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(expected, target, protocmp.Transform()); diff != "" {
		t.Errorf("response didn't match (-want / +got)\n%s", diff)
	}
}
