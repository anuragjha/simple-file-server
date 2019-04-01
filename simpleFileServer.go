package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func load(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func renderHtmlTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	fmt.Println("in render Template")
	cwd, _ := os.Getwd()
	editFileAddress := filepath.Join(cwd, "/roughWork/edit.html")
	viewFileAddress := filepath.Join(cwd, "/roughWork/view.html")
	var templates = template.Must(template.ParseFiles(editFileAddress, viewFileAddress))

	//fileaddress := filepath.Join(cwd, "/"+tmpl+".html")
	//fileaddress := filepath.Join(tmpl + ".html")
	fmt.Println("File address  : ", tmpl+".html")
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	ls := strings.Split(r.URL.Path, "/")
	w.Write([]byte("I love this - " + string(ls[1])))
}

//func viewHandler(w http.ResponseWriter, r *http.Request) {
func viewHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := load(title)

	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	} else {
		renderHtmlTemplate(w, "view", p)
	}

}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {

	p, err := load(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderHtmlTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {

	body := r.FormValue("body") // getting form value

	p := &Page{ // creating struct Page
		Title: title,
		Body:  []byte(body),
	}

	err := p.save() // saving Page struct to file

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc { //read function literals
	return func(w http.ResponseWriter, r *http.Request) {
		validPath := regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$") // functionality of getTitle
		valid := validPath.FindStringSubmatch(r.URL.Path)
		if valid == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, valid[2])
	}
}

func main() {

	//http.HandleFunc("/", handler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":3000", nil)) //running server

}
