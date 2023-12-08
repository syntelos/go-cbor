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
	"fmt"
	"testing"
)

func TestTag(t *testing.T){
	var o Object = Object{0b01011000}

	fmt.Printf("0x%02X 0b%08b [%s] \"%s\"\n", o.Tag(), o.Tag(), o.MajorString(), o.String())
}
func TestText(t *testing.T){
	var o Object = Object{0x6D,0x68,0x65,0x6C,0x6C,0x6F,0x2C,0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E}

	fmt.Printf("[%s] \"%s\"\n", o.MajorString(), o.Text())
}
