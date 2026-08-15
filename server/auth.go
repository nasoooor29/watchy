package server

import (
	"backend/db"
	"backend/utils"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (s *Server) loginPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "login.html", pageData{}, http.StatusOK)
}
func (s *Server) registerPage(w http.ResponseWriter, r *http.Request) {
	s.render(w, r, "register.html", pageData{}, http.StatusOK)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	email, password := strings.ToLower(strings.TrimSpace(r.FormValue("email"))), r.FormValue("password")
	user, err := s.DB.GetUserByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		s.render(w, r, "login.html", pageData{Email: email, Error: "Email or password is incorrect."}, http.StatusUnauthorized)
		return
	}
	remmberMe := r.FormValue("remember") != ""
	err = utils.GenerateCookie(w, s.DB, user.Id, remmberMe)
	if err != nil {
		slog.Error("generating cookie failed", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	s.redirect(w, r, "/me")
}

type registerInput struct {
	Name     string `validate:"required,min=2,max=100"`
	Email    string `validate:"required,email,max=254"`
	Password string `validate:"required,min=8,max=72"`
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	input := registerInput{
		Name:     strings.TrimSpace(r.FormValue("name")),
		Email:    strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		Password: r.FormValue("password"),
	}
	data := pageData{Name: input.Name, Email: input.Email}
	if err := validate.Struct(input); err != nil {
		fmt.Printf("err: %v\n", err)
		data.Error = "Enter a name, a valid email address, and a password of at least 8 characters."
		s.render(w, r, "register.html", data, http.StatusBadRequest)
		return
	}
	_, err := s.DB.GetUserByEmail(input.Email)
	if err == nil {
		data.Error = "An account with that email already exists."
		s.render(w, r, "register.html", data, http.StatusConflict)
		return
	}
	if !db.IsNotFound(err) {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	user, err := s.DB.CreateUser(input.Name, input.Email, string(hash))
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	err = utils.GenerateCookie(w, s.DB, user.Id, false)
	if err != nil {
		slog.Error("generating cookie failed", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	s.redirect(w, r, "/me")
}

func (s *Server) mePage(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserByCookie(r, s.DB)
	s.render(w, r, "me.html", pageData{Name: user.Name, Email: user.Email}, http.StatusOK)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	if cookie := utils.GetCookie(r, utils.AUTH_COOKE_NAME); cookie != nil {
		_ = s.DB.DeleteCookie(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{Name: utils.AUTH_COOKE_NAME, Value: "", Path: "/", MaxAge: -1, HttpOnly: true, SameSite: http.SameSiteLaxMode})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
