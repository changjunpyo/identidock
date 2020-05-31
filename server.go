package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/labstack/echo/v4"
	"html/template"
	"io"
	"net/http"
	"strings"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func redisInjection(next echo.HandlerFunc,redis *redis.Client ) echo.HandlerFunc{
	return func(c echo.Context) error{
		c.Set("redis",redis )
		return next(c)
	}
 }
var client *redis.Client
func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	}
	client = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
		Password: "",
		DB :0,
	})

	e.Renderer = renderer

	e.GET("/", handleIndex)
	e.POST("/", handleIndex)

	e.GET("/monster/:name", getIdenticon)

	e.Logger.Fatal(e.Start("0.0.0.0:5000"))
}

func getIdenticon(c echo.Context) error {
	var ret io.Reader
	name := c.Param("name")
	image, err := client.Get(name).Result()
	if err == redis.Nil {
		resp, _ := http.Get(`http://dnmonster:8080/monster/` + name + `?size=80`)
		buf := new(strings.Builder)
		_, err := io.Copy(buf, resp.Body)
		if err != nil {
			fmt.Printf("error in string to io.Reader")
		}
		if err := client.Set(name, buf.String(), 0).Err(); err != nil {
			return err
		}
		ret = strings.NewReader(buf.String())
	} else if err != nil{
		panic(err)
	} else{
		ret = strings.NewReader(image)
	}
	return c.Stream(http.StatusOK, "image/png", ret)
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
