/*
 * CBOR Text Encode
 * Copyright 2023 John Douglas Pritchard, Syntelos
 */
package main

import (
	"fmt"
	"github.com/syntelos/go-cbor"
)

func main(){
	var object cbor.Object

	object.Encode("hello, world.")

	fmt.Println(object.Text())
}
