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

    table [print|enumerate|list]

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
func (this Line) read(inl []byte) (Line) {

	var lhs, rhs []byte

	var n uint64
	var e error

	if '\t' == inl[4] {
		lhs = inl[0:4]
		rhs = inl[5:]

		this.description = rewrite(rhs)

		var f string = string(lhs[2:4])

		n, e = strconv.ParseUint(f,16,8)
		if nil == e {

			this.first = uint8(n)
			this.last = this.first
		}

	} else if '\t' == inl[9] {
		lhs = inl[0:9]
		rhs = inl[10:]

		this.description = rewrite(rhs)

		var f string = string(lhs[2:4])
		var l string = string(lhs[7:])

		n, e = strconv.ParseUint(f,16,8)
		if nil == e {

			this.first = uint8(n)

			n, e = strconv.ParseUint(l,16,8)
			if nil == e {

				this.last = uint8(n)
			}
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
func (this Line) list(){
	if this.first == this.last {

		fmt.Printf("0x%02X\n",this.first)

	} else {
		var x, y uint8 = this.first, this.last

		for ; x <= y; x++ {

			fmt.Printf("0x%02X\n",x)
		}
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
				this.records = append(this.records,lin.read(inl))

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
func (this *Table) list(){

	var count int = table.size()
	var index int = 0
	for ; index < count; index++ {
		this.records[index].list()
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

			case "list":
				e := table.read(fin)
				if nil != e {
					fmt.Fprintf(os.Stderr,"table: %v\n",e);
					os.Exit(1)
				} else {
					table.list()

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
