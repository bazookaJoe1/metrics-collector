package datacompressor

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"github.com/labstack/echo/v4"
	"strings"
)

const ContentEncodingGZIP = "gzip"
const ContentEncodingDeflate = "deflate"

var GZIPHeaderParams = map[string]string{
	echo.HeaderContentType:     echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding:  ContentEncodingGZIP + ", " + ContentEncodingDeflate,
	echo.HeaderContentEncoding: ContentEncodingGZIP,
}

var DeflateHeaderParams = map[string]string{
	echo.HeaderContentType:     echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding:  ContentEncodingGZIP + ", " + ContentEncodingDeflate,
	echo.HeaderContentEncoding: ContentEncodingDeflate,
}

var UncompressedHeaderParams = map[string]string{
	echo.HeaderContentType:    echo.MIMEApplicationJSONCharsetUTF8,
	echo.HeaderAcceptEncoding: ContentEncodingGZIP + ", " + ContentEncodingDeflate,
}

// GZIPCompress compresses data in d with gzip algo.
// Return value is compressed data and error if occurs
func GZIPCompress(d []byte) ([]byte, error) {
	var gzData bytes.Buffer
	gz, err := gzip.NewWriterLevel(&gzData, gzip.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = gz.Write(d)
	if err != nil {
		return nil, err
	}
	err = gz.Close()
	if err != nil {
		return nil, err
	}

	return gzData.Bytes(), nil
}

// DeflateCompress compresses data in d with zlib algo.
// Return value is compressed data and error if occurs
func DeflateCompress(d []byte) ([]byte, error) {
	var zlibData bytes.Buffer
	zl, err := zlib.NewWriterLevel(&zlibData, zlib.BestCompression)
	if err != nil {
		return nil, err
	}
	_, err = zl.Write(d)
	if err != nil {
		return nil, err
	}
	err = zl.Close()
	if err != nil {
		return nil, err
	}

	return zlibData.Bytes(), nil
}

// ServerCompress performs compressing of data according to "Accept-Encoding" value
// of request. It looks for the first occasion of "Accept-Encoding" param and compresses input data with that algorithm.
// Allowed encodings are [gzip, deflate]. If compressing is not supported, func returns input data and error
func ServerCompress(c echo.Context, d []byte) ([]byte, error) {
	acptEnc := c.Request().Header.Get(echo.HeaderAcceptEncoding)
	if acptEnc == EmptyString {
		for key, value := range UncompressedHeaderParams {
			c.Response().Header().Set(key, value)
		}
		return d, nil
	}
	encs := strings.Split(acptEnc, ",")
	for _, enc := range encs {
		switch enc {
		case "gzip":
			gzData, err := GZIPCompress(d)
			if err != nil {
				for key, value := range UncompressedHeaderParams {
					c.Response().Header().Set(key, value)
				}
				return d, err
			}

			for key, value := range GZIPHeaderParams {
				c.Response().Header().Set(key, value)
			}
			return gzData, nil

		case "deflate":
			zlibData, err := DeflateCompress(d)
			if err != nil {
				for key, value := range UncompressedHeaderParams {
					c.Response().Header().Set(key, value)
				}
				return d, err
			}

			for key, value := range DeflateHeaderParams {
				c.Response().Header().Set(key, value)
			}
			return zlibData, nil
		}
	}

	for key, value := range UncompressedHeaderParams {
		c.Response().Header().Set(key, value)
	}
	return d, fmt.Errorf("compressing %v is not supported", encs)
}
