package compressing

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
)

// GZIPCompress compresses data with gzip algorithm.
// Return value is compressed data and error if occurs.
func GZIPCompress(data []byte) ([]byte, error) {
	var gzData bytes.Buffer
	gz, err := gzip.NewWriterLevel(&gzData, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(data)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return gzData.Bytes(), nil
}

// DeflateCompress compresses data with zlib algo.
// Return value is compressed data and error if occurs.
func DeflateCompress(data []byte) ([]byte, error) {
	var zlibData bytes.Buffer
	zl, err := zlib.NewWriterLevel(&zlibData, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = zl.Write(data)
	if err != nil {
		return nil, err
	}
	err = zl.Close()
	if err != nil {
		return nil, err
	}

	return zlibData.Bytes(), nil
}
