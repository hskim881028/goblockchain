package explorer

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/hskim881028/goblockchain/blockchain"
)

var templates *template.Template

const (
	port        string = ":4000"
	templateDir string = "explorer/templates/"
)

type homeData struct {
	PageTitle string
	Blocks    []*blockchain.Block
}

func home(rw http.ResponseWriter, r *http.Request) {
	data := homeData{"home", blockchain.GetBlcokchain().AllBlocks()}
	templates.ExecuteTemplate(rw, "home", data)
}

func add(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		templates.ExecuteTemplate(rw, "add", nil)
	case "POST":
		r.ParseForm()
		data := r.Form.Get("blockData")
		blockchain.GetBlcokchain().AddBlock(data)
		http.Redirect(rw, r, "/", http.StatusPermanentRedirect)
	}
}

func Start() {
	templates = template.Must(template.ParseGlob(templateDir + "pages/*.gohtml"))
	templates = template.Must(templates.ParseGlob(templateDir + "partials/*.gohtml"))

	http.HandleFunc("/", home)
	http.HandleFunc("/add", add)
	fmt.Printf("listening on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
