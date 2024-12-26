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
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Viewing title: "+title)
	p, err := LoadPage(title)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to load page: %s\n", err.Error())
		renderError("not-found", w, err)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Editing title: "+title)
	p, err := LoadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Saving title: "+title)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to save page: %s\n", err.Error())
		renderError("internal-error", w, err)
		return
	}
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Redirecting to view/"+title)
	http.Redirect(w, r, "/view/"+title, http.StatusPermanentRedirect)
	// renderTemplate(w, "view", p)
}

func renderTemplate(w http.ResponseWriter, tName string, p *Page) {
	t, err := template.ParseFiles("templates/" + tName + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderError(name string, w http.ResponseWriter, e error) {
	var responseCodes map[string]int = map[string]int{
		"not-found":      http.StatusNotFound,
		"internal-error": http.StatusInternalServerError,
	}
	status, ok := responseCodes[name]
	if !ok {
		errMsg := "error " + name + " does not have an appropriate status code"
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles("templates/error/" + name + ".html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	t.Execute(w, struct{ Message string }{Message: e.Error()})
}
