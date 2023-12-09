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
	"errors"
	"fmt"
	"io"
	"github.com/syntelos/go-endian"
)
/*
 */
var Break error = errors.New("CBOR Break")

var ErrorUnrecognizedTag error = errors.New("Unrecognized CBOR Tag")
var ErrorMissingData error = errors.New("Missing CBOR Data")

/*
 * Principal user interface.
 */
type IO interface {

	Write(io.Writer) (error)

	Read(io.Reader) (error)
}
/*
 * Encoded data set.
 */
type Object []byte
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
 */
func (this Object) Write(w io.Writer) (e error){
	_, e = w.Write(this)
	return e
}
/*
 */
func (this Object) Read(r io.Reader) (e error){
	var tag []byte = make([]byte,1)
	var m, n int

	n, e = r.Read(tag)
	if nil != e {
		return e
	} else if 1 != n {
		return fmt.Errorf("Read (%d) expected (1).",n)
	} else {
		var d []byte
		var t byte = tag[0]
		var a, b Object

		switch t {
		case 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17:
			/* unsigned integer 0x00..0x17 (0..23)
			 */
			this = tag
			return nil

		case 0x18:
			/* unsigned integer (one-byte uint8_t follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x19:
			/* unsigned integer (two-byte uint16_t follows)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x1A:
			/* unsigned integer (four-byte uint32_t follows)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x1B:
			/* unsigned integer (eight-byte uint64_t follows)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2A, 0x2B, 0x2C, 0x2D, 0x2E, 0x2F, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37:
			/* negative integer -1-0x00..-1-0x17 (-1..-24)
			 */
			this = tag
			return nil

		case 0x38:
			/* negative integer -1-n (one-byte uint8_t for n follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x39:
			/* negative integer -1-n (two-byte uint16_t for n follows)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x3A:
			/* negative integer -1-n (four-byte uint32_t for n follows)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint32(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x3B:
			/* negative integer -1-n (eight-byte uint64_t for n follows)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint64(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
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
				return fmt.Errorf("Data: %w",e)
			} else if m != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x58:
			/* byte string (one-byte uint8_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x59:
			/* byte string (two-byte uint16_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x5A:
			/* byte string (four-byte uint32_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint32(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x5B:
			/* byte string (eight-byte uint64_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint64(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x5F:
			/* byte string, byte strings follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = make([]byte,0)
				e = a.Read(r)
				if nil == e {

					this = concatenate(this,a)

				} else if Break == e {
					e = nil
					break
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}
			return nil

		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			/* UTF-8 string (0x00..0x17 bytes follow)
			 */
			this = tag
			m = int(t-0x60)
			d = make([]byte,m)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if m != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0x78:
			/* UTF-8 string (one-byte uint8_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var z int = int(d[0])
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x79:
			/* UTF-8 string (two-byte uint16_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var z int = int(endian.BigEndian.DecodeUint16(d))
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != n {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x7A:
			/* UTF-8 string (four-byte uint32_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var z uint32 = endian.BigEndian.DecodeUint32(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint32(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}
			}

		case 0x7B:
			/* UTF-8 string (eight-byte uint64_t for n, and then n bytes follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var z uint64 = endian.BigEndian.DecodeUint64(d)
				var p []byte = make([]byte,z)
				n, e = r.Read(p)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else if z != uint64(n) {
					return ErrorMissingData
				} else {
					d = concatenate(d,p)
					this = concatenate(this,d)
					return nil
				}	
			}

		case 0x7F:
			/* UTF-8 string, UTF-8 strings follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = make([]byte,0)
				e = a.Read(r)
				if nil == e {

					this = concatenate(this,a)

				} else if Break == e {
					e = nil
					break
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}
			return nil

		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89, 0x8A, 0x8B, 0x8C, 0x8D, 0x8E, 0x8F, 0x90, 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97:
			/* array (0x00..0x17 data items follow)
			 */
			this = tag
			m = int(t-0x80)
			for n = 0; n < m; n++ {
				a = make([]byte,0)
				e = a.Read(r)
				if nil == e {

					this = concatenate(this,a)

				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}
			return nil

		case 0x98:
			/* array (one-byte uint8_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var z int = int(d[0])
				for n = 0; n < z; n++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil == e {

						this = concatenate(this,a)

					} else {
						return fmt.Errorf("Data: %w",e)
					}
				}
				return nil
			}

		case 0x99:
			/* array (two-byte uint16_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint16 = 0, endian.BigEndian.DecodeUint16(d)
				for ; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil == e {

						this = concatenate(this,a)

					} else {
						return fmt.Errorf("Data: %w",e)
					}
				}
				return nil
			}

		case 0x9A:
			/* array (four-byte uint32_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint32 = 0, endian.BigEndian.DecodeUint32(d)
				for ; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil == e {

						this = concatenate(this,a)

					} else {
						return fmt.Errorf("Data: %w",e)
					}
				}
				return nil
			}

		case 0x9B:
			/* array (eight-byte uint64_t for n, and then n data items follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint64 = 0, endian.BigEndian.DecodeUint64(d)
				for ; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil == e {

						this = concatenate(this,a)

					} else {
						return fmt.Errorf("Data: %w",e)
					}
				}
				return nil
			}

		case 0x9F:
			/* array, data items follow, terminated by 'break'
			 */
			this = tag
			for nil == e {
				a = make([]byte,0)
				e = a.Read(r)
				if nil == e {

					this = concatenate(this,a)

				} else if Break == e {
					e = nil
					break
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}
			return nil

		case 0xA0, 0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8, 0xA9, 0xAA, 0xAB, 0xAC, 0xAD, 0xAE, 0xAF, 0xB0, 0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7:
			/* map (0x00..0x17 pairs of data items follow)
			 */
			this = tag
			m, n = 0, int(t-0xA0)
			for ; m < n; m++ {
				a = make([]byte,0)
				e = a.Read(r)
				if nil != e {
					return fmt.Errorf("Data: %w",e)
				} else {
					this = concatenate(this,a)

					b = make([]byte,0)
					e = b.Read(r)
					if nil != e {
						return fmt.Errorf("Data: %w",e)
					} else {
						this = concatenate(this,b)
					}	
				}
			}
			return nil

		case 0xB8:
			/* map (one-byte uint8_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint8 = 0, uint8(d[0])
				for x = 0; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil != e {
						return fmt.Errorf("Data: %w",e)
					} else {
						this = concatenate(this,a)
						b = make([]byte,0)
						e = b.Read(r)
						if nil != e {
							return fmt.Errorf("Data: %w",e)
						} else {
							this = concatenate(this,b)
						}	
					}
				}
				return nil
			}

		case 0xB9:
			/* map (two-byte uint16_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint16 = 0, endian.BigEndian.DecodeUint16(d)
				for x = 0; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil != e {
						return fmt.Errorf("Data: %w",e)
					} else {
						this = concatenate(this,a)
						b = make([]byte,0)
						e = b.Read(r)
						if nil != e {
							return fmt.Errorf("Data: %w",e)
						} else {
							this = concatenate(this,b)
						}	
					}
				}
				return nil
			}

		case 0xBA:
			/* map (four-byte uint32_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint32 = 0, endian.BigEndian.DecodeUint32(d)
				for x = 0; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil != e {
						return fmt.Errorf("Data: %w",e)
					} else {
						this = concatenate(this,a)
						b = make([]byte,0)
						e = b.Read(r)
						if nil != e {
							return fmt.Errorf("Data: %w",e)
						} else {
							this = concatenate(this,b)
						}	
					}
				}
				return nil
			}

		case 0xBB:
			/* map (eight-byte uint64_t for n, and then n pairs of data items follow)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				var x, z uint64 = 0, endian.BigEndian.DecodeUint64(d)
				for x = 0; x < z; x++ {
					a = make([]byte,0)
					e = a.Read(r)
					if nil != e {
						return fmt.Errorf("Data: %w",e)
					} else {
						this = concatenate(this,a)
						b = make([]byte,0)
						e = b.Read(r)
						if nil != e {
							return fmt.Errorf("Data: %w",e)
						} else {
							this = concatenate(this,b)
						}	
					}
				}
				return nil
			}

		case 0xBF:
			/* map, pairs of data items follow, terminated by 'break'
			 */
			this = tag

			for nil == e {
				a = make([]byte,0)
				e = a.Read(r)
				if nil == e {
					this = concatenate(this,a)

					b = make([]byte,0)
					e = b.Read(r)
					if nil == e {
						this = concatenate(this,b)

					} else {
						return fmt.Errorf("Data: %w",e)
					}
				} else if Break == e {
					e = nil
					break
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}
			return nil

		case 0xC0, 0xC1:
			/* date/time (data item follows; see Section 3.4.1 and 3.4.2)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xC2:
			/* unsigned bignum (data item 'byte string' follows)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xC3:
			/* negative bignum (data item 'byte string' follows)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xC4:
			/* decimal Fraction (data item 'array' follows; see Section 3.4.4)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xC5:
			/* bigfloat (data item 'array' follows; see Section 3.4.4)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xC6, 0xC7, 0xC8, 0xC9, 0xCA, 0xCB, 0xCC, 0xCD, 0xCE, 0xCF, 0xD0, 0xD1, 0xD2, 0xD3, 0xD4:
			/* (tag)
			 */
			this = tag
			return nil

		case 0xD5, 0xD6, 0xD7:
			/* expected conversion (data item follows; see Section 3.4.5.2)
			 */
			this = tag
			a = make([]byte,0)
			e = a.Read(r)
			if nil == e {
				this = concatenate(this,a)
				return nil
			} else {
				return fmt.Errorf("Data: %w",e)
			}

		case 0xD8:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,1)
			n, e = r.Read(a)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return fmt.Errorf("Data expected (1) found (%d).",n)
			} else {
				this = concatenate(this,a)

				b = make([]byte,0)
				e = b.Read(r)
				if nil == e {
					this = concatenate(this,b)

					return nil
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}

		case 0xD9:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,2)
			n, e = r.Read(a)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return fmt.Errorf("Data expected (2) found (%d).",n)
			} else {
				this = concatenate(this,a)

				b = make([]byte,0)
				e = b.Read(r)
				if nil == e {
					this = concatenate(this,b)

					return nil
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}

		case 0xDA:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,4)
			n, e = r.Read(a)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return fmt.Errorf("Data expected (4) found (%d).",n)
			} else {
				this = concatenate(this,a)

				b = make([]byte,0)
				e = b.Read(r)
				if nil == e {
					this = concatenate(this,b)

					return nil
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}

		case 0xDB:
			/* (more tags; 1/2/4/8 bytes of tag number and then a data item follow)
			 */
			this = tag
			a = make([]byte,8)
			n, e = r.Read(a)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return fmt.Errorf("Data expected (8) found (%d).",n)
			} else {
				this = concatenate(this,a)

				b = make([]byte,0)
				e = b.Read(r)
				if nil == e {
					this = concatenate(this,b)

					return nil
				} else {
					return fmt.Errorf("Data: %w",e)
				}
			}

		case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF, 0xF0, 0xF1, 0xF2, 0xF3:
			/* (simple value)
			 */
			this = tag
			return nil

		case 0xF4:
			/* "false"
			 */
			this = tag
			return nil

		case 0xF5:
			/* "true"
			 */
			this = tag
			return nil

		case 0xF6:
			/* "null"
			 */
			this = tag
			return nil

		case 0xF7:
			/* "undefined"
			 */
			this = tag
			return nil

		case 0xF8:
			/* (simple value, one byte follows)
			 */
			this = tag
			d = make([]byte,1)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 1 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0xF9:
			/* half-precision float (two-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,2)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 2 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0xFA:
			/* single-precision float (four-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,4)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 4 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0xFB:
			/* double-precision float (eight-byte IEEE 754)
			 */
			this = tag
			d = make([]byte,8)
			n, e = r.Read(d)
			if nil != e {
				return fmt.Errorf("Data: %w",e)
			} else if 8 != n {
				return ErrorMissingData
			} else {
				this = concatenate(this,d)
				return nil
			}

		case 0xFF:
			/* 'break' stop code"
			 */
			this = tag
			return Break

		default:
			return ErrorUnrecognizedTag
		}
	}
}
/*
 */
