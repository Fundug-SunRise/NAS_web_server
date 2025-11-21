package main

import (
	"html/template"
	"nasweb/service"
	"net/http"
	"path"
)

var (
	logger   service.Logger
	AuthUser bool = false
)

func main() {

	http.Handle("/template/",
		http.StripPrefix("/template/",
			http.FileServer(http.Dir("./template"))))

	http.HandleFunc("/", redirectionHanderLogin)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/main", mainHandler)

	logger.PrintINFO("Сервис запущен на: localhost:8080")
	http.ListenAndServe(service.GetConfig().Port, nil)
}

func redirectionHanderLogin(w http.ResponseWriter, r *http.Request) {
	if !service.GetConfig().AuthEnabled || AuthUser {

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

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Ошибка parsing формы", http.StatusBadRequest)
			logger.PrintFATAL("Ошибка parsing формы: 400")
			return
		}

		AuthUser = checkUser(r.FormValue("username"), r.FormValue("password"))

		if AuthUser {
			logger.PrintINFO("Переадресация на /main")
			http.Redirect(w, r, "/main", http.StatusFound)
			return
		} else {
			logger.PrintINFO("Неудачная попытка входа")
			renderTemplate(w, getTemplate(w, "login"))
		}

	case "GET":
		renderTemplate(w, getTemplate(w, "login"))
	default:
		logger.PrintFATAL("LoginHandler Метод не поддерживается")
	}
}

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

	templatePath := path.Join("template", name, "index.html")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	return tmpl
}
