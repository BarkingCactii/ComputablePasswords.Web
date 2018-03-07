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

/*
func applicationHandler(w http.ResponseWriter, r *http.Request, title string) {
	application, err := loadApplication(title)
	if err == nil {
		// load existing application
		renderTemplateApplication(w, "application", application)
		return
	}

	a := &Application{Title: title, Position: "", Company: "", DateApplied: "", Followup: "", Action: "", Notes: []byte(""), NotesRows: 5}
	renderTemplateApplication(w, "application", a)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	company := r.FormValue("company")
	contact := r.FormValue("contact")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	logo := r.FormValue("logo")
	body := r.FormValue("body")
	p := &Contact{Title: title, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
*/

/*
func PopulateTemplates() *template.Template {
	result := template.New("templates")
	const BasePath = "templates"
	template.Must(result.ParseGlob(BasePath + "/*.html"))
	return result
}
*/
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
		fmt.Println(fi.Name())
	}
	return result
}

func saveHandler(w http.ResponseWriter, r *http.Request) { //, title string) {

	context.WebURL = r.FormValue("url")
	fmt.Printf("InputM = [%s]\n", context.WebURL)
	context.Result = "Password: " + calc(context.WebURL)
	//p := &Contact{Title: title, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: []byte(body)}
	//err := p.save()
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	http.Redirect(w, r, "/main", http.StatusFound)
	/*
		company := r.FormValue("company")
		contact := r.FormValue("contact")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		logo := r.FormValue("logo")
		body := r.FormValue("body")
		p := &Contact{Title: title, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: []byte(body)}
		err := p.save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusFound)
	*/
}

var context = viewmodel.NewBase()

func main() {
	templates := PopulateTemplates()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requestedFile := r.URL.Path[1:]
		template := templates[requestedFile+".html"]
		//context = viewmodel.NewBase()
		/*
			var context interface{}

			switch requestedFile {
			case "shop":
				context = viewmodel.NewShop()
			default:
				context = viewmodel.NewHome()
			}
		*/
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
	http.HandleFunc("/generate/", saveHandler)
	http.ListenAndServe(":8080", nil)

	//	http.HandleFunc("/img/", serveResource)
	//	http.HandleFunc("/css/", serveResource)
	//http.HandleFunc("/", makeHandler(listHandler))
	/*
		templates := PopulateTemplates()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			requestedFile := r.URL.Path[1:]
			t := templates[requestedFile+".html"]
			//	log.Println(requestedFile)
			if t != nil {
				fmt.Println("hi")
				err := t.Execute(w, nil)
				if err != nil {
					//	log.Println(err)
				}
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
	*/
	/*
		f, err := os.Open("public" + r.URL.Path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}

		defer f.Close()
		io.Copy(w, f)
	*/
	//	})

	//	http.Handle("/img/", http.FileServer(http.Dir("public")))
	//	http.Handle("/css/", http.FileServer(http.Dir("public")))

	//	http.ListenAndServe(":8080", nil)
}
