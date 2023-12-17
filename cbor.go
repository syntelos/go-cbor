/*
 * CBOR RFC8949 I/O
 * Copyright 2023 John Douglas Pritchard, Syntelos
 *
 *
 * References
 *
 * https://tools.ietf.org/html/rfc8949
 */
package cbor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"github.com/syntelos/go-endian"
	"math"
	"math/big"
	"reflect"
)
/*
 * Encoded data set content object.
 */
type Object []byte
/*
 * Content object user interface.  CBOR binary code is
 * consumed and produced in these interfaces.
 */
type IO interface {
	/*
	 * The CBOR producer is replicating.
	 */
	Write(io.Writer) (error)
	/*
	 * The CBOR consumer is validating.
	 */
	Read(io.Reader) (Object, error)
}
/*
 * Eight bits of Tag.  See Appendix B Table 7 [RFC8949].
 * See also ./doc/cbor-rfc8949-table.go
 */
type Tag byte
/*
 * High three bits of Tag shifted onto Major Type (0-7).
 * See Section 3.1 [RFC8949].
 */
type Major byte
/*
 * MajorBits = (0b111 << 5)
 */
var MajorUint Major   = Major(0)
var MajorSint Major   = Major(1)
var MajorBlob Major   = Major(2)
var MajorText Major   = Major(3)
var MajorArray Major  = Major(4)
var MajorMap Major    = Major(5)
var MajorTagged Major = Major(6)
var MajorSimple Major = Major(7)
/*
 * A package external struct type can extend this package by
 * implementing this interface.
 */
type Coder interface {
	/*
	 * The member type category produces a CBOR Object
	 * may perform byte encoding by calling
	 * "cbor.Encode" on a GOPL primitive, i.e. "map".
	 */
	Encode() (Object)
	/*
	 * The member type category consumes a CBOR binary
	 * by calling "cbor.Decode" to yield a GOPL
	 * primitive, i.e. "map".
	 */
	Decode(Object) (any)
}
/*
 * Internal CBOR Break
 */
var Break error = errors.New("CBOR Break")
/*
 * Validation errors produced by <Object#Read>.
 */
const ErrorWrapRead string = "CBOR Data: %w"
var ErrorUnrecognizedTag error = errors.New("Unrecognized CBOR Tag")
var ErrorMissingData error = errors.New("Missing CBOR Data")
/*
 */
func (this Object) Write(w io.Writer) (e error){
	_, e = w.Write(this)
	return e
}
/*
 */
