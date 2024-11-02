package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

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
	baseUrl := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")
	maxItemsCountValue := os.Getenv("MAX_ITEMS_COUNT")
	maxItemLengthValue := os.Getenv("MAX_ITEM_LENGTH")
	MaxWishlistsPerDayValue := os.Getenv("MAX_WISHLISTS_PER_DAY")
	maxBodySizeValue := os.Getenv("MAX_BODY_SIZE")
	useInMemoryDb := os.Getenv("USE_IN_MEMORY_DB")

	maxItemsCount, err := strconv.Atoi(maxItemsCountValue)
	if err != nil {
		panic(err)
	}

	maxItemLength, err := strconv.Atoi(maxItemLengthValue)
	if err != nil {
		panic(err)
	}

	maxWishlistsPerDay, err := strconv.Atoi(MaxWishlistsPerDayValue)
	if err != nil {
		panic(err)
	}

	validationConfig := handlers.ValidationConfig{
		MaxItemsCount:      maxItemsCount,
		MaxItemLength:      maxItemLength,
		MaxWishlistsPerDay: maxWishlistsPerDay,
	}

	e.Use(middleware.BodyLimit(maxBodySizeValue))
	e.Use(CreateGlobalRateLimiter())
	e.Use(middleware.Logger())
	e.Use(UserIdCookieHandler)
	e.Use(LanguageHandler)
	e.Renderer = NewTemplate()

	var dataRepository repository.WishlistRepository
	if strings.ToLower(useInMemoryDb) == "true" {
		dataRepository = repository.NewInMemorySqliteRepository()
	} else {
		dbPath := os.Getenv("SQLITE_DB_NAME")
		dataRepository = repository.NewSqliteRepository(dbPath)
	}
	handlers.SetupRoutes(e, dataRepository, baseUrl, validationConfig)
	e.Static("/css", "static/css")
	e.Static("/icons", "static/icons")
	e.Logger.Fatal(e.Start(":" + port))
}

func CreateGlobalRateLimiter() echo.MiddlewareFunc {
	return middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(5), Burst: 5, ExpiresIn: 3 * time.Minute},
	))
}
