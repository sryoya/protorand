package main

import (
	"fmt"

	"github.com/sryoya/protorand"
	testpb "github.com/sryoya/protorand/testdata"
)

func main() {
	pr := protorand.New()
	pb := &testpb.TestMessage{} // the base type to generate rand values

	out1, err := pr.Gen(pb)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out1)

	out2, err := pr.Gen(pb)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(out2)
}
