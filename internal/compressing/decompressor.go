package compressing

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
)

const gzipMethod = "gzip"
const deflateMethod = "deflate"

// DecompressorWrapper wraps io.ReadCloser with appropriate decompressor depending on input method.
func DecompressorWrapper(method string, reader io.ReadCloser) (io.ReadCloser, error) {
	switch method {
	case gzipMethod:
		gzReader, err := gzip.NewReader(reader)
		if err != nil {
			return nil, err
		}
		return gzReader, nil
	case deflateMethod:
		deflateReader, err := zlib.NewReader(reader)
		if err != nil {
			return nil, err
		}
		return deflateReader, nil
	case "":
		return reader, nil
	}
	return nil, fmt.Errorf("unsupported compression method: %s", method)
}
