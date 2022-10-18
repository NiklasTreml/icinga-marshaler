package marshaler

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestMarshal(t *testing.T) {
	type Basic struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
	}
	type BasicNested struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
		DeepNested  Basic
	}
	type testStruct struct {
		StringValue string
		IntValue    int32
		BoolValue   bool
		FloatValue  float32
		Nested      BasicNested
	}

	type pointer struct {
		Pointer     *Basic
		StringValue string
	}

	type recursive struct {
		StringValue string
		Recursive   *recursive
	}

	data := []struct {
		args any
		want []byte

		name string
	}{
		{name: "Marshals with pointer", args: pointer{
			StringValue: "Hello",
			Pointer: &Basic{
				StringValue: "PointerString",
				IntValue:    50, BoolValue: true, FloatValue: 50.5},
		},
			want: []byte("Pointer.StringValue=PointerString Pointer.IntValue=50 Pointer.BoolValue=true Pointer.FloatValue=50.5 StringValue=Hello")},
		{name: "Marshals unnested", args: Basic{StringValue: "MyString", IntValue: 50, BoolValue: true, FloatValue: 50.5}, want: []byte("StringValue=MyString IntValue=50 BoolValue=true FloatValue=50.5")},
		{name: "Marshals nested", args: testStruct{StringValue: "MyString", IntValue: 50, BoolValue: true, FloatValue: 5.0, Nested: BasicNested{StringValue: "myNestedString", IntValue: 100, BoolValue: true, FloatValue: 10.5, DeepNested: Basic{StringValue: "myNestedString", IntValue: 100, BoolValue: true, FloatValue: 10.5}}}, want: []byte("StringValue=MyString IntValue=50 BoolValue=true FloatValue=5 Nested.StringValue=myNestedString Nested.IntValue=100 Nested.BoolValue=true Nested.FloatValue=10.5 Nested.DeepNested.StringValue=myNestedString Nested.DeepNested.IntValue=100 Nested.DeepNested.BoolValue=true Nested.DeepNested.FloatValue=10.5")},
		{name: "Marshals empty", args: struct{}{}, want: []byte("")},
		{name: "Marshals nil pointer", args: recursive{StringValue: "Top", Recursive: nil}, want: []byte("StringValue=Top")},
		{name: "Marshals recursive", args: recursive{StringValue: "L1", Recursive: &recursive{StringValue: "L2", Recursive: &recursive{StringValue: "L3", Recursive: nil}}}, want: []byte("StringValue=L1 Recursive.StringValue=L2 Recursive.Recursive.StringValue=L3")},
	}

	for _, tt := range data {
		t.Run(tt.name, func(t *testing.T) {
			result := Marshal(tt.args)
			if !reflect.DeepEqual(result, tt.want) {
				t.Fatalf("\nExpected:\t %v \ngot:\t\t %v", string(tt.want), string(result))
				t.Fatalf("\nExpected:\t %v \ngot:\t\t %v", tt.want, result)
			}

			log.Println(string(result))
		})
	}
}

func ExampleMarshal() {
	type Check struct {
		Status string
		Memory int64
	}

	status := Check{
		Status: "WARN",
		Memory: 1024,
	}

	bytes := Marshal(status)
	// Output: Status=WARN Memory=1024
	fmt.Println(string(bytes))

}
