package main

import (
	"html/template"
	"nasweb/service"
	"net/http"
	"path"
)

var (
	logger service.Logger
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
	if service.GetConfig().AuthEnabled {
		http.Redirect(w, r, "login", http.StatusFound)
	} else {
		http.Redirect(w, r, "main", http.StatusFound)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Main handler"))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, getTemplate(w, "login"))
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
