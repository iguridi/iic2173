package main

import (
    "fmt"
	"io/ioutil"
    "net/http"
    "log"
    "html/template"
)

type Page struct {
    Title string
    Body  []byte
}

func (p *Page) save() error {
    filename := p.Title
    return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
    filename := title
    body, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    return &Page{Title: title, Body: body}, nil
}

// func main() {
//     p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
//     p1.save()
//     p2, _ := loadPage("TestPage")
//     fmt.Println(string(p2.Body))
// }



// func viewHandler(w http.ResponseWriter, r *http.Request) {
//     title := r.URL.Path[len("/view/"):]
//     p, _ := loadPage(title)
//     fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
// }


func readFile(title string) (*Page) {
    p, err := loadPage(title)
    if err != nil {
        fmt.Println("Err: %v", err)
        p = &Page{Title: "Error:" , Body: []byte("File could not be read")}
    } 
    return p
}

func index(w http.ResponseWriter, r *http.Request) {
    title := "comments"
    switch r.Method {
        case "GET":     
            p, err := loadPage(title)
            if err != nil {
                p = &Page{Title: title}
            }
            t, _ := template.ParseFiles("index.html")
            t.Execute(w, p)
        case "POST":
            if err := r.ParseForm(); err != nil {
                fmt.Fprintf(w, "ParseForm() err: %v", err)
                return
            }
            extra := r.FormValue("comment")
            p2 := readFile(title)
            body := string(p2.Body) + "\n" + extra
            // fmt.Println(string(body))
            p := &Page{Title: title, Body: []byte(body)}
            p.save()
            http.Redirect(w, r, "/", http.StatusFound)
        default:
            fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
    }

    
}

// func handler(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

func saveHandler(w http.ResponseWriter, r *http.Request) {
    // title := r.URL.Path[len("/new"):]
    body := r.FormValue("comment")
    p := &Page{Title: "comments", Body: []byte(body)}
    p.save()
    // index(w, r)
    // http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
    http.HandleFunc("/", index)
    log.Fatal(http.ListenAndServe(":8080", nil))
    // http.HandleFunc("/new", saveHandler)
    // log.Fatal(http.ListenAndServe(":8080", nil))
}