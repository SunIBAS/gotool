package flag

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"testing"
)

func TestFlag(t *testing.T) {
	var first string
	var second string
	first = *flag.String("first", "vFirst", "usage first")
	flag.StringVar(&second, "second", "vSecond", "Usage second")
	fmt.Println(first)
	flag.Visit(func(f *flag.Flag) {
		fmt.Println(f)
	})
	a := os.File{}
	json.Marshal(a)
	flag.Parse()
}
