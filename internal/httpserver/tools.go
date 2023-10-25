package httpserver

import (
	"github.com/bazookajoe1/metrics-collector/internal/compressing"
	"github.com/labstack/echo/v4"
	"strings"
)

// Compressor compresses data according to "Accept-Encoding" field in request header and sets response header "Content-Encoding". It chooses the first
// appropriate method. If there are no appropriate methods or error occurs in compressing, input data and response header will not be changed.
func Compressor(c echo.Context, data []byte) []byte {
	requestAllowEncodings := c.Request().Header.Get(echo.HeaderAcceptEncoding)

	if requestAllowEncodings == emptyString { // request host doesn't support encodings
		return data
	}

	acceptEncodings := strings.Split(requestAllowEncodings, ",")

	for _, encoding := range acceptEncodings {
		switch encoding {
		case "gzip":
			gzData, err := compressing.GZIPCompress(data)
			if err != nil {
				return data
			}

			c.Response().Header().Set(echo.HeaderContentEncoding, "gzip")
			return gzData
		case "deflate":
			deflateData, err := compressing.DeflateCompress(data)
			if err != nil {
				return data
			}

			c.Response().Header().Set(echo.HeaderContentEncoding, "deflate")
			return deflateData
		default:
			return data
		}
	}
	return data
}
