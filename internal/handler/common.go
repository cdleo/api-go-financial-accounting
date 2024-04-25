package handler

import "github.com/gorilla/mux"

type Handler interface {
	RegisterEndpoints(router *mux.Router)
}
