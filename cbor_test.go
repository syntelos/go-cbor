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

	source string

	target []byte
}

var TypeTestCoderObject TypeTestCoder = TypeTestCoder{source: TestStringDatum, target: TestStringCode}

func (this TypeTestCoder) Equals(that TypeTestCoder) (bool) {
	if this.source == that.source {
		var m, n, o int = 0, len(this.target), len(that.target)
		if n == o {
			for ; n < m; n++ {
				if this.target[n] != that.target[n] {
					return false
				}
			}
			return true
		}
	}
	return false
}
func (this TypeTestCoder) Encode() (code Object) {
	var text map[string]any = map[string]any{ "source": this.source, "target": this.target}

	code = Encode(text)
	return code
}
func (this TypeTestCoder) Decode(cbor Object) (TypeTestCoder) {
	this.source = ""
	this.target = nil

	var text map[string]any = cbor.Decode().(map[string]any)

	this.source = text["source"].(string)
	this.target = text["target"].([]byte)

	return this
}

func TestCoder(t *testing.T){
	var text TypeTestCoder = TypeTestCoderObject

	var code Object = text.Encode()

	var check TypeTestCoder = text.Decode(code)
	
	if !TypeTestCoderObject.Equals(check) {

		t.Error("Decoding")
	}
}
