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
	"math/big"
)
/*
 */
type CborIO interface {

	Write(io.Writer) (uint64, error)

	Read(io.Reader) (error)
}
/*
 */
type CborSignedInteger = int64
type CborTag = uint8
type CborFloat = float64
type CborBytes = []byte
type CborChars = string
type CborArray = []any
type CborMap = map[any]any
type CborObject struct {
	tag CborTag
	dat any
}
const CborTrue bool = true
const CborFalse bool = false
const CborNull byte = 0
const CborUndefined byte = 0

type CborBigInt = big.Int
type CborBigFloat = big.Float