func (this Object) Read(r io.Reader) (Object, error){
	var tag []byte = make([]byte,1)
	var m, n int
	var e error

	n, e = r.Read(tag)
	if nil != e {
		return nil, e
	} else if 1 != n {
		return nil, fmt.Errorf("Read (%d) expected (1).",n)
	} else {
		var d []byte
		var t byte = tag[0]
		var a, b Object

		switch t {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
			/* unsigned integer 0x00..0x17 (0..23)
			 */
			this = tag
			return this, nil

		case 0x18:
			/* unsigned integer (one-byte uint8_t follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x19:
			/* unsigned integer (two-byte uint16_t follows)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x1A:
			/* unsigned integer (four-byte uint32_t follows)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x1B:
			/* unsigned integer (eight-byte uint64_t follows)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
			/* negative integer -1-0x00..-1-0x17 (-1..-24)
			 */
			this = tag
			return this, nil

		case 0x38:
			/* negative integer -1-n (one-byte uint8_t for n follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x39:
			/* negative integer -1-n (two-byte uint16_t for n follows)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x3A:
			/* negative integer -1-n (four-byte uint32_t for n follows)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint32(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x3B:
			/* negative integer -1-n (eight-byte uint64_t for n follows)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint64(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
			/* byte string (0x00..0x17 bytes follow)
			 */
			this = tag
			m = int(t-0x40)
			d = make([]byte,m)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if m != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x58:
			/* byte string (one-byte uint8_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x59:
			/* byte string (two-byte uint16_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x5A:
			/* byte string (four-byte uint32_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint32(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x5B:
			/* byte string (eight-byte uint64_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint64(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x5F:
			/* byte string, byte strings follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = Object{}
				a, e = a.Read(r)
				if nil == e {
					this = this.Concatenate(a)
				} else if Break == e {
					e = nil
					break
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}
			return this, nil

		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			/* UTF-8 string (0x00..0x17 bytes follow)
			 */
			this = tag
			m = int(t-0x60)
			d = make([]byte,m)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if m != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0x78:
			/* UTF-8 string (one-byte uint8_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x79:
			/* UTF-8 string (two-byte uint16_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != n {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(p)
					return this, nil
				}	
			}

		case 0x7A:
			/* UTF-8 string (four-byte uint32_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint32(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(d)
					return this, nil
				}
			}

		case 0x7B:
			/* UTF-8 string (eight-byte uint64_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else if z != uint64(n) {
					return nil, ErrorMissingData
				} else {
					this = this.Concatenate(d)
					return this, nil
				}	
			}

		case 0x7F:
			/* UTF-8 string, UTF-8 strings follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = Object{}
				a, e = a.Read(r)
				if nil == e {
					this = this.Concatenate(a)
				} else if Break == e {
					e = nil
					break
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}
			return this, nil

		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			/* array (0x00..0x17 data items follow)
			 */
			this = tag
			m = int(t-0x80)
			for n = 0; n < m; n++ {
				a = Object{}
				a, e = a.Read(r)
				if nil == e {
					this = this.Concatenate(a)
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}
			return this, nil

		case 0x98:
			/* array (one-byte uint8_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var z int = int(d[0])
				for n = 0; n < z; n++ {
					a = Object{}
					a, e = a.Read(r)
					if nil == e {
						this = this.Concatenate(a)
					} else {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					}
				}
				return this, nil
			}

		case 0x99:
			/* array (two-byte uint16_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint16 = 0, endian.BigEndian.DecodeUint16(d)
				for ; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil == e {
						this = this.Concatenate(a)
					} else {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					}
				}
				return this, nil
			}

		case 0x9A:
			/* array (four-byte uint32_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint32 = 0, endian.BigEndian.DecodeUint32(d)
				for ; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil == e {
						this = this.Concatenate(a)
					} else {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					}
				}
				return this, nil
			}

		case 0x9B:
			/* array (eight-byte uint64_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint64 = 0, endian.BigEndian.DecodeUint64(d)
				for ; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil == e {
						this = this.Concatenate(a)
					} else {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					}
				}
				return this, nil
			}

		case 0x9F:
			/* array, data items follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = Object{}
				a, e = a.Read(r)
				if nil == e {
					this = this.Concatenate(a)
				} else if Break == e {
					e = nil
					break
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}
			return this, nil

		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			/* map (0x00..0x17 pairs of data items follow)
			 */
			this = tag
			m, n = 0, int(t-0xA0)
			for ; m < n; m++ {
				a = Object{}
				a, e = a.Read(r)
				if nil != e {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				} else {
					this = this.Concatenate(a)
					b = make([]byte,0)
					b, e = b.Read(r)
					if nil != e {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					} else {
						this = this.Concatenate(b)
					}	
				}
			}
			return this, nil

		case 0xB8:
			/* map (one-byte uint8_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint8 = 0, uint8(d[0])
				for x = 0; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil != e {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					} else {
						this = this.Concatenate(a)
						b = make([]byte,0)
						b, e = b.Read(r)
						if nil != e {
							return nil, fmt.Errorf(ErrorWrapRead,e)
						} else {
							this = this.Concatenate(b)
						}	
					}
				}
				return this, nil
			}

		case 0xB9:
			/* map (two-byte uint16_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint16 = 0, endian.BigEndian.DecodeUint16(d)
				for x = 0; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil != e {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					} else {
						this = this.Concatenate(a)
						b = make([]byte,0)
						b, e = b.Read(r)
						if nil != e {
							return nil, fmt.Errorf(ErrorWrapRead,e)
						} else {
							this = this.Concatenate(b)
						}	
					}
				}
				return this, nil
			}

		case 0xBA:
			/* map (four-byte uint32_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint32 = 0, endian.BigEndian.DecodeUint32(d)
				for x = 0; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil != e {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					} else {
						this = this.Concatenate(a)
						b = make([]byte,0)
						b, e = b.Read(r)
						if nil != e {
							return nil, fmt.Errorf(ErrorWrapRead,e)
						} else {
							this = this.Concatenate(b)
						}	
					}
				}
				return this, nil
			}

		case 0xBB:
			/* map (eight-byte uint64_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				var x, z uint64 = 0, endian.BigEndian.DecodeUint64(d)
				for x = 0; x < z; x++ {
					a = Object{}
					a, e = a.Read(r)
					if nil != e {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					} else {
						this = this.Concatenate(a)
						b = make([]byte,0)
						b, e = b.Read(r)
						if nil != e {
							return nil, fmt.Errorf(ErrorWrapRead,e)
						} else {
							this = this.Concatenate(b)
						}	
					}
				}
				return this, nil
			}

		case 0xBF:
			/* map, pairs of data items follow, terminated by 'break'
			 */
			this = tag

			for nil == e {
				a = Object{}
				a, e = a.Read(r)
				if nil == e {
					this = this.Concatenate(a)

					b = make([]byte,0)
					b, e = b.Read(r)
					if nil == e {
						this = this.Concatenate(b)
					} else {
						return nil, fmt.Errorf(ErrorWrapRead,e)
					}
				} else if Break == e {
					e = nil
					break
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}
			return this, nil

		case 0xC0, 0xC1:
			/* date/time (data item follows; see Section 3.4.1 and 3.4.2)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xC2:
			/* unsigned bignum (data item 'byte string' follows)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xC3:
			/* negative bignum (data item 'byte string' follows)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xC4:
			/* decimal Fraction (data item 'array' follows; see Section 3.4.4)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xC5:
			/* bigfloat (data item 'array' follows; see Section 3.4.4)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4:
			/* (tag)
			 */
			this = tag
			return this, nil

		case 0xD5, 0xD6, 0xD7:
			/* expected conversion (data item follows; see Section 3.4.5.2)
			 */
			this = tag
			a = Object{}
			a, e = a.Read(r)
			if nil == e {
				this = this.Concatenate(a)
				return this, nil
			} else {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			}

		case 0xD8:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,1)
			n, e = r.Read(a)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, fmt.Errorf("Data expected (1) found (%d).",n)
			} else {
				this = this.Concatenate(a)
				b = make([]byte,0)
				b, e = b.Read(r)
				if nil == e {
					this = this.Concatenate(b)
					return this, nil
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}

		case 0xD9:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,2)
			n, e = r.Read(a)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, fmt.Errorf("Data expected (2) found (%d).",n)
			} else {
				this = this.Concatenate(a)
				b = make([]byte,0)
				b, e = b.Read(r)
				if nil == e {
					this = this.Concatenate(b)
					return this, nil
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}

		case 0xDA:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,4)
			n, e = r.Read(a)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, fmt.Errorf("Data expected (4) found (%d).",n)
			} else {
				this = this.Concatenate(a)
				b = make([]byte,0)
				b, e = b.Read(r)
				if nil == e {
					this = this.Concatenate(b)
					return this, nil
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}

		case 0xDB:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,8)
			n, e = r.Read(a)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, fmt.Errorf("Data expected (8) found (%d).",n)
			} else {
				this = this.Concatenate(a)
				b = make([]byte,0)
				b, e = b.Read(r)
				if nil == e {
					this = this.Concatenate(b)
					return this, nil
				} else {
					return nil, fmt.Errorf(ErrorWrapRead,e)
				}
			}

		case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3:
			/* (simple value)
			 */
			this = tag
			return this, nil

		case 0xF4:
			/* "false"
			 */
			this = tag
			return this, nil

		case 0xF5:
			/* "true"
			 */
			this = tag
			return this, nil

		case 0xF6:
			/* "null"
			 */
			this = tag
			return this, nil

		case 0xF7:
			/* "undefined"
			 */
			this = tag
			return this, nil

		case 0xF8:
			/* (simple value, one byte follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 1 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0xF9:
			/* half-precision float (two-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 2 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0xFA:
			/* single-precision float (four-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 4 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0xFB:
			/* double-precision float (eight-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return nil, fmt.Errorf(ErrorWrapRead,e)
			} else if 8 != n {
				return nil, ErrorMissingData
			} else {
				this = this.Concatenate(d)
				return this, nil
			}

		case 0xFF:
			/* 'break' stop code"
			 */
			this = tag
			return nil, Break

		default:
			return nil, ErrorUnrecognizedTag
		}
	}
}
/*
 */
