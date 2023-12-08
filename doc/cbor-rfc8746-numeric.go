/*
 * CBOR TAG RFC8746
 * Copyright 2023 John Douglas Pritchard, Syntelos
 */
package main

import (
	"fmt"
)

type CborTag = byte

const CborTagNumMajor byte = (0b010 << 5)

const CborTagNumLen8  byte = 0
const CborTagNumLen16  byte = 1
const CborTagNumLen32  byte = 2
const CborTagNumLen64 byte = 3

const CborTagNumFmtInt byte = 0
const CborTagNumFmtFlt byte = 1

const CborTagNumSigU byte = 0
const CborTagNumSigS byte = 1

const CborTagNumEndBig byte = 0
const CborTagNumEndLil byte = 1

type CborTagNum struct {
	fmt, sig, end, len CborTag
}

func (this CborTagNum) String() string {

	switch this.fmt {
	case CborTagNumFmtInt:
		switch this.sig {
		case CborTagNumSigU:
			switch this.end {
			case CborTagNumEndBig:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-int-u-big-8"
				case CborTagNumLen16:
					return "cbor-num-int-u-big-16"
				case CborTagNumLen32:
					return "cbor-num-int-u-big-32"
				case CborTagNumLen64:
					return "cbor-num-int-u-big-64"
				}
			case CborTagNumEndLil:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-int-u-lil-8"
				case CborTagNumLen16:
					return "cbor-num-int-u-lil-16"
				case CborTagNumLen32:
					return "cbor-num-int-u-lil-32"
				case CborTagNumLen64:
					return "cbor-num-int-u-lil-64"
				}
			}

		case CborTagNumSigS:
			switch this.end {
			case CborTagNumEndBig:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-int-s-big-8"
				case CborTagNumLen16:
					return "cbor-num-int-s-big-16"
				case CborTagNumLen32:
					return "cbor-num-int-s-big-32"
				case CborTagNumLen64:
					return "cbor-num-int-s-big-64"
				}
			case CborTagNumEndLil:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-int-s-lil-8"
				case CborTagNumLen16:
					return "cbor-num-int-s-lil-16"
				case CborTagNumLen32:
					return "cbor-num-int-s-lil-32"
				case CborTagNumLen64:
					return "cbor-num-int-s-lil-64"
				}
			}
		}

	case CborTagNumFmtFlt:
		switch this.sig {
		case CborTagNumSigU:
			switch this.end {
			case CborTagNumEndBig:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-flt-int-u-big-8"
				case CborTagNumLen16:
					return "cbor-flt-int-u-big-16"
				case CborTagNumLen32:
					return "cbor-flt-int-u-big-32"
				case CborTagNumLen64:
					return "cbor-flt-int-u-big-64"
				}
			case CborTagNumEndLil:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-flt-int-u-lil-8"
				case CborTagNumLen16:
					return "cbor-flt-int-u-lil-16"
				case CborTagNumLen32:
					return "cbor-flt-int-u-lil-32"
				case CborTagNumLen64:
					return "cbor-flt-int-u-lil-64"
				}
			}

		case CborTagNumSigS:
			switch this.end {
			case CborTagNumEndBig:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-flt-s-big-8"
				case CborTagNumLen16:
					return "cbor-num-flt-s-big-16"
				case CborTagNumLen32:
					return "cbor-num-flt-s-big-32"
				case CborTagNumLen64:
					return "cbor-num-flt-s-big-64"
				}
			case CborTagNumEndLil:
				switch this.len {
				case CborTagNumLen8:
					return "cbor-num-flt-s-lil-8"
				case CborTagNumLen16:
					return "cbor-num-flt-s-lil-16"
				case CborTagNumLen32:
					return "cbor-num-flt-s-lil-32"
				case CborTagNumLen64:
					return "cbor-num-flt-s-lil-64"
				}
			}
		}
	}

	return "cbor-num"
}

func (this CborTagNum) Binary() (tag CborTag) {
	tag = CborTagNumMajor
	if 0 != this.fmt {
		tag |= (this.fmt << 4)
	}
	if 0 != this.sig {
		tag |= (this.sig << 3)
	}
	if 0 != this.end {
		tag |= (this.end << 2)
	}
	if 0 != this.len {
		tag |= this.len
	}
	return tag
}

func Enumerate() (a []CborTagNum){
	var list_f = []byte {CborTagNumFmtInt, CborTagNumFmtFlt}
	var list_s = []byte {CborTagNumSigU, CborTagNumSigS}
 	var list_e = []byte {CborTagNumEndBig, CborTagNumEndLil}
	var list_l = []byte {CborTagNumLen8, CborTagNumLen16, CborTagNumLen32, CborTagNumLen64}
	for _, f := range list_f {
		for _, s := range list_s {
			for _, e := range list_e {
				for _, l := range list_l {
					var o CborTagNum
					o.fmt = f
					o.sig = s
					o.end = e
					o.len = l
					a = append(a,o)
				}
			}
		}
	}
	return a
}

func main(){

	for _, o := range Enumerate() {
		var s string = o.String()
		var b CborTag = o.Binary()

		fmt.Printf("%s\t0x%X\t0b%b\n", s, b, b)
	}
}
