/*
 * CBOR table interpreter in GO
 * Copyright 2023 John Douglas Pritchard.  All rights reserved.
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
)
/*
 */
func usage(){
	fmt.Fprintln(os.Stderr,`
Synopsis

    table [print|enumerate]

Description

    Read table, and print or enumerate.

`)
	os.Exit(1)
}
/*
 */
type Line struct {
	first, last uint8
	description string
}
func rewrite(src []byte) (string) {
	var ch byte
	var idx, cnt int = 0, len(src)

	var tgt []byte = make([]byte,cnt)

	for ; idx < cnt; idx++ {
		ch = src[idx]

		if '"' == ch {
			tgt[idx] = '\''

		} else {
			tgt[idx] = ch
		}
	}

	return string(tgt)
}
func (this Line) ctor(inl []byte) (Line) {
	var x, z int = 0, len(inl)
	var lhs []byte
	var f, l string
	var n, m uint64
	var e error

	for ; x < z; x++ {
		if ' ' == inl[x] {

			lhs = inl[0:x]

			this.description = rewrite(inl[x+1:])
			{
				switch len(lhs) {
				case 4:
					f = string(lhs[2:4])
					n, e = strconv.ParseUint(f,16,8)
					if nil != e {
						this.first = 0xFF
						this.last = 0xFF
					} else {
						this.first = uint8(n)
						this.last = uint8(n)
					}
					
				case 10:
					f = string(lhs[2:4])
					l = string(lhs[8:])
					n, e = strconv.ParseUint(f,16,8)
					if nil != e {
						this.first = 0xFF
					} else {
						this.first = uint8(n)

						m, e = strconv.ParseUint(l,16,8)
						if nil != e {
							this.last = 0xFF
						} else {
							this.last = uint8(m)
						}
					}
				}
			}
			return this
		}
	}
	return this
}
func (this Line) print(){
	if this.first == this.last {
		fmt.Printf("0x%02X\t%s\n",this.first,this.description)
	} else {
		fmt.Printf("0x%02X-0x%02X\t%s\n",this.first,this.last,this.description)
	}
}
func (this Line) enumerate(){
	if this.first == this.last {

		fmt.Printf("case 0x%02X:\n\treturn \"%s\"\n",this.first,this.description)

	} else {
		var x, y uint8 = this.first, this.last

		fmt.Printf("case ")

		for ; x <= y; x++ {

			if this.first == x {

				fmt.Printf("0x%02X",x)
			} else {
				fmt.Printf(", 0x%02X",x)
			}
		}
		fmt.Printf(":\n\treturn \"%s\"\n",this.description)
	}
}
/*
 */
type Table struct {
	filename string
	records []Line
}

var table *Table = new(Table)

func (this *Table) size() (z int){

	return len(this.records)
}
func (this *Table) read(filename string) (e error){
	this.filename = filename

	var file *os.File
	file, e = os.Open(filename)
	if nil != e {
		e = fmt.Errorf("Error opening '%s': %w",filename,e)
		return e
	} else {
		defer file.Close()

		var reader *bufio.Reader = bufio.NewReader(file)

		var inl []byte
		var isp bool
		var lin Line
		inl, isp, e = reader.ReadLine()

		for true {
			if nil != e {
				if io.EOF == e {

					return nil
				} else {
					return fmt.Errorf("Error reading '%s': %w",filename,e)
				}
			} else if isp {
				return fmt.Errorf("Error reading '%s'.",filename)
			} else {
				this.records = append(this.records,lin.ctor(inl))

				inl, isp, e = reader.ReadLine()
			}
		}
		return nil
	}
}
func (this *Table) print(){

	var count int = table.size()
	fmt.Printf("# %s %d\n",this.filename,count)

	var index int = 0
	for ; index < count; index++ {
		this.records[index].print()
	}
}
func (this *Table) enumerate(){

	var count int = table.size()
	var index int = 0
	for ; index < count; index++ {
		this.records[index].enumerate()
	}
}
/*
 */
const location_rel string = "cbor-rfc8949-table.txt"
const location_doc string = "doc/cbor-rfc8949-table.txt"

func (this *Table) location() (string, error) {
	_, er := os.Stat("doc")
	if nil == er {
		_, er := os.Stat(location_doc)
		if nil == er {
			return location_doc, nil
		} else {
			return "", er
		}
	} else {
		_, er := os.Stat(location_rel)
		if nil == er {
			return location_rel, nil
		} else {
			return "", er
		}
	}
}
/*
 */
func main(){
	var argc int = len(os.Args)
	var argx int = 1
	if argx < argc {
		var opr string = os.Args[argx]

		var fin string
		var err error
		fin, err = table.location()
		if nil != err {
			fmt.Fprintf(os.Stderr,"table: file not found '%s'.\n",location_doc);
		} else {
			switch opr {

			case "print":
				e := table.read(fin)
				if nil != e {
					fmt.Fprintf(os.Stderr,"table: %v\n",e);
					os.Exit(1)
				} else {
					table.print()

					os.Exit(0)
				}

			case "enumerate":
				e := table.read(fin)
				if nil != e {
					fmt.Fprintf(os.Stderr,"table: %v\n",e);
					os.Exit(1)
				} else {
					table.enumerate()

					os.Exit(0)
				}

			default:
				usage()
			}
		}
	} else {
		usage()
	}
}
