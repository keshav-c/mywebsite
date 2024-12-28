package main

import (
	"fmt"
	"os"
	"strings"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := "./pages/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := "./pages/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func (p *Page) Validate() error {
	title := strings.TrimSpace(p.Title)
	if len(title) == 0 {
		return fmt.Errorf("Page Title cannot be empty")
	}
	body := strings.TrimSpace(string(p.Body))
	if len(body) == 0 {
		return fmt.Errorf("Page Body cannot be empty")
	}
	return nil
}
