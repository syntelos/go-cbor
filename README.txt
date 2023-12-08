CBOR I/O for GOPL


  type CborObject interface {

	  Write(io.Writer) (error)

	  Read(io.Reader) (error)
  }


References

  [CBOR] https://tools.ietf.org/html/rfc8949
  [GOPL] https://go.dev/

