package controller

import (
	"net/http"
	"os"
	"session-6/model"
	"session-6/repository"
	"session-6/utils.go"

	"github.com/labstack/echo/v4"
)

func (s *Server) Login(c echo.Context) error {
	var data model.User

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	userRepo := repository.UserRepository{
		DB: s.DB,
	}

	user, err := userRepo.GetUserByUsername(data.Username)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	if !utils.CompareHashPassword(data.Password, user.Password) {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "invalid username/password",
			"error":   nil,
		})
	}

	dbSession, _ := s.DBSession.Get(c.Request(), os.Getenv("SESSION_ID"))
	dbSession.Values["user"] = user.Username
	err = dbSession.Save(c.Request(), c.Response())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed to login",
			"error":   err.Error(),
		})
	}

	return c.Redirect(http.StatusFound, "/homepage")
}

func (s *Server) Logout(c echo.Context) error {
	session, _ := s.DBSession.Get(c.Request(), os.Getenv("SESSION_ID"))
	session.Options.MaxAge = -1
	session.Save(c.Request(), c.Response())

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "logout succesfully",
		"error":   nil,
	})
}

func (s *Server) Register(c echo.Context) error {
	var data model.User

	err := c.Bind(&data)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	userRepo := repository.UserRepository{
		DB: s.DB,
	}

	hashedPassword, err := utils.HashPassword(data.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	data.Password = hashedPassword
	err = userRepo.RegisterUser(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "failed",
			"error":   err.Error(),
		})
	}

	return c.Redirect(http.StatusFound, "/login")
}
