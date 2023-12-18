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

var TestStringCode []byte = []byte{0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}

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

var TypeTestCoderObject TypeTestCoder = TypeTestCoder{name: TestStringDatum, count: len(TestStringDatum), data: TestStringCode}

func (this TypeTestCoder) Equals(that TypeTestCoder) (bool) {
	if this.name == that.name && this.count == that.count {
		var m, n, o int = 0, len(this.data), len(that.data)
		if n == o {
			for ; n < m; n++ {
				if this.data[n] != that.data[n] {
					return false
				}
			}
			return true
		}
	}
	return false
}
func (this TypeTestCoder) String() (string) {
	var by []byte
	var er error
	by, er = json.Marshal(this)
	if nil == er {
		return string(by)
	} else {
		return ""
	}
}
func (this TypeTestCoder) Encode() (code Object) {
	var text map[string]any = map[string]any{ "name": this.name, "count": this.count, "data": this.data}

	code = Encode(text) // [TODO] [BREAKPOINT]
	return code
}
func (this TypeTestCoder) Decode(cbor Object) (TypeTestCoder) {
	this.name = ""
	this.count = 0
	this.data = nil

	var text map[string]any = cbor.Decode().(map[string]any) // [TODO] [BUG]

	this.name = text["name"].(string)
	this.count = text["count"].(int)
	this.data = text["data"].([]byte)

	return this
}

func TestDescribe(t *testing.T){
	var text TypeTestCoder = TypeTestCoderObject

	var code Object = text.Encode()

	var structure string = code.Describe()

	fmt.Println(structure)
}

func TestObject(t *testing.T){
	var text TypeTestCoder = TypeTestCoderObject

	var code Object = text.Encode()

	var check TypeTestCoder = text.Decode(code)
	
	if !TypeTestCoderObject.Equals(check) {

		t.Error("Result of decoding.")
	} else {
		fmt.Println(json.Marshal(check))
	}
}
