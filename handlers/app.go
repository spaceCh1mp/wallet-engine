package handlers

import (
	"log"
	"net/http"
	"wallet-engine/errors"
)

// AppHandler defines the custom http.Handler used in the app
type AppHandler func(http.ResponseWriter, *http.Request) error

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := fn(w, r)
	if e != nil {
		log.Println(e)
		err, ok := e.(*errors.AppError)
		if ok {
			http.Error(w, err.Message, err.Code)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
