package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	http.HandleFunc("/error/", errorHandler)
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := LoadPage(title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load page: %s\n", err.Error())
		http.Redirect(w, r, "/error/not-found", http.StatusTemporaryRedirect)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	p, err := LoadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Redirect(w, r, "/error/internal-error", http.StatusTemporaryRedirect)
		return
	}
	renderTemplate(w, "view", p)
}

func renderTemplate(w http.ResponseWriter, tName string, p *Page) {
	t, err := template.ParseFiles("templates/" + tName + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, p)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	var responseCodes map[string]int = map[string]int{
		"not-found":      http.StatusNotFound,
		"internal-error": http.StatusInternalServerError,
	}
	eName := r.URL.Path[len("/error/"):]
	t, _ := template.ParseFiles("templates/error/" + eName + ".html")
	status := responseCodes[eName]
	w.WriteHeader(status)
	t.Execute(w, Page{})
}
