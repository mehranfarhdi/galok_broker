package controllers

import (
	"net/http"

	"github.com/mehranfarhdi/galok_broker/src/api/response"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
