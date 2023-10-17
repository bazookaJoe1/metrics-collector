package datacompressor

import (
	"compress/gzip"
	"compress/zlib"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

const EmptyString = ""

// ServerDecompressor is a middleware that defines decompressing method
// with request header "Content-Encoding" and wraps request body with appropriate
// io.ReadCloser. If request "Content-Encoding" is not supported
// it returns http.StatusBadRequest according to request "Content-Type"
func ServerDecompressor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		compressReader := decompressorSwitcher(c.Request())
		if compressReader == nil {
			return DefineResponseType(c)(c, http.StatusUnsupportedMediaType, EmptyString)
		}
		c.Request().Body = compressReader
		err := next(c)
		return err
	}
}

// decompressorSwitcher returns appropriate io.ReadCloser
// according to request "Content-Encoding". If "Content-Encoding"
// is not "gzip", "deflate" or "", it returns nil
func decompressorSwitcher(r *http.Request) io.ReadCloser {
	switch r.Header.Get(echo.HeaderContentEncoding) {
	case "gzip":
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil
		}
		return gz
	case "deflate":
		def, err := zlib.NewReader(r.Body)
		if err != nil {
			return nil
		}
		return def
	case EmptyString:
		return r.Body
	default:
		return nil
	}
}
