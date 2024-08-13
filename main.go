package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/adrg/frontmatter"
	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/patrickmn/go-cache"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
)

var tpl *template.Template
var tplHome *template.Template

var c = cache.New(5*time.Minute, 10*time.Minute)

func init() {
	var err error
	tpl, err = template.ParseFiles("./templates/post.gohtml")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
	tplHome, err = template.ParseFiles("./templates/home.gohtml")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", HomeHandler())
	mux.HandleFunc("GET /posts/{slug}", PostHandler(FileReader{}))

	// err := http.ListenAndServe(":"+strconv.Itoa(HTTP_PORT), mux)
	// if err != nil {
	// 	log.Fatal("ListenAndServe: ", err)
	// }

	err := http.ListenAndServeTLS(":443", "./certs/server.crt", "./certs/server.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func getPosts() PostContainer {
	var pageData PostContainer

	if x, found := c.Get("posts"); found {
		posts := x.(*PostContainer)
		return *posts
	} else {
		log.Println("cache unavailable, fetching new data...")
	}

	entries, err := os.ReadDir("./posts")
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range entries {
		var post Post
		slug := strings.Split(e.Name(), ".")[0]
		postMarkdown, err := FileReader{}.Read(slug)
		if err != nil {
			return pageData
		}
		frontmatter.Parse(strings.NewReader(postMarkdown), &post)
		post.Slug = slug
		post.StringDate = post.Date.Format("Jan 02 2006")
		pageData.Posts = append(pageData.Posts, post)
	}

	sort.Slice(pageData.Posts, func(i, j int) bool {
		return pageData.Posts[i].Date.Unix() > pageData.Posts[j].Date.Unix()
	})

	c.Set("posts", &pageData, cache.DefaultExpiration)

	return pageData
}

func HomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var pageData = getPosts()
		err := tplHome.Execute(w, pageData)
		if err != nil {
			log.Println(err)
			http.Error(w, "Error populating template", http.StatusInternalServerError)
			return
		}
	}
}

type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(POSTS_FOLDER + slug + ".md")
	if err != nil {
		return "", err
	}
	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func PostHandler(sl SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post Post
		post.Slug = r.PathValue("slug")
		postMarkdown, err := sl.Read(post.Slug)
		if err != nil {
			// TODO: Handle different errors in the future
			http.Error(w, "Post not found", http.StatusNotFound)
			return
		}
		rest, err := frontmatter.Parse(strings.NewReader(postMarkdown), &post)
		if err != nil {
			http.Error(w, "Error parsing frontmatter", http.StatusInternalServerError)
			return
		}
		post.StringDate = post.Date.Format("Jan 02 2006")
		mdRenderer := goldmark.New(
			goldmark.WithExtensions(
				highlighting.NewHighlighting(
					highlighting.WithStyle("xcode-dark"),
					highlighting.WithFormatOptions(
						chromahtml.WithLineNumbers(true),
					),
				),
			),
		)
		var buf bytes.Buffer
		err = mdRenderer.Convert(rest, &buf)
		if err != nil {
			http.Error(w, "Error converting markdown", http.StatusInternalServerError)
			return
		}
		post.Content = template.HTML(buf.String())
		err = tpl.Execute(w, post)
		if err != nil {
			http.Error(w, "Error populating template", http.StatusInternalServerError)
			return
		}
	}
}
