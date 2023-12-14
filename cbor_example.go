/*
 * CBOR Examples
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
)

func ExampleTags(){
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