func (this *Object) String() string {
	if this.HasTag() {
		var tag Tag = this.Tag()
		switch tag {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
			return "unsigned integer 0x00..0x17 (0..23)"
		case 0x18:
			return "unsigned integer (one-byte uint8_t follows)"
		case 0x19:
			return "unsigned integer (two-byte uint16_t follows)"
		case 0x1A:
			return "unsigned integer (four-byte uint32_t follows)"
		case 0x1B:
			return "unsigned integer (eight-byte uint64_t follows)"
		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
			return "negative integer -1-0x00..-1-0x17 (-1..-24)"
		case 0x38:
			return "negative integer -1-n (one-byte uint8_t for n follows)"
		case 0x39:
			return "negative integer -1-n (two-byte uint16_t for n follows)"
		case 0x3A:
			return "negative integer -1-n (four-byte uint32_t for n follows)"
		case 0x3B:
			return "negative integer -1-n (eight-byte uint64_t for n follows)"
		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
			return "byte string (0x00..0x17 bytes follow)"
		case 0x58:
			return "byte string (one-byte uint8_t for n, and then n bytes follow)"
		case 0x59:
			return "byte string (two-byte uint16_t for n, and then n bytes follow)"
		case 0x5A:
			return "byte string (four-byte uint32_t for n, and then n bytes follow)"
		case 0x5B:
			return "byte string (eight-byte uint64_t for n, and then n bytes follow)"
		case 0x5F:
			return "byte string, byte strings follow, terminated by 'break'"
		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			return "UTF-8 string (0x00..0x17 bytes follow)"
		case 0x78:
			return "UTF-8 string (one-byte uint8_t for n, and then n bytes follow)"
		case 0x79:
			return "UTF-8 string (two-byte uint16_t for n, and then n bytes follow)"
		case 0x7A:
			return "UTF-8 string (four-byte uint32_t for n, and then n bytes follow)"
		case 0x7B:
			return "UTF-8 string (eight-byte uint64_t for n, and then n bytes follow)"
		case 0x7F:
			return "UTF-8 string, UTF-8 strings follow, terminated by 'break'"
		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			return "array (0x00..0x17 data items follow)"
		case 0x98:
			return "array (one-byte uint8_t for n, and then n data items follow)"
		case 0x99:
			return "array (two-byte uint16_t for n, and then n data items follow)"
		case 0x9A:
			return "array (four-byte uint32_t for n, and then n data items follow)"
		case 0x9B:
			return "array (eight-byte uint64_t for n, and then n data items follow)"
		case 0x9F:
			return "array, data items follow, terminated by 'break'"
		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			return "map (0x00..0x17 pairs of data items follow)"
		case 0xB8:
			return "map (one-byte uint8_t for n, and then n pairs of data items follow)"
		case 0xB9:
			return "map (two-byte uint16_t for n, and then n pairs of data items follow)"
		case 0xBA:
			return "map (four-byte uint32_t for n, and then n pairs of data items follow)"
		case 0xBB:
			return "map (eight-byte uint64_t for n, and then n pairs of data items follow)"
		case 0xBF:
			return "map, pairs of data items follow, terminated by 'break'"
		case 0xC0:
			return "text-based date/time (data item follows; see Section 3.4.1)"
		case 0xC1:
			return "epoch-based date/time (data item follows; see Section 3.4.2)"
		case 0xC2:
			return "unsigned bignum (data item 'byte string' follows)"
		case 0xC3:
			return "negative bignum (data item 'byte string' follows)"
		case 0xC4:
			return "decimal Fraction (data item 'array' follows; see Section 3.4.4)"
		case 0xC5:
			return "bigfloat (data item 'array' follows; see Section 3.4.4)"
		case 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4:
			return "(tag) "
		case 0xD5, 0xD6, 0xD7:
			return "expected conversion (data item follows; see Section 3.4.5.2)"
		case 0xD8, 0xD9, 0xDA, 0xDB:
			return "(more tags; 1/2/4/8 bytes of tag number and then a data item follow)"
		case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3:
			return "(simple value)"
		case 0xF4:
			return "false"
		case 0xF5:
			return "true"
		case 0xF6:
			return "null"
		case 0xF7:
			return "undefined"
		case 0xF8:
			return "(simple value, one byte follows)"
		case 0xF9:
			return "half-precision float (two-byte IEEE 754)"
		case 0xFA:
			return "single-precision float (four-byte IEEE 754)"
		case 0xFB:
			return "double-precision float (eight-byte IEEE 754)"
		case 0xFF:
			return "'break' stop code"
		default:
			return ""
		}
	} else {
		return ""
	}
}
/*
 * Resolve tag structure of object.
 */
