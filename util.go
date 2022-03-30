package MonkeSockets

import "github.com/labstack/echo/v4"

func ServeStaticClientCode(e *echo.Echo, url string) {
	e.File("url", "MonkeSocket.js")
}
