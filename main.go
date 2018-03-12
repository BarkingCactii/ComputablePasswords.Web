package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"./viewmodel"
)

func PopulateTemplates() map[string]*template.Template {
	result := make(map[string]*template.Template)
	const basePath = "templates"
	layout := template.Must(template.ParseFiles(basePath + "/_layout.html"))
	template.Must(layout.ParseFiles(basePath+"/_header.html", basePath+"/_footer.html"))
	dir, err := os.Open(basePath + "/content")

	if err != nil {
		panic("Failed to open template blocks directory: " + err.Error())
	}
	fis, err := dir.Readdir(-1)
	if err != nil {
		panic("Failed to read contents of content directory: " + err.Error())
	}
	for _, fi := range fis {
		f, err := os.Open(basePath + "/content/" + fi.Name())
		if err != nil {
			panic("Failed to open template '" + fi.Name() + "'")
		}
		content, err := ioutil.ReadAll(f)
		if err != nil {
			panic("Failed to read content from file '" + fi.Name() + "'")
		}
		f.Close()
		tmpl := template.Must(layout.Clone())
		_, err = tmpl.Parse(string(content))
		if err != nil {
			panic("Failed to parse contents of '" + fi.Name() + "' as template")
		}
		result[fi.Name()] = tmpl
	}
	return result
}

func saveHandler(w http.ResponseWriter, r *http.Request) { //, title string) {

	context.WebURL = r.FormValue("url")
	context.PublicKey = r.FormValue("pepper")
	fmt.Printf("InputM = [%s], public key = [%s]\n", context.WebURL, context.PublicKey)
	context.Result = "Password: " + calc(context.WebURL, context.PublicKey)

	http.Redirect(w, r, "/main", http.StatusFound)
}

var context = viewmodel.NewBase()

func main() {
	templates := PopulateTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedFile := r.URL.Path[1:]
		template := templates[requestedFile+".html"]

		if template != nil {
			err := template.Execute(w, context) //context)
			if err != nil {
				log.Println(err)
			}
		} else {
			w.WriteHeader(404)
		}
	})
	http.Handle("/img/", http.FileServer(http.Dir("public")))
	http.Handle("/css/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/generate", saveHandler)
	http.ListenAndServe(":8080", nil)
}
