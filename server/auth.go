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

type authPageData struct{ Name, Email, Error string }

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	email, password := strings.ToLower(strings.TrimSpace(r.FormValue("email"))), r.FormValue("password")
	user, err := s.DB.GetUserByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		s.render(w, r, http.StatusUnauthorized, "login.html", authPageData{
			Email: email,
			Error: "Email or password is incorrect.",
		})
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
	data := authPageData{Name: input.Name, Email: input.Email}
	if err := validate.Struct(input); err != nil {
		fmt.Printf("err: %v\n", err)
		data.Error = "Enter a name, a valid email address, and a password of at least 8 characters."
		s.render(w, r, http.StatusBadRequest, "register.html", data)
		return
	}
	_, err := s.DB.GetUserByEmail(input.Email)
	if err == nil {
		data.Error = "An account with that email already exists."
		s.render(w, r, http.StatusConflict, "register.html", data)
		return
	}
	if !db.IsNotFound(err) {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	user, err := s.DB.UpsertUser("", input.Name, input.Email, input.Password)
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

type updateUserInput struct {
	Name     string `validate:"required,min=2,max=100"`
	Email    string `validate:"required,email,max=254"`
	Password string `validate:"omitempty,min=8,max=72"`
}

func (s *Server) UpdateUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUserByCookie(r, s.DB)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	uinp := updateUserInput{
		Name:     strings.TrimSpace(r.FormValue("name")),
		Email:    strings.ToLower(strings.TrimSpace(r.FormValue("email"))),
		Password: r.FormValue("password"),
	}

	if err := validate.Struct(uinp); err != nil {
		fmt.Printf("err: %v\n", err)
		s.render(w, r, http.StatusBadRequest, "me.html", authPageData{
			Name:  uinp.Name,
			Email: uinp.Email,
			Error: "Enter a name, a valid email address, and a password of at least 8 characters.",
		})
		return
	}

	user, err = s.DB.UpsertUser(user.Id, uinp.Name, uinp.Email, uinp.Password)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	s.render(w, r, http.StatusOK, "me.html", authPageData{Name: user.Name, Email: user.Email})
}

func (s *Server) mePage(w http.ResponseWriter, r *http.Request) {
	user, _ := utils.GetUserByCookie(r, s.DB)
	s.render(w, r, http.StatusOK, "me.html", authPageData{Name: user.Name, Email: user.Email})
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	cookie := utils.GetCookie(r, utils.AUTH_COOKE_NAME)
	err := s.DB.DeleteCookie(cookie.Value)
	if err != nil {
		slog.Error("", "err", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     utils.AUTH_COOKE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
