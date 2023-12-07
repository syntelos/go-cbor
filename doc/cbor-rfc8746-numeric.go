/*
 * CBOR TAG RFC8746
 * Copyright 2023 John Douglas Pritchard, Syntelos
 */
package main

import (
	"fmt"
)

type CborTagNum = byte

const CborTagNumMask byte = 0b01000000
const CborTagNumMaskFmt byte = 0b01010000
const CborTagNumMaskSig byte = 0b01001000
const CborTagNumMaskEnd byte = 0b01000100
const CborTagNumMaskLen byte = 0b01000011

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

func CborTagNumString(fmt, sig, end, len byte) (s string){
	switch fmt {
	case CborTagNumFmtInt:
		switch sig {
		case CborTagNumSigU:
			switch end {
			case CborTagNumEndBig:
				switch len {
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
				switch len {
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
			switch end {
			case CborTagNumEndBig:
				switch len {
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
				switch len {
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
		switch sig {
		case CborTagNumSigU:
			switch end {
			case CborTagNumEndBig:
				switch len {
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
				switch len {
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
			switch end {
			case CborTagNumEndBig:
				switch len {
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
				switch len {
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

func main(){
	const head string = "0b010"
	var list_f = []byte {CborTagNumFmtInt, CborTagNumFmtFlt}
	var list_s = []byte {CborTagNumSigU, CborTagNumSigS}
 	var list_e = []byte {CborTagNumEndBig, CborTagNumEndLil}
	var list_l = []byte {CborTagNumLen8, CborTagNumLen16, CborTagNumLen32, CborTagNumLen64}
	for _, f := range list_f {
		for _, s := range list_s {
			for _, e := range list_e {
				for _, l := range list_l {
					fmt.Printf("%s\t%s%b%b%b%02b\n",CborTagNumString(f,s,e,l),head,f,s,e,l)
				}
			}
		}
	}
}
