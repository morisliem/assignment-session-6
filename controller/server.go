package controller

import (
	"database/sql"

	"github.com/gorilla/sessions"
	"github.com/srinathgs/mysqlstore"
)

type Server struct {
	DB          *sql.DB
	CookieStore *sessions.CookieStore
	DBSession   *mysqlstore.MySQLStore
}
