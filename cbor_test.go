/*
 * CBOR Test
 * Copyright 2023 John Douglas Pritchard, Syntelos
 *
 *
 * References
 *
 * https://datatracker.ietf.org/doc/html/rfc8949
 */
package cbor

import (
	"encoding/json"
	"fmt"
	"testing"
)

const TestStringDatum string = "hello, world."

func TestString(t *testing.T){
	var s string = TestStringDatum
	var o Object = Encode(s)

	if MajorText == o.Major() {

		if TestStringDatum == o.Text() {
			fmt.Printf("[%s] \"%s\".\n", o.MajorString(),o.Text())
		} else {
			t.Errorf("Expected test vector '%s', found '%s'.",TestStringDatum,o.Text())
		}
	} else {
		t.Errorf("Expected major type [text], found '%s'.",o.MajorString())
	}

}

type TypeTestCoder struct {

	name string
	count int
	data []byte
}

func (this TypeTestCoder) Encode() (code Object) {
	var text map[string]any = map[string]any{ "name": this.name, "count": this.count, "data": this.data}

	code = Encode(text)
	return code
}
func (this TypeTestCoder) Decode(cbor Object) (TypeTestCoder) {
	this.name = ""
	this.count = 0
	this.data = nil

	var text map[string]any = cbor.Decode().(map[string]any) // [TODO] [BUG]

	var by []byte
	var er error
	by, er = json.Marshal(text)
	if nil == er {
		this.name = text["name"].(string)
		this.count = text["count"].(int)
		this.data = text["data"].([]byte)

		fmt.Println(string(by))
	}
	return this
}

func TestDescribe(t *testing.T){
	var text TypeTestCoder = TypeTestCoder{name: TestStringDatum, count: 13, data: []byte{0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}}

	var code Object = text.Encode()

	var structure string = code.Describe()

	fmt.Println(structure)
}

func TestObject(t *testing.T){
	var text TypeTestCoder = TypeTestCoder{name: TestStringDatum, count: 13, data: []byte{0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}}

	var code Object = text.Encode()

	var check TypeTestCoder = text.Decode(code) // [TODO] [BREAKPOINT]
	
	if 0 == len(check.name) || 0 == check.count || 0 == len(check.data) {

		t.Error("Empty result of decoding.")
	} else {
		fmt.Println(json.Marshal(check))
	}
}
