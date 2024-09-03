package main

import (
	"context"
	"flag"
	"html/template"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
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
	e.Logger.SetLevel(log.DEBUG)

	//middleware
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("35K"))
	e.Use(middleware.Secure())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: 5 * time.Second}))
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: `{"time":"${time_rfc3339_nano}",` +
		`"remote_ip":"${remote_ip}",` +
		`"host":"${host}",` +
		`"method":"${method}",` +
		`"uri":"${uri}",` +
		`"status":${status},` +
		`"error":"${error}",` +
		`"latency_human":"${latency_human}"` +
		`}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000"}))

	//flag parsing
	port := flag.String("port", "4000", "port for app")
	flag.Parse()

	//static files
	e.Static("/static", "static")

	//initialise render
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}

	e.Renderer = renderer
	// routes
	e.GET("/", listContacts)

	//handle gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//start server in a different thread
	go func() {
		if err := e.Start(":" + *port); err != http.ErrServerClosed {
			e.Logger.Fatal("Could not start the server. Shutting down")
		}
	}()
	//block untill the context channel is closed (maybe due to os.Interrupt signal)
	<-ctx.Done()

	//buy 10 seconds to gracefully shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal("Could not not shutdown the server gracefully, ", err)
	}

}

//handlers

func listContacts(c echo.Context) error {
	return c.Render(http.StatusOK, "contacts.html", contacts)
}
