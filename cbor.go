/*
 * CBOR I/O
 * Copyright 2023 John Douglas Pritchard, Syntelos
 *
 *
 * References
 *
 * https://tools.ietf.org/html/rfc8949
 */
package cbor

import (
	"io"
)
/*
 */
type CborObject interface {

	Write(io.Writer) (uint64, error)

	Read(io.Reader) (error)
}
