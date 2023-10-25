package httpserver

import (
	"github.com/bazookajoe1/metrics-collector/internal/compressing"
	"github.com/labstack/echo/v4"
	"net/http"
)

// Decompressor is the middleware that defines request "Content-Encoding" and replaces original
// body io.Reader with new.
func Decompressor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		method := c.Request().Header.Get(echo.HeaderContentEncoding)

		newBody, err := compressing.DecompressorWrapper(method, c.Request().Body) // wrapping decompressor

		if err != nil { // if not supported return error

			return c.NoContent(http.StatusUnsupportedMediaType) // DecompressorWrapper can return not only unsupported method,
			// error can occur in creating io.Reader, but this should not happen.
			// So, I don't handle this case.
		}

		c.Request().Body = newBody
		err = next(c)
		return err
	}
}
