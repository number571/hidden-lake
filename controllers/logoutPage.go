package controllers

import (
    "net/http"
    "../connect"
    "../settings"
)

func logoutPage(w http.ResponseWriter, r *http.Request) {
	if settings.User.Auth {
		connect.Logout()
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