func (this Object) HasTag() bool {
	var z int = len(this)
	return (0 < z)
}
/*
 * Resolve tag value of object.
 */
func (this Object) Tag() Tag {
	if this.HasTag() {
		return Tag(this[0])
	} else {
		return 0
	}
}
/*
 * Resolve major type from tag.
 */
func (this Object) Major() Major {
	if this.HasTag() {
		var tag byte = byte(this.Tag())
		var major byte = ((tag & 0xE0)>>5)
		return Major(major)
	} else {
		return Major(0)
	}
}
/*
 * Describe major type of tag.
 */
func (this Object) MajorString() string {
	if this.HasTag() {
		switch this.Major() {
		case MajorUint:
			return "unsigned integer"
		case MajorSint:
			return "signed integer"
		case MajorBlob:
			return "blob"
		case MajorText:
			return "text"
		case MajorArray:
			return "array"
		case MajorMap:
			return "map"
		case MajorTagged:
			return "tagged data item"
		default:
			return "float, simple, break"
		}
	} else {
		return ""
	}
}
/*
 * Resolve text object type.
 */
func (this Object) HasText() bool {
	return (this.HasTag() && MajorText == this.Major())
}
/*
 * Resolve text object content.
 */
func (this Object) Text() (s string) {
	var a any = this.Decode()
	if nil != a {
		return a.(string)
	} else {
		return ""
	}
}
/*
 */
func copier(dst []byte, dx, dz int, src []byte, sx, sz int) ([]byte) {
	for dx < dz && sx < sz {

		dst[dx] = src[sx]
		dx += 1
		sx += 1
	}
	return dst
}
func (this Object) Concatenate(b []byte) (Object) {
	var a []byte = this
	var a_len int = len(a)
	if 0 == a_len {
		this = b
	} else {
		var b_len int = len(b)
		if 0 == b_len {
			this = b
		} else {
			var c_len int = (a_len+b_len)

			var c []byte = make([]byte,c_len)

			c = copier(c,0,c_len,a,0,a_len)

			c = copier(c,a_len,c_len,b,0,b_len)

			this = c
		}
	}
	return this
}
/*
 * Define object as major type tag.
 */
func Define(m Major) (this Object) {

	var major byte = ((byte(m) & 7) << 5)

	this = Object{major}

	return this
}
/*
 * Define object as major type tag refined by size (in octet
 * count).
 */
func (this Object) Refine(size uint64) (Object) {

	var major Major = this.Major()
	switch major {
	case MajorUint:
		if 0x17 >= size {
			this[0] = byte(size)
		} else if 0xFF >= size {
			this[0] = 0x18
		} else if 0xFFFF >= size {
			this[0] = 0x19
		} else if 0xFFFFFFFF >= size {
			this[0] = 0x1A
		} else {
			this[0] = 0x1B
		}
		return this

	case MajorSint:
		if 0x17 >= size {
			this[0] = byte(size)+0x20
		} else if 0xFF >= size {
			this[0] = 0x38
		} else if 0xFFFF >= size {
			this[0] = 0x39
		} else if 0xFFFFFFFF >= size {
			this[0] = 0x3A
		} else {
			this[0] = 0x3B
		}
		return this

	case MajorBlob:
		if 0x17 >= size {
			this[0] = byte(size)+0x40
		} else if 0xFF >= size {
			this[0] = 0x58
		} else if 0xFFFF >= size {
			this[0] = 0x59
		} else if 0xFFFFFFFF >= size {
			this[0] = 0x5A
		} else {
			this[0] = 0x5B
		}
		return this

	case MajorText:
		if 0x17 >= size {
			this[0] = byte(size)+0x60
		} else if 0xFF >= size {
			this[0] = 0x78
		} else if 0xFFFF >= size {
			this[0] = 0x79
		} else if 0xFFFFFFFF >= size {
			this[0] = 0x7A
		} else {
			this[0] = 0x7B
		}
		return this

	case MajorArray:
		if 0x17 >= size {
			this[0] = byte(size)+0x80
		} else if 0xFF >= size {
			this[0] = 0x98
		} else if 0xFFFF >= size {
			this[0] = 0x99
		} else if 0xFFFFFFFF >= size {
			this[0] = 0x9A
		} else {
			this[0] = 0x9B
		}
		return this

	case MajorMap:
		if 0x17 >= size {
			this[0] = byte(size)+0xA0
		} else if 0xFF >= size {
			this[0] = 0xB8
		} else if 0xFFFF >= size {
			this[0] = 0xB9
		} else if 0xFFFFFFFF >= size {
			this[0] = 0xBA
		} else {
			this[0] = 0xBB
		}
		return this
	}
	return this
}
/*
 * Define object content.
 */
