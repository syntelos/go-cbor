/*
 * GOPL type struct object coding requires type binding.
 */
package cbor

import (
	"encoding/json"
	"fmt"
)

type TypeExampleCoder struct {

	source string

	target []byte
}

const ExampleStringDatum string = "hello, world."

var ExampleStringCode []byte = []byte{0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}

var TypeExampleCoderObject TypeExampleCoder = TypeExampleCoder{source: ExampleStringDatum, target: ExampleStringCode}

func (this TypeExampleCoder) Equals(that TypeExampleCoder) (bool) {
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
func (this TypeExampleCoder) String() (string) {
	var by []byte
	var er error
	by, er = json.Marshal(this)
	if nil == er {
		return string(by)
	} else {
		return ""
	}
}
func (this TypeExampleCoder) Encode() (code Object) {
	var text map[string]any = map[string]any{ "source": this.source, "target": this.target}

	code = Encode(text)
	return code
}
func (this TypeExampleCoder) Decode(cbor Object) (TypeExampleCoder) {
	this.source = ""
	this.target = nil

	var text map[string]any = cbor.Decode().(map[string]any)

	this.source = text["source"].(string)
	this.target = text["target"].([]byte)

	return this
}
func ExampleDescribe(){
	var text TypeExampleCoder = TypeExampleCoderObject

	var code Object = text.Encode()

	var content string = code.String()

	var encoding string = code.Describe()

	fmt.Printf("%s\t%s\n",content,encoding)
}
