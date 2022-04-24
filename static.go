package MonkeSockets

import (
	"github.com/labstack/echo/v4"
)

func Static(e *echo.Echo, path string) {
	e.File(path, "MonkeSockets/MonkeSocket.js")
}
