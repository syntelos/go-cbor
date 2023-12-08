CBOR I/O for GOPL

  type Object []byte

  type IO interface {

	  Write(io.Writer) (error)

	  Read(io.Reader) (error)
  }


hello, world.


  package main

  import (
	  "fmt"
	  "github.com/syntelos/go-cbor"
  )

  func main(){
	  var object cbor.Object = cbor.Object{
	      	     0x6D,0x68,0x65,0x6C,0x6C,0x6F,0x2C,
		     0x20,0x77,0x6f,0x72,0x6C,0x64,0x2E
		}

	  fmt.Println(object.Text())
  }


References

  [CBOR] https://tools.ietf.org/html/rfc8949
  [GOPL] https://go.dev/

