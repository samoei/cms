package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Contact struct {
	ID    int
	Name  string
	Email string
	Phone string
}

var contacts = []Contact{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Phone: "123-456-7890"},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Phone: "7890-456-123"},
	{ID: 3, Name: "Phil Cole", Email: "phil@example.com", Phone: "456-123-7890"},
}

var contactMap = make(map[int]Contact)

func init() {
	for _, contact := range contacts {
		contactMap[contact.ID] = contact
	}
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	//middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//static files
	e.Static("/static", "static")

	//initialise render
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	templ, err := template.ParseGlob("templates/*.html")

	if err != nil {
		log.Fatal("Error in parsing templates: %v", err)
	}

	for _, t := range templ.Templates() {
		fmt.Println("Template Name:", t.Name())
	}

	e.Renderer = renderer

	// routes
	e.GET("/", listContacts)

	//start server
	e.Logger.Fatal(e.Start(":8080"))
}

//handlers

func listContacts(c echo.Context) error {
	return c.Render(http.StatusOK, "contacts.html", contacts)
}
