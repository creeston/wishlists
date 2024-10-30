package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"creeston/lists/internal/handlers"
	"creeston/lists/internal/repository"
	"creeston/lists/internal/utils"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/time/rate"
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
				userId = utils.GenerateUUID()
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

func getLanguageFromRequest(c echo.Context) language.Tag {
	cookieLang, err := c.Cookie("lang")
	if err == nil {
		cookieValue := cookieLang.Value
		lang, err := language.Parse(cookieValue)
		if err == nil {
			return lang
		}
	}

	acceptLang := c.Request().Header.Get("Accept-Language")
	acceptLang = acceptLang[:5]
	lang := message.MatchLanguage(acceptLang)
	return lang
	// lang, err := language.Parse(acceptLang)
	// if err == nil {
	// 	return lang
	// }

	// return language.English
}

func LanguageHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		lang := getLanguageFromRequest(c)
		p := message.NewPrinter(lang)
		c.Set("clientLanguage", lang.String()[:5])
		c.Set("i18n", p)
		return next(c)
	}
}

func main() {
	godotenv.Load()
	e := echo.New()
	repository := repository.NewInMemoryRepository()
	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	maxItemsCountValue := os.Getenv("MAX_ITEMS_COUNT")
	maxItemLengthValue := os.Getenv("MAX_ITEM_LENGTH")
	maxBodySizeValue := os.Getenv("MAX_BODY_SIZE")
	maxRateLimitValue := os.Getenv("MAX_RATE_LIMIT")

	maxItemsCount, err := strconv.Atoi(maxItemsCountValue)
	if err != nil {
		panic(err)
	}

	maxItemLength, err := strconv.Atoi(maxItemLengthValue)
	if err != nil {
		panic(err)
	}

	maxRateLimit, err := strconv.Atoi(maxRateLimitValue)
	if err != nil {
		panic(err)
	}

	validationConfig := handlers.ValidationConfig{
		MaxItemsCount: maxItemsCount,
		MaxItemLength: maxItemLength,
	}

	e.Use(middleware.BodyLimit(maxBodySizeValue))
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(rate.Limit(maxRateLimit))))
	e.Use(middleware.Logger())
	e.Use(UserIdCookieHandler)
	e.Use(LanguageHandler)
	e.Renderer = NewTemplate()
	handlers.SetupRoutes(e, repository, baseUrl, validationConfig)
	e.Static("/css", "css")
	e.Logger.Fatal(e.Start(":" + port))
}
