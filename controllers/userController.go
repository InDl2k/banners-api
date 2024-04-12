package controllers

import (
	u "banners/internal/utils"
	"banners/internal/utils/jwt"
	"github.com/gorilla/mux"
	"net/http"
)

var GetToken = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	token, err := jwt.CreateToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		u.Respond(w, u.MessageError("Ошибка при создании токена"))
		return
	}

	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{"token": token}
	u.Respond(w, &resp)
}
