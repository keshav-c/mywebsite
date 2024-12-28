package main

import (
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var templates = readTemplates("./templates")

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
	if len(title) == 0 {
		fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Editing New Page")
		renderTemplate(w, "edit", &Page{})
	} else {
		fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Editing title: "+title)
		p, err := LoadPage(title)
		if err != nil {
			p = &Page{Title: title}
		}
		renderTemplate(w, "edit", p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/save/"):]
	if title == "" {
		title = r.FormValue("title")
	}
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Saving title: "+title)
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.Validate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to save page: %s\n", err.Error())
		renderError("bad-request", w, err)
		return
	}
	err = p.Save()
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] Failed to save page: %s\n", err.Error())
		renderError("internal-error", w, err)
		return
	}
	fmt.Fprintf(os.Stdout, "[INFO] %s\n", "Redirecting to view/"+title)
	http.Redirect(w, r, "/view/"+title, http.StatusPermanentRedirect)
}

func renderTemplate(w http.ResponseWriter, tName string, p *Page) {
	err := templates.ExecuteTemplate(w, tName+".html", p)
	if err != nil {
		renderError("internal-error", w, err)
	}
}

func renderError(name string, w http.ResponseWriter, e error) {
	fmt.Fprintf(os.Stderr, "[ERROR] %+v\n", e)
	var responseCodes map[string]int = map[string]int{
		"not-found":      http.StatusNotFound,
		"internal-error": http.StatusInternalServerError,
		"bad-request":    http.StatusBadRequest,
	}
	status, ok := responseCodes[name]
	if !ok {
		errMsg := fmt.Sprintf("[ERROR] %s does not have an appropriate status code", name)
		fmt.Fprint(os.Stderr, errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	err := templates.ExecuteTemplate(w, name+".html", BlogError{Message: e.Error()})
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] template execution error for %s : %s\n", name+".html", err.Error())
		fmt.Fprint(os.Stderr, errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
	}
}

func readTemplates(root string) *template.Template {
	templateNames := make([]string, 0, 10)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			templateNames = append(templateNames, path)
		}
		return err
	})
	if err != nil {
		log.Fatalf("[ERROR] Failed to read app templates: %s\n", err.Error())
	}
	return template.Must(template.ParseFiles(templateNames...))
}
