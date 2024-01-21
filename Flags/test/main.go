package main

import (
	"flag"
	"fmt"
	"reflect"
	"time"
)

type Text struct {
	Content []byte
}

func (t *Text) MarshalText() (text []byte, err error) {
	return text, nil
}
func (t *Text) UnmarshalText(text []byte) error {
	t.Content = text
	return nil
}

func main() {
	var second string
	flag.StringVar(&second, "string", "str", "string usage")
	var boolean bool
	flag.BoolVar(&boolean, "boolean", true, "boolean usage")
	var dur time.Duration
	flag.DurationVar(&dur, "duration", 0, "duration usage")
	var text = Text{
		Content: []byte{1, 2, 3},
	}
	flag.TextVar(&text, "text", &text, "text usage")
	flag.Func("name", "usage", func(s string) error {
		fmt.Println(s)
		return nil
	})
	//json.Marshal()
	flag.VisitAll(func(f *flag.Flag) {
		fmt.Println(reflect.TypeOf(f.Value).String())
		//switch reflect.TypeOf(f.Value) {
		//
		//}
		fmt.Println(f)
	})
	flag.Parse()
	test()
}

type BaseObject interface {
	getName() string
}
type Builder interface {
	BaseObject
	build()
}
type StringBuilder struct{}

func (sb StringBuilder) getName() string {
	return "StringBuilder"
}
func (sb StringBuilder) build() {
	fmt.Println("StringBuilder build")
}

type Cat struct {
	Name string
}

func (cat Cat) getName() string {
	return cat.Name
}

type BuilderCollection[T Builder] struct {
	builder T
}

func (bc BuilderCollection[T]) runBuild() {
	bc.builder.build()
}

func testA() {
	bc1 := BuilderCollection[StringBuilder]{}
	bc2 := BuilderCollection[Cat]{}
	bc1.runBuild()
	bc2.runBuild()
}

type MySlice[T int | float32 | string] struct {
	value []T
}

func (mySlice MySlice[T]) Add(a T, b T) T {
	return a + b
}
func (mySlice MySlice[T]) Sum() T {
	//switch mySlice.value[0].(type) {
	//
	//}
	var sum T
	for _, value := range mySlice.value {
		sum += value
	}
	return sum
}

type MyDataType struct {
}

func (bName MyDataType) getName() string {
	return "MyDataType"
}

type Name interface {
	~string | ~[]byte | interface{ MyDataType }
	getName() string
}

type baseName[T Name] struct {
	name T
}

func (bName baseName[T]) test() string {
	return bName.name.getName()
}

func test() {
	bName := baseName[MyDataType]{}
	fmt.Println(bName.test())
}
