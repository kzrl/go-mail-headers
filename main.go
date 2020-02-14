package main

import (
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/kzrl/go-mail-headers/headers"
)

func main() {

	port := 8082
	r := mux.NewRouter()

	r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	r.HandleFunc("/", rootHandler)
	r.HandleFunc("/headers", formHandler)

	log.Println("*********")
	log.Printf("gomailheaders listening on port %d\n", port)
	log.Println("*********")
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)

	/*body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", body)*/

}

type Page struct {
	Header     mail.Header
	Hops       []headers.Hop
	TotalDelay int
}

// Render the form
func rootHandler(w http.ResponseWriter, req *http.Request) {
	t := template.Must(template.ParseFiles("templates/form.html"))
	err := t.Execute(w, nil)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}

func formHandler(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		panic(err) //TODO
	}

	msg := req.Form.Get("headers")
	r := strings.NewReader(msg)
	m, err := mail.ReadMessage(r)
	if err != nil {
		log.Fatal(err)
	}

	//	header := m.Header
	/*fmt.Println("Date:", header.Get("Date"))
	fmt.Println("From:", header.Get("From"))
	fmt.Println("To:", header.Get("To"))
	fmt.Println("Subject:", header.Get("Subject"))
	fmt.Println("Atl-Mail-ID:", header.Get("Atl-Mail-ID"))
	*/

	var p Page
	p.Header = m.Header
	p.Hops = headers.Hops(m.Header["Received"])

	/*for i, hop := range p.Hops {
		fmt.Printf("%d %s\n", i, hop)
	}*/

	t := template.Must(template.ParseFiles("templates/results.html"))
	err = t.Execute(w, p)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
