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
	"github.com/samoei/cms/models"
)

type application struct {
	contacts []models.Contact
}

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	app := &application{
		contacts: []models.Contact{
			{ID: 1, FirstName: "Philemon", LastName: "Samoei", Email: "phil.samoei@gmail.com", Phone: "0722511046", City: "Eldoret", DOB: "12/12/1992"},
			{ID: 2, FirstName: "James", LastName: "Mwangi", Email: "j.mwangi@gmail.com", Phone: "0724541046", City: "Kiambu", DOB: "19/08/1996"},
			{ID: 3, FirstName: "Alphas", LastName: "Koech", Email: "aphas3467@gmail.com", Phone: "0732511078", City: "Kapsabet", DOB: "01/03/1999"},
			{ID: 4, FirstName: "Maryann", LastName: "Mwangi", Email: "mmwangi77@gmail.com", Phone: "0722420046", City: "Kikuyu", DOB: "07/07/1999"},
			{ID: 5, FirstName: "Carol", LastName: "Kihara", Email: "carol.kihara99@gmail.com", Phone: "0723611302", City: "Limuru", DOB: "12/12/1992"},
			{ID: 6, FirstName: "Mary", LastName: "Wamaitha", Email: "mary.mesh@gmail.com", Phone: "0722511046", City: "Eldoret", DOB: "12/12/1992"},
			{ID: 7, FirstName: "Maggie", LastName: "Maina", Email: "maina.maggz@gmail.com", Phone: "0722511046", City: "Nyahurur", DOB: "12/12/1992"},
			{ID: 8, FirstName: "Angela", LastName: "Wendo", Email: "wndanglalove@gmail.com", Phone: "0722511046", City: "Mombasa", DOB: "12/12/1992"},
			{ID: 9, FirstName: "Peris", LastName: "Chepkoech", Email: "pkoech77@gmail.com", Phone: "0722511046", City: "Moi's Bridge", DOB: "12/12/1992"},
			{ID: 10, FirstName: "Susan", LastName: "Githaiga", Email: "susie-mwaks@gmail.com", Phone: "0722511046", City: "Webuye", DOB: "12/12/1992"},
			{ID: 11, FirstName: "Clinton", LastName: "Mwale", Email: "clinto.mwenyewe@gmail.com", Phone: "0722511046", City: "Machakos", DOB: "12/12/1992"},
			{ID: 12, FirstName: "Karua", LastName: "Kihara", Email: "mrkarua.kihara@gmail.com", Phone: "0722511046", City: "Kisii", DOB: "12/12/1992"},
		},
	}
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
	e.GET("/", app.listContacts)

	//handle gracefull shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	//start server in a different thread
	go func() {
		if err := e.Start(":" + *port); err != http.ErrServerClosed {
			e.Logger.Fatal("Could not start the server. Shutting down,", err)
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

func (a *application) listContacts(c echo.Context) error {
	return c.Render(http.StatusOK, "contacts.html", a.contacts)
}
