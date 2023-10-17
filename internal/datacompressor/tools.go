package datacompressor

import (
	"github.com/labstack/echo/v4"
)

type ResponseFunc func(echo.Context, int, interface{}) error

func DefineResponseType(c echo.Context) ResponseFunc {
	switch c.Request().Header.Get(echo.HeaderContentType) {
	case echo.MIMEApplicationJSONCharsetUTF8:
		return func(c echo.Context, statusCode int, data interface{}) error {
			return c.JSON(statusCode, data)
		}
	case echo.MIMEApplicationJSON:
		return func(c echo.Context, statusCode int, data interface{}) error {
			return c.JSON(statusCode, data)
		}
	default:
		return func(c echo.Context, statusCode int, data interface{}) error {
			return c.String(statusCode, data.(string))
		}
	}
}
