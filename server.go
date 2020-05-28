package main

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	e.Renderer = renderer

	e.GET("/", handleIndex)
	e.POST("/", handleIndex)

	e.GET("/monster/:name", getIdenticon)

	e.Logger.Fatal(e.Start("0.0.0.0:5000"))
}

func getIdenticon(c echo.Context) error {
	name := c.Param("name")
	resp, _ := http.Get(`http://dnmonster:8080/monster/` + name + `?size=80`)
	return c.Stream(http.StatusOK, "image/png", resp.Body)
}

func handleIndex(c echo.Context) error {
	salt := "UNIQUE"
	h := sha256.New()
	saltedName := salt + "Joe malone"
	if c.FormValue("name") != "" {
		saltedName = salt + c.FormValue("name")
	}
	h.Write([]byte(saltedName))
	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"name":      c.FormValue("name"),
		"name_hash": hex.EncodeToString(h.Sum(nil)),
	})
}
