package main

import (
	"fmt"

	"github.com/sryoya/protorand"
	testpb "github.com/sryoya/protorand/testdata"
)

func main() {
	pr := protorand.New()
	pb := &testpb.TestMessage{}
	pr.EmbedValues(pb)
	fmt.Println(pb)
}
