package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"session-6/controller"

	"github.com/srinathgs/mysqlstore"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/go-sql-driver/mysql"

	_ "session-6/docs"

	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Assignment 2
// @version 1.0
// @description Simple Login & Register service
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4444
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {
	e := echo.New()

	key := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(key)

	connString := fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?parseTime=true",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)

	db, err := sql.Open("mysql", connString)
	if err != nil {
		panic(fmt.Errorf("error when open sql connection. error: %w", err))
	}

	err = db.Ping()
	if err != nil {
		panic(fmt.Errorf("error when open ping the database. error: %w", err))
	}

	dbSession, err := mysqlstore.NewMySQLStoreFromConnection(db, os.Getenv("SESSION_DB"), "/", 3600, []byte(os.Getenv("SESSION_KEY")))
	if err != nil {
		panic(fmt.Errorf("error when NewMySQLStoreFromConnection. error: %w", err))
	}

	server := controller.Server{
		DB:          db,
		CookieStore: store,
		DBSession:   dbSession,
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "server is alive ...")
	})

	v1 := e.Group("/api/v1")

	v1.POST("/register", server.Register)
	v1.POST("/login", server.Login)
	v1.POST("/logout", server.Logout)

	e.File("/register", "statics/register.html")
	e.File("/login", "statics/login.html")
	e.File("/homepage", "statics/homepage.html", server.CheckSession)
	// e.File("/homepage", "statics/homepage.html")
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Start(fmt.Sprintf(":%v", os.Getenv("SERVER_PORT")))

}
