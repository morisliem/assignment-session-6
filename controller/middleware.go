package controller

import (
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
)

func (s *Server) CheckSession(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		session, _ := s.DBSession.Get(c.Request(), os.Getenv("SESSION_ID"))
		if len(session.Values) == 0 {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"status": "failed",
				"error":  "no session",
			})
		}

		return c.HTML(http.StatusOK, `
        <!doctype html>
        <html>

        <head>
            <meta charset="utf-8">
            <title>HomePage</title>
            <h1>Homepage</h1>
        </head>

        <h2 class="user-name"></h2>

        <form action="/api/v1/logout" method="POST">
            <button>Logout</button>
        </form>

        <body>
        </body>

        </html>
        `)
	}
}
