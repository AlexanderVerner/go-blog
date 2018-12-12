package main

import (
	"fmt"
	"net/http"

	"./models"
	"./utils"

	"github.com/codegangsta/martini"
	"github.com/martini-contrib/render"
)

var posts map[string]*models.Post

func IndexHandler(rnd render.Render) {
	/*t, err := template.ParseFiles("templates/index.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}*/

	fmt.Println(posts)

	rnd.HTML(200, "index", posts)

	//t.ExecuteTemplate(w, "index", posts)
}

func writeHandler(rnd render.Render) {
	/*t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	t.ExecuteTemplate(w, "write", nil)*/

	rnd.HTML(200, "write", nil)
}

func editHandler(rnd render.Render, r *http.Request) {
	/*t, err := template.ParseFiles("templates/write.html", "templates/header.html", "templates/footer.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}*/
	id := r.FormValue("id")
	post, found := posts[id]
	if !found {
		rnd.Redirect("/")
	}

	//t.ExecuteTemplate(w, "write", post)

	rnd.HTML(200, "write", post)
}

func savePostHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	content := r.FormValue("content")

	var post *models.Post
	if id != "" {
		post = posts[id]
		post.Title = title
		post.Content = content
	} else {
		id = utils.GenerateId()
		post := models.NewPost(id, title, content)
		posts[post.Id] = post
	}

	//http.Redirect(w, r, "/", 302)
	rnd.Redirect("/")
}

func deleteHandler(rnd render.Render, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		//http.NotFound(w, r)
		rnd.Redirect("/")
	}

	delete(posts, id)

	//http.Redirect(w, r, "/", 302)
	rnd.Redirect("/")
}

func main() {
	fmt.Println("Listening on port :3000")

	posts = make(map[string]*models.Post, 0)

	m := martini.Classic()

	m.Use(render.Renderer(render.Options{
		Directory:  "templates",                // Specify what path to load the templates from.
		Layout:     "layout",                   // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		//Funcs:           []template.FuncMap{AppHelpers}, // Specify helper function maps for templates to access.
		//Delims:          render.Delims{"{[{", "}]}"},    // Sets delimiters to the specified strings.
		Charset:    "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		IndentJSON: true,    // Output human readable JSON
		IndentXML:  true,    // Output human readable XML
		//HTMLContentType: "application/xhtml+xml",        // Output XHTML content type instead of default "text/html"
	}))

	//http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	staticOptions := martini.StaticOptions{Prefix: "assets"}
	m.Use(martini.Static("assets", staticOptions))
	m.Get("/", IndexHandler)
	m.Get("/write", writeHandler)
	m.Get("/edit", editHandler)
	m.Get("/delete", deleteHandler)
	m.Post("/SavePost", savePostHandler)

	//http.ListenAndServe(":3000", nil)
	m.Run()
}
