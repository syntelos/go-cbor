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

func TestCbor(t *testing.T){
	var b CborTag = []byte{0b01011000}

	fmt.Println(b.String())
}
