package main

import (
	"html/template"
	"io"
	"sort"
	"strconv"

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

type Count struct {
	Count int
}

type WishlistItem struct {
	Text  string
	Index int
}

type Wishlist struct {
	Items []WishlistItem
	Id    int
}

type WishlistFormData struct {
	Items    []WishlistItem
	HasItems bool
	HasId    bool
	Id       int
}

func NewWishlist(items []WishlistItem, id int) Wishlist {
	return Wishlist{
		Items: items,
		Id:    id,
	}
}

type Wishlists = []Wishlist

type Data struct {
	Wishlists Wishlists
}

func newData() *Data {
	return &Data{
		Wishlists: Wishlists{},
	}
}

func main() {
	data := newData()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", WishlistFormData{HasItems: false, HasId: false})
	})

	e.POST("/wishlist", func(c echo.Context) error {
		items := []WishlistItem{}
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		for key, value := range params {
			if key[:4] != "item" {
				continue
			}

			if len(value) == 0 {
				continue
			}

			text := value[0]
			if text == "" {
				continue
			}

			indexStringValue := key[4:]
			index, error := strconv.Atoi(indexStringValue)
			if error != nil {
				return error
			}

			items = append(items, WishlistItem{Text: value[0], Index: index})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Index < items[j].Index
		})

		wishlistId := len(data.Wishlists)
		wishlist := NewWishlist(items, wishlistId)
		data.Wishlists = append(data.Wishlists, wishlist)

		// redirect to the wishlist page
		// by setting hx-redirect header in response
		// c.Response().Header().Set("hx-redirect", "/wishlist/"+strconv.Itoa(wishlistId))
		// return c.Redirect(302, "/wishlist/"+strconv.Itoa(wishlistId))
		return c.Render(200, "wishlist-form", WishlistFormData{
			Items:    items,
			HasItems: true,
			HasId:    true,
			Id:       wishlistId,
		})
	})

	e.PUT("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		if (id < 0) || (id >= len(data.Wishlists)) {
			c.Response().Header().Set("hx-redirect", "/")
			return c.String(404, "Wishlist not found")
		}

		items := []WishlistItem{}
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		for key, value := range params {
			if key[:4] != "item" {
				continue
			}

			if len(value) == 0 {
				continue
			}

			text := value[0]
			if text == "" {
				continue
			}

			indexStringValue := key[4:]
			index, error := strconv.Atoi(indexStringValue)
			if error != nil {
				return error
			}

			items = append(items, WishlistItem{Text: value[0], Index: index})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Index < items[j].Index
		})

		data.Wishlists[id].Items = items

		return c.Render(200, "wishlist-form", WishlistFormData{
			Items:    items,
			HasItems: true,
			HasId:    true,
			Id:       id,
		})
	})

	e.GET("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		if (id < 0) || (id >= len(data.Wishlists)) {
			c.Response().Header().Set("hx-redirect", "/")
			return c.String(404, "Wishlist not found")
		}

		wishlist := data.Wishlists[id]
		return c.Render(200, "wishlist", wishlist)
	})

	e.Static("/css", "css")
	e.Static("/js", "js")

	e.Logger.Fatal(e.Start(":1323"))
}
