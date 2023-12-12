/*
 * CBOR I/O
 * Copyright 2023 John Douglas Pritchard, Syntelos
 *
 *
 * References
 *
 * https://datatracker.ietf.org/doc/html/rfc8949
 */
package cbor

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

func _TestConstructor(t *testing.T){
	var o Object = Object{0x6D,0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}

	fmt.Printf("[%s] \"%s\"\n", o.MajorString(), o.Text())
}

func _TestEncoder(t *testing.T){
	var s string = "hello, world"
	var o Object = Encode(s)

	if MajorText == o.Major() {
		fmt.Printf("[%s] \"%s\".\n", o.MajorString(),o.Text())
	} else {
		t.Errorf("Expected major type [text], found '%s'.",o.MajorString())
	}

}

func _TestTags(t *testing.T){
	var list []Object = []Object{
		Object{0x00}, Object{0x18}, Object{0x20},
		Object{0x38}, Object{0x40}, Object{0x58},
		Object{0x60}, Object{0x78}, Object{0x80},
		Object{0x98}, Object{0xA0}, Object{0xB8},
		Object{0xC0}, Object{0xC1}, Object{0xC2},
		Object{0xC3}, Object{0xC4}, Object{0xC5},
		Object{0xC6}, Object{0xD5}, Object{0xDB},
		Object{0xE0}, Object{0xF4}, Object{0xF5},
		Object{0xF6}, Object{0xF7}, Object{0xF8},
		Object{0xF9}, Object{0xFA}, Object{0xFB},
		Object{0xFF},
	}

	for _, o := range list {

		fmt.Printf("0x%02X 0b%08b [%s] \"%s\"\n", o.Tag(), o.Tag(), o.MajorString(), o.String())
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
func (this TypeTestCoder) Decode(cbor Object){

	var text map[string]any = cbor.Decode().(map[string]any) // [TODO] BUG

	var by []byte
	var er error
	by, er = json.Marshal(text)
	if nil == er {
		this.name = text["name"].(string)
		this.count = text["count"].(int)
		this.data = text["data"].([]byte)

		fmt.Println(string(by))
	}
}

func TestCoderEncode(t *testing.T){
	var text TypeTestCoder = TypeTestCoder{name: "hello, world", count: 13, data: []byte{0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}}

	var code Object = text.Encode()

	if 0 == len(code) {

		t.Error("Encoded map to empty code.")
	} else {

		fmt.Println(hex.EncodeToString(code))
	}
}

func TestCoderDecode(t *testing.T){
	var code Object = Object{0x64,0x64,0x61,0x74,0x61,0x4d,0x68,0x65,0x6c,0x6c,0x6f,0x2c,0x20,0x77,0x6f,0x72,0x6c,0x64,0x2e,0x64,0x6e,0x61,0x6d,0x65,0x6c,0x68,0x65,0x6c,0x6c,0x6f,0x2c,0x20,0x77,0x6f,0x72,0x6c,0x64}
	var text TypeTestCoder
	text.Decode(code)
	if 0 == len(text.name) || 0 == text.count || 0 == len(text.data) {

		t.Error("Empty result of decoding.")
	} else {
		fmt.Println(json.Marshal(text))
	}
}