func Encode(a any) (this Object) {
	if nil != a {
		switch a.(type) {

		case uint8: // (eq byte)
			this = Define(MajorUint).Refine(1)
			var hbo []byte = []byte{a.(byte)}

			this = this.Concatenate(hbo)
		case uint16:
			this = Define(MajorUint).Refine(2)
			var hbo []byte = endian.BigEndian.EncodeUint16(a.(uint16))
			this = this.Concatenate(hbo)
		case uint32:
			this = Define(MajorUint).Refine(4)
			var hbo []byte = endian.BigEndian.EncodeUint32(a.(uint32))
			this = this.Concatenate(hbo)
		case uint64:
			this = Define(MajorUint).Refine(8)
			var hbo []byte = endian.BigEndian.EncodeUint64(a.(uint64))
			this = this.Concatenate(hbo)

		case int8:
			this = Define(MajorSint).Refine(1)
			var hbo []byte = []byte{a.(byte)}
			this = this.Concatenate(hbo)
		case int16:
			this = Define(MajorSint).Refine(2)
			var hbo []byte = endian.BigEndian.EncodeUint16(a.(uint16))
			this = this.Concatenate(hbo)
		case int32:
			this = Define(MajorSint).Refine(4)
			var hbo []byte = endian.BigEndian.EncodeUint32(a.(uint32))
			this = this.Concatenate(hbo)
		case int64:
			this = Define(MajorSint).Refine(8)
			var hbo []byte = endian.BigEndian.EncodeUint64(a.(uint64))
			this = this.Concatenate(hbo)

		case int:
			var val int = a.(int)
			var typ reflect.Type = reflect.TypeOf(a)
			var siz uint64 = uint64(typ.Size())
			switch siz {
			case 4:
				this = Define(MajorSint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint32(uint32(val))
				this = this.Concatenate(hbo)
			case 8:
				this = Define(MajorSint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint64(uint64(val))
				this = this.Concatenate(hbo)
			}

		case uint:
			var val uint = a.(uint)
			var typ reflect.Type = reflect.TypeOf(a)
			var siz uint64 = uint64(typ.Size())
			switch siz {
			case 4:
				this = Define(MajorUint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint32(uint32(val))
				this = this.Concatenate(hbo)
			case 8:
				this = Define(MajorUint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint64(uint64(val))
				this = this.Concatenate(hbo)
			}

		case uintptr:
			var val uintptr = a.(uintptr)
			var typ reflect.Type = reflect.TypeOf(a)
			var siz uint64 = uint64(typ.Size())
			switch siz {
			case 4:
				this = Define(MajorUint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint32(uint32(val))
				this = this.Concatenate(hbo)
			case 8:
				this = Define(MajorUint).Refine(siz)
				var hbo []byte = endian.BigEndian.EncodeUint64(uint64(val))
				this = this.Concatenate(hbo)
			}


		case []byte:
			this = Define(MajorBlob)
			var bry []byte = a.([]byte)
			var brz uint64 = uint64(len(bry))
			this = this.Refine(brz)
			switch this.Tag() {
			case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
				this = this.Concatenate(bry)
			case 0x58:
				var cnt uint8 = uint8(brz)
				var brc []byte = []byte{cnt}
				this = this.Concatenate(brc)
				this = this.Concatenate(bry)
			case 0x59:
				var cnt uint16 = uint16(brz)
				var brc []byte = endian.BigEndian.EncodeUint16(cnt)
				this = this.Concatenate(brc)
				this = this.Concatenate(bry)
			case 0x5A:
				var cnt uint32 = uint32(brz)
				var brc []byte = endian.BigEndian.EncodeUint32(cnt)
				this = this.Concatenate(brc)
				this = this.Concatenate(bry)
			case 0x5B:
				var cnt uint64 = brz
				var brc []byte = endian.BigEndian.EncodeUint64(cnt)
				this = this.Concatenate(brc)
				this = this.Concatenate(bry)
			}


		case string:
			this = Define(MajorText)
			var str string = a.(string)
			var sty []byte = []byte(str)
			var stz uint64 = uint64(len(sty))
			this = this.Refine(stz)
			switch this.Tag() {
			case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
				this = this.Concatenate(sty)
			case 0x78:
				var cnt uint8 = uint8(stz)
				var stc []byte = []byte{cnt}
				this = this.Concatenate(stc)
				this = this.Concatenate(sty)
			case 0x79:
				var cnt uint16 = uint16(stz)
				var stc []byte = endian.BigEndian.EncodeUint16(cnt)
				this = this.Concatenate(stc)
				this = this.Concatenate(sty)
			case 0x7A:
				var cnt uint32 = uint32(stz)
				var stc []byte = endian.BigEndian.EncodeUint32(cnt)
				this = this.Concatenate(stc)
				this = this.Concatenate(sty)
			case 0x7B:
				var cnt uint64 = stz
				var stc []byte = endian.BigEndian.EncodeUint64(cnt)
				this = this.Concatenate(stc)
				this = this.Concatenate(sty)
			}

		case []any:
			this = Define(MajorArray)
			var ary []any = a.([]any)
			var arz uint64 = uint64(len(ary))
			this = this.Refine(arz)
			switch this.Tag() {
			case 0x98:
				var cnt uint8 = uint8(arz)
				var arc []byte = []byte{cnt}
				this = this.Concatenate(arc)
			case 0x99:
				var cnt uint16 = uint16(arz)
				var arc []byte = endian.BigEndian.EncodeUint16(cnt)
				this = this.Concatenate(arc)
			case 0x9A:
				var cnt uint32 = uint32(arz)
				var arc []byte = endian.BigEndian.EncodeUint32(cnt)
				this = this.Concatenate(arc)
			case 0x9B:
				var cnt uint64 = uint64(arz)
				var arc []byte = endian.BigEndian.EncodeUint64(cnt)
				this = this.Concatenate(arc)
			}
			for _, v := range ary {
				var vo Object = Encode(v)
				this = this.Concatenate([]byte(vo))
			}

		case map[string]any:
			this = Define(MajorMap)
			var mmm map[string]any = a.(map[string]any)
			var mmz uint64 = uint64(len(mmm))
			this = this.Refine(mmz)
			switch this.Tag() {
			case 0xB8:
				var cnt uint8 = uint8(mmz)
				var mmc []byte = []byte{cnt}
				this = this.Concatenate(mmc)
			case 0xB9:
				var cnt uint16 = uint16(mmz)
				var mmc []byte = endian.BigEndian.EncodeUint16(cnt)
				this = this.Concatenate(mmc)
			case 0xBA:
				var cnt uint32 = uint32(mmz)
				var mmc []byte = endian.BigEndian.EncodeUint32(cnt)
				this = this.Concatenate(mmc)
			case 0xBB:
				var cnt uint64 = uint64(mmz)
				var mmc []byte = endian.BigEndian.EncodeUint64(cnt)
				this = this.Concatenate(mmc)
			}
			for k, v := range mmm {
				var ko Object = Encode(k)
				this = this.Concatenate([]byte(ko))

				var vo Object = Encode(v)
				this = this.Concatenate([]byte(vo))
			}

		case Coder:
			var coder Coder = a.(Coder)
			this = coder.Encode()

		default:
			var undefined Object = Object{0xF7}
			this = undefined
		}
	} else {
		var null Object = Object{0xF6}
		this = null
	}
	return this
}
/*
 * Resolve object content.
 */
func (this Object) Decode() (a any) {
	if this.HasTag() {
		var tag Tag = this.Tag()
		switch tag {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
			return uint8(tag)
		case 0x18:
			var cnt uint8 = this[1]
			var text []byte = this[2:(2+cnt)]
			switch cnt {
			case 2:
				return endian.BigEndian.DecodeUint16(text)
			case 4:
				return endian.BigEndian.DecodeUint32(text)
			case 8:
				return endian.BigEndian.DecodeUint64(text)
			default:
				var value big.Int
				value.SetBytes(text)
				return value
			}
		case 0x19:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			var text []byte = this[3:(3+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x1A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			var text []byte = this[5:(5+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x1B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			var text []byte = this[9:(9+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
			var delta int = (int(tag)-0x20)
			return (-1-delta)
		case 0x38:
			var cnt uint8 = this[1]
			var text []byte = this[2:(2+cnt)]
			switch cnt {
			case 2:
				var value uint16 = endian.BigEndian.DecodeUint16(text)
				return int16(value)
			case 4:
				var value uint32 = endian.BigEndian.DecodeUint32(text)
				return int32(value)
			case 8:
				var value uint64 = endian.BigEndian.DecodeUint64(text)
				return int64(value)
			default:
				var value big.Int
				value.SetBytes(text)
				return value
			}
		case 0x39:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			var text []byte = this[3:(3+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x3A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			var text []byte = this[5:(5+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x3B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			var text []byte = this[9:(9+cnt)]
			var value big.Int
			value.SetBytes(text)
			return value
		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
			var m int = int(tag-0x40)
			var text []byte = this[1:(m+1)]
			return text
		case 0x58:
			var cnt uint8 = this[1]
			var text []byte = this[2:(3+cnt)]
			return text
		case 0x59:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			var text []byte = this[3:(3+cnt)]
			return text
		case 0x5A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			var text []byte = this[5:(5+cnt)]
			return text
		case 0x5B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			var text []byte = this[9:(9+cnt)]
			return text
		case 0x5F:
			var bary Object
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			for true {
				var o Object = Object{}
				var e error
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a = o.Decode()
					if nil != a {
						var src []byte = a.([]byte)
						bary.Concatenate(src)
					}
				}
			}
			return bary
		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			var m int = int(tag-0x60)
			var text []byte = this[1:(m+1)]
			return string(text)
		case 0x78:
			var cnt uint8 = this[1]
			var text []byte = this[2:(2+cnt)]
			return string(text)
		case 0x79:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			var text []byte = this[3:(3+cnt)]
			return string(text)
		case 0x7A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			var text []byte = this[5:(5+cnt)]
			return string(text)
		case 0x7B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			var text []byte = this[9:(9+cnt)]
			return string(text)
		case 0x7F:
			var bary Object
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			for true {
				var o Object = Object{}
				var e error
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a = o.Decode()
					if nil != a {
						var src []byte = a.([]byte)
						bary.Concatenate(src)
					}
				}
			}
			return string(bary)
		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			var m, n int = int(tag-0x80), 0
			var a []any = make([]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a[n] = o.Decode()
				}
			}
			return a
		case 0x98:
			var m, n uint8 = uint8(this[1]), 0
			var a []any = make([]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[2:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a[n] = o.Decode()
				}
			}
			return a
		case 0x99:
			var m, n uint16 = endian.BigEndian.DecodeUint16(this[1:2]), 0
			var a []any = make([]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[3:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a[n] = o.Decode()
				}
			}
			return a
		case 0x9A:
			var m, n uint32 = endian.BigEndian.DecodeUint32(this[1:4]), 0
			var a []any = make([]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[5:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a[n] = o.Decode()
				}
			}
			return a
		case 0x9B:
			var m, n uint64 = endian.BigEndian.DecodeUint64(this[1:8]), 0
			var a []any = make([]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[9:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a[n] = o.Decode()
				}
			}
			return a
		case 0x9F:
			var a []any = make([]any,0)
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for true {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					a = append(a, o.Decode())
				}
			}
			return a
		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			var m, n int = int(tag-0xA0), 0
			var o map[string]any = make(map[string]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xB8:
			var m, n uint8 = uint8(this[1]), 0
			var o map[string]any = make(map[string]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[2:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xB9:
			var m, n uint16 = endian.BigEndian.DecodeUint16(this[1:2]), 0
			var o map[string]any = make(map[string]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[3:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xBA:
			var m, n uint32 = endian.BigEndian.DecodeUint32(this[1:4]), 0
			var o map[string]any = make(map[string]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[5:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xBB:
			var m, n uint64 = endian.BigEndian.DecodeUint64(this[1:8]), 0
			var o map[string]any = make(map[string]any,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[9:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xBF:
			var o map[string]any = make(map[string]any,1)
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error = nil
			for nil == e {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						a = ko.Decode()
						if nil != a {
							var k string = a.(string)
							o[k] = vo.Decode()
						}
					}
				}
			}
			return o
		case 0xC0, 0xC1:
			var a Object = Object{}
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			a, e = a.Read(b)
			if nil == e {
				return a.Decode()
			} 
		case 0xC2, 0xC3:
			var a big.Int
			a.SetBytes(this[1:])
			return a
		case 0xC4:
			// [TODO] rational
		case 0xC5:
			// [TODO] bigfloat
		case 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4:
			// [TODO] tag (content hints)
		case 0xD5, 0xD6, 0xD7:
			// [TODO] expected conversion (encoding/base)
		case 0xD8, 0xD9, 0xDA, 0xDB:
			// [TODO] tagged data
		case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3:
			// [TODO] simple value
		case 0xF4:
			return false
		case 0xF5:
			return true
		case 0xF6, 0xF7:
			return nil   // "null" and "undefined"
		case 0xF8:
			var a uint8 = this[1]
			return a
		case 0xF9:
			// [TODO] float16
		case 0xFA:
			var text []byte = this[1:4]
			var bits uint32 = endian.BigEndian.DecodeUint32(text)
			return math.Float32frombits(bits)

		case 0xFB:
			var text []byte = this[1:8]
			var bits uint64 = endian.BigEndian.DecodeUint64(text)
			return math.Float64frombits(bits)

		case 0xFF:
			return Break
		}
	}
	return nil
}
/*
 * Represent object structure.
 */
func (this Object) Describe() (string) {
	if this.HasTag() {
		var tag Tag = this.Tag()
		var desc string = fmt.Sprintf("<tag:%s>",this.MajorString())
		switch tag { 
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
		case 0x18:
			var cnt uint8 = this[1]
			desc = fmt.Sprintf("%s<uint8>",desc)
			switch cnt {
			case 2:
				desc = fmt.Sprintf("%s<uint16>",desc)
			case 4:
				desc = fmt.Sprintf("%s<uint32>",desc)
			case 8:
				desc = fmt.Sprintf("%s<uint64>",desc)
			default:
				desc = fmt.Sprintf("%s<int[%d]>",desc,cnt)
			}
		case 0x19:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			desc = fmt.Sprintf("%s<uint16><int[%d]>",desc,cnt)
		case 0x1A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			desc = fmt.Sprintf("%s<uint32><byte[%d]>",desc,cnt)
		case 0x1B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			desc = fmt.Sprintf("%s<uint64><byte[%d]>",desc,cnt)
		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
		case 0x38:
			var cnt uint8 = this[1]
			switch cnt {
			case 2:
				desc = fmt.Sprintf("%s<uint8><uint16>",desc)
			case 4:
				desc = fmt.Sprintf("%s<uint8><uint32>",desc)
			case 8:
				desc = fmt.Sprintf("%s<uint8><uint64>",desc)
			default:
				desc = fmt.Sprintf("%s<uint8><byte[%d]>",desc,cnt)
			}
		case 0x39:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			desc = fmt.Sprintf("%s<uint16><byte[%d]>",desc,cnt)
		case 0x3A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			desc = fmt.Sprintf("%s<uint32><byte[%d]>",desc,cnt)
		case 0x3B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			desc = fmt.Sprintf("%s<uint64><byte[%d]>",desc,cnt)
		case 0x40, 0x41, 0x42, 0x43, 0x44, 0x45, 0x46, 0x47, 0x48, 0x49, 0x4A, 0x4B, 0x4C, 0x4D, 0x4E, 0x4F, 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
			var m int = int(tag-0x40)
			desc = fmt.Sprintf("%s<byte[%d]>",desc,m)
		case 0x58:
			var cnt uint8 = this[1]
			desc = fmt.Sprintf("%s<uint8><byte[%d]>",desc,cnt)
		case 0x59:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			desc = fmt.Sprintf("%s<uint16><byte[%d]>",desc,cnt)
		case 0x5A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			desc = fmt.Sprintf("%s<uint32><byte[%d]>",desc,cnt)
		case 0x5B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			desc = fmt.Sprintf("%s<uint64><byte[%d]>",desc,cnt)
		case 0x5F:
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			for true {
				var o Object = Object{}
				var e error
				o, e = o.Read(b)
				if nil != e {
					if Break == e {
						desc = fmt.Sprintf("%s<break>",desc)
					}
					return desc
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc

		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			var m int = int(tag-0x60)
			desc = fmt.Sprintf("%s<byte[%d]>",desc,m)
		case 0x78:
			var cnt uint8 = this[1]
			desc = fmt.Sprintf("%s<uint8><byte[%d]>",desc,cnt)
		case 0x79:
			var cnt_ary []byte = this[1:2]
			var cnt uint16 = endian.BigEndian.DecodeUint16(cnt_ary)
			desc = fmt.Sprintf("%s<uint16><byte[%d]>",desc,cnt)
		case 0x7A:
			var cnt_ary []byte = this[1:4]
			var cnt uint32 = endian.BigEndian.DecodeUint32(cnt_ary)
			desc = fmt.Sprintf("%s<uint32><byte[%d]>",desc,cnt)
		case 0x7B:
			var cnt_ary []byte = this[1:8]
			var cnt uint64 = endian.BigEndian.DecodeUint64(cnt_ary)
			desc = fmt.Sprintf("%s<uint64><byte[%d]>",desc,cnt)
		case 0x7F:
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			for true {
				var o Object = Object{}
				var e error
				o, e = o.Read(b)
				if nil != e {
					desc = fmt.Sprintf(desc,"<break>")
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			var m, n int = int(tag-0x80), 0
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x98:
			var m, n uint8 = uint8(this[1]), 0
			desc = fmt.Sprintf("%s<uint8[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[2:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x99:
			var m, n uint16 = endian.BigEndian.DecodeUint16(this[1:2]), 0
			desc = fmt.Sprintf("%s<uint16[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[3:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x9A:
			var m, n uint32 = endian.BigEndian.DecodeUint32(this[1:4]), 0
			desc = fmt.Sprintf("%s<uint32[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[5:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x9B:
			var m, n uint64 = endian.BigEndian.DecodeUint64(this[1:8]), 0
			desc = fmt.Sprintf("%s<uint64[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[9:])
			var e error
			for n = 0; n < m; n++ {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0x9F:
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for true {
				var o Object = Object{}
				o, e = o.Read(b)
				if nil != e {
					desc = fmt.Sprintf("%s<break>",desc)
					break
				} else {
					desc = fmt.Sprintf("%s%s",desc,o.Describe())
				}
			}
			return desc
		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			var m, n int = int(tag-0xA0), 0
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xB8:
			var m, n uint8 = uint8(this[1]), 0
			desc = fmt.Sprintf("%s<uint8[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[2:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xB9:
			var m, n uint16 = endian.BigEndian.DecodeUint16(this[1:2]), 0
			desc = fmt.Sprintf("%s<uint16[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[3:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xBA:
			var m, n uint32 = endian.BigEndian.DecodeUint32(this[1:4]), 0
			desc = fmt.Sprintf("%s<uint32[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[5:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xBB:
			var m, n uint64 = endian.BigEndian.DecodeUint64(this[1:8]), 0
			desc = fmt.Sprintf("%s<uint64[%d]>",desc,m)
			var b *bytes.Buffer = bytes.NewBuffer(this[9:])
			var e error
			for n = 0; n < m; n++ {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xBF:
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error = nil
			for nil == e {
				var ko Object = Object{}
				ko, e = ko.Read(b)
				if nil != e {
					if Break == e {
						desc = fmt.Sprintf("%s<break>",desc)
					}
					break
				} else {
					var vo Object = Object{}
					vo, e = vo.Read(b)
					if nil != e {
						break
					} else {
						desc = fmt.Sprintf("%s%s",desc,ko.Describe())

						desc = fmt.Sprintf("%s%s",desc,vo.Describe())
					}
				}
			}
			return desc
		case 0xC0, 0xC1:
			var a Object = Object{}
			var b *bytes.Buffer = bytes.NewBuffer(this[1:])
			var e error
			a, e = a.Read(b)
			if nil == e {
				desc = fmt.Sprintf("%s%s",desc,a.Describe())
			}
			return desc
		case 0xC2, 0xC3:
		case 0xC4:
		case 0xC5:
		case 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4:
		case 0xD5, 0xD6, 0xD7:
		case 0xD8, 0xD9, 0xDA, 0xDB:
		case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3:
		case 0xF4:
		case 0xF5:
		case 0xF6, 0xF7:
		case 0xF8:
		case 0xF9:
		case 0xFA:
		case 0xFB:
		case 0xFF:
		}

		return desc
	} else {
		return ""
	}	
}
