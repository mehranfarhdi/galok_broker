package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"   // sqlite database driver
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

//func (server *Server) Run(s string) {
//
//}

func (server *Server) Initialize(db *gorm.DB) {

	server.DB = db

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)
	fmt.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(addr, corsHandler(server.Router)))
}
