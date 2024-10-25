package main

import (
	"html/template"
	"io"
	"net/http"
	"os"

	"creeston/lists/internal/handlers"
	"creeston/lists/internal/repository"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func UserIdCookieHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := ""
		userIdCookie, error := c.Cookie("wishlist_uid")
		if error != nil {
			if error == http.ErrNoCookie {
				userId = uuid.New().String()
				c.SetCookie(&http.Cookie{
					Name:  "wishlist_uid",
					Value: userId,
				})
			}
		} else {
			userId = userIdCookie.Value
		}

		c.Set("userId", userId)
		return next(c)
	}
}

func main() {
	godotenv.Load()
	e := echo.New()
	data := repository.NewData()
	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	e.Use(middleware.Logger())
	e.Use(UserIdCookieHandler)
	e.Renderer = NewTemplate()
	handlers.SetupRoutes(e, data, baseUrl)
	e.Static("/css", "css")
	e.Logger.Fatal(e.Start(":" + port))
}