func (this Object) String() string {
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
 */
func (this Object) HasTag() bool {
	var z int = len(this)
	return (0 < z)
}
/*
 */
func (this Object) Tag() Tag {
	if this.HasTag() {
		return Tag(this[0])
	} else {
		return 0
	}
}
/*
 */
func (this Object) Major() Major {
	if this.HasTag() {
		var tag Tag = this.Tag()
		var major Major = Major((tag & 0xE0)>>5)
		return major
	} else {
		return Major(0)
	}
}
/*
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
		case MajorSimple:
			return "float, simple, break"
		default:
			return ""
		}
	} else {
		return ""
	}
}
/*
 */
func (this Object) HasText() bool {
	return (this.HasTag() && MajorText == this.Major())
}
/*
 */
func (this Object) Text() (s string) {
	if this.HasText() {
		var tag Tag = this.Tag()
		switch tag {
		case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77:
			var m int = int(tag-0x60)
			var text []byte = this[1:(m+1)]
			return string(text)

		case 0x78:
		case 0x79:
		case 0x7A:
		case 0x7B:

		}
	}
	return ""
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
func concatenate(a []byte, b []byte) (c []byte) {

	var a_len int = len(a)
	if 0 == a_len {
		return b
	} else {
		var b_len int = len(b)
		if 0 == b_len {
			return b
		} else {
			var c_len int = (a_len+b_len)

			c = make([]byte,c_len)

			copier(c,0,c_len,a,0,a_len)

			return copier(c,a_len,c_len,b,0,b_len)
		}
	}
}
