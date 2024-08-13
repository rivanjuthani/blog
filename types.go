package main

import (
	"html/template"
	"time"
)

type Post struct {
	Title      string `toml:"title"`
	Slug       string `toml:"slug"`
	Content    template.HTML
	Date       time.Time `toml:"date"`
	Tags       []string  `toml:"tags"`
	StringDate string
	Author     Author `toml:"author"`
}

type Author struct {
	Name  string `toml:"name"`
	Email string `toml:"email"`
}

type PostContainer struct {
	Posts []Post
}
