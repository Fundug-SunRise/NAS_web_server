package main

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"nasweb/service"
	"net/http"
	"path"
)

var (
	logger   service.Logger
	sessions []string
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", redirectionHanderLogin)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/main", mainHandler)

	logger.PrintINFO("Сервис запущен на: localhost:8080")
	http.ListenAndServe(service.GetConfig().Port, nil)
}

func redirectionHanderLogin(w http.ResponseWriter, r *http.Request) {
	if !service.GetConfig().AuthEnabled || checkSessionCookie(r, "session") {

		http.Redirect(w, r, "/main", http.StatusFound)
		logger.PrintINFO("Переадресация на /main")
	} else {

		http.Redirect(w, r, "/login", http.StatusFound)
		logger.PrintINFO("Переадресация на /login")
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main handler"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "POST":

		AuthUser := checkUser(r.FormValue("username"), r.FormValue("password"))

		if AuthUser {

			if !checkSessionCookie(r, "session") {
				sessionsID, _ := generateSessionID()
				sessions = append(sessions, sessionsID)

				cookie := &http.Cookie{
					Name:     "session",
					Value:    sessionsID,
					Path:     "/login",
					HttpOnly: true,
				}

				http.SetCookie(w, cookie)
			}

			logger.PrintINFO("Переадресация на /main")
			http.Redirect(w, r, "/main", http.StatusFound)
		} else {
			logger.PrintINFO("Неудачная попытка входа")
			http.Redirect(w, r, "/login?error=auth_failed", http.StatusSeeOther)
		}

	case "GET":
		renderTemplate(w, getTemplate(w, "login"))
	default:
		logger.PrintFATAL("LoginHandler Метод не поддерживается")
	}
}

//Функции надо будет кудато вынести ибо они мне не встрались тут

func checkUser(username, password string) bool {
	return username == service.GetConfig().Login && password == service.GetConfig().Password
}

func renderTemplate(w http.ResponseWriter, tmpl *template.Template) {
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTemplate(w http.ResponseWriter, name string) *template.Template {

	templatePath := path.Join("static", name+".html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return tmpl
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func checkSessionCookie(r *http.Request, name string) bool {
	cookie, err := r.Cookie(name)
	if err != nil {
		return false
	}
	for _, sess := range sessions {
		if sess == cookie.Value {
			return true
		}
	}
	return false
}
