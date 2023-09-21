package privatebin

import (
	"bytes"
	"compress/flate"
	"io"
)

func zlibDecompress(content []byte) []byte {
	reader := flate.NewReader(bytes.NewBuffer(content))
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	io.Copy(buffer, reader)

	return buffer.Bytes()
}
