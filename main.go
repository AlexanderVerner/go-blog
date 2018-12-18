package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./db/documents"
	"./models"
	"./utils"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
	"github.com/tobyzxj/mgo"
)

var postsCollection *mgo.Collection

func getLoginHandler(rnd render.Render) {
	rnd.HTML(200, "login", nil)
}

func IndexHandler(rnd render.Render) {
	postDocuments := []documents.PostDocument{}
	postsCollection.Find(nil).All(&postDocuments)
	posts := []models.Post{}
	for _, doc := range postDocuments {
		post := models.Post{doc.Id, doc.Title, doc.ContentHtml, doc.ContentMarkdown}
		posts = append(posts, post)
	}

	fmt.Println(posts)

	rnd.HTML(200, "index", posts)
}

func writeHandler(rnd render.Render) {
	post := models.Post{}
	rnd.HTML(200, "write", post)
}

func editHandler(rnd render.Render, r *http.Request, params martini.Params) {

	id := params["id"]
	postDocument := documents.PostDocument{}
	err := postsCollection.FindId(id).One(&postDocument)
	if err != nil {
		rnd.Redirect("/")
	}
	post := models.Post{postDocument.Id, postDocument.Title, postDocument.ContentHtml, postDocument.ContentMarkdown}

	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	contentMarkdown := r.FormValue("content")
	contentHtml := utils.ConvertMarkdownToHtml(contentMarkdown)

	postDocument := documents.PostDocument{id, title, contentHtml, contentMarkdown}
	if id != "" {
		postsCollection.UpdateId(id, postDocument)
	} else {
		id = utils.GenerateId()
		postDocument.Id = id
		postsCollection.Insert(postDocument)
	}

	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, r *http.Request, params martini.Params) {
	id := params["id"]
	if id == "" {
		rnd.Redirect("/")
	}

	postsCollection.RemoveId(id)

	rnd.Redirect("/")
}

func getHtmlHandler(rnd render.Render, r *http.Request) {
	md := r.FormValue("md")
	html := utils.ConvertMarkdownToHtml(md)

	rnd.JSON(200, map[string]interface{}{"html": html})
}

func unescape(x string) interface{} {
	return template.HTML(x)
}

func main() {
	fmt.Println("Listening on port :3000")

	sessions, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	postsCollection = sessions.DB("blog").C("posts")

	m := martini.Classic()

	unescapeFuncMap := template.FuncMap{"unescape": unescape}

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                         // Specify what path to load the templates from.
		Layout:     "layout",                            // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"},          // Specify extensions to load for templates.
		Funcs:      []template.FuncMap{unescapeFuncMap}, // Specify helper function maps for templates to access.
		Charset:    "UTF-8",                             // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,                                // Output human readable JSON
		IndentXML:  true,                                // Output human readable XML
	}))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", IndexHandler)
	m.Get("/login", getLoginHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit/:id", editHandler)
	m.Get("/delete/:id", deleteHandler)
	m.Post("/SavePost", savePostHandler)
	m.Post("/write", getHtmlHandler)
	m.Run()
}
