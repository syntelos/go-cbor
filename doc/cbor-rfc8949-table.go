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

    table [print|enumerate] <file>

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
func (this Line) ctor(inl []byte) (Line) {
	var x, z int = 0, len(inl)
	var lhs []byte
	var f, l string
	var n, m uint64
	var e error

	for ; x < z; x++ {
		if ' ' == inl[x] {

			lhs = inl[0:x]

			this.description = string(inl[x+1:])
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
		fmt.Printf("0x%02X\t%s\n",this.first,this.description)
	} else {
		var x, y uint8 = this.first, this.last
		
		for ; x <= y; x++ {

			fmt.Printf("0x%02X\t%s\n",x,this.description)
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
func main(){
	var argc int = len(os.Args)
	var argx int = 1
	if argx < argc {
		var opr string = os.Args[argx]

		argx += 1
		if argx < argc {
			var filename string = os.Args[argx]

			switch opr {

			case "print":
				e := table.read(filename)
				if nil != e {
					fmt.Fprintf(os.Stderr,"table: %v\n",e);
					os.Exit(1)
				} else {
					table.print()

					os.Exit(0)
				}

			case "enumerate":
				e := table.read(filename)
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
		} else {
			usage()
		}
	} else {
		usage()
	}
}
