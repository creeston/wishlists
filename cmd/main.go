package main

import (
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/google/uuid"
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

func ParseWishlistFormDataToWishlistItems(data map[string][]string) []*WishlistItem {
	items := []*WishlistItem{}
	for key, value := range data {
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
			return nil
		}

		items = append(items, &WishlistItem{Text: value[0], Index: index})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Index < items[j].Index
	})

	return items
}

func main() {
	data := NewData()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(UserIdCookieHandler)
	e.Renderer = NewTemplate()

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", WishlistFormData{HasItems: false, HasId: false})
	})

	e.POST("/wishlist", func(c echo.Context) error {
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		items := ParseWishlistFormDataToWishlistItems(params)
		userId := c.Get("userId").(string)
		wishlist := data.AddWishlist(userId, items)
		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(wishlist))
	})

	e.PUT("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		wishlist := data.GetWishlistByIdOrNull(id)

		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		userId := c.Get("userId").(string)
		if userId != wishlist.CreatorId {
			return c.String(403, "Forbidden")
		}

		params, error := c.FormParams()
		if error != nil {
			return error
		}

		items := ParseWishlistFormDataToWishlistItems(params)
		data.UpdateWishlistWithItems(id, items)
		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(wishlist))
	})

	e.PUT("/wishlist/:id/:itemId", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		itemId, error := strconv.Atoi(c.Param("itemId"))
		if error != nil {
			return error
		}

		wishlist := data.GetWishlistByIdOrNull(id)
		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		if (itemId < 0) || (itemId >= len(wishlist.Items)) {
			return c.String(404, "Item not found")
		}

		checkRequest := c.FormValue(("flag")) == "on"
		userId := c.Get("userId").(string)
		wishlistItem := wishlist.Items[itemId]
		formData := WishlistCheckedItemData{
			Index: wishlistItem.Index,
			Text:  wishlistItem.Text,
			Id:    wishlist.Id,
		}

		if checkRequest && wishlistItem.Checked && wishlistItem.CheckedById == userId {
			return c.Render(200, "wishlist-checked-item", formData)
		}

		if checkRequest && wishlistItem.Checked && wishlistItem.CheckedById != userId {
			return c.Render(200, "wishlist-already-checked-item-with-popup", formData)
		}

		if !checkRequest && !wishlistItem.Checked {
			return c.Render(200, "wishlist-not-checked-item", formData)
		}

		if !checkRequest && wishlistItem.Checked && wishlistItem.CheckedById != userId {
			return c.Render(200, "wishlist-already-checked-item", formData)
		}

		if !checkRequest && wishlistItem.Checked && wishlistItem.CheckedById == userId {
			wishlistItem.Checked = false
			wishlistItem.CheckedById = ""
			return c.Render(200, "wishlist-not-checked-item", formData)
		}

		wishlistItem.Checked = true
		wishlistItem.CheckedById = userId
		return c.Render(200, "wishlist-checked-item", formData)
	})

	e.GET("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		if (id < 0) || (id >= len(data.Wishlists)) {
			return c.Render(200, "not-found", nil)
		}

		wishlist := data.Wishlists[id]
		userId := c.Get("userId").(string)
		if userId != wishlist.CreatorId {
			return c.Render(200, "wishlist", MapWishlistToWishlistViewFormData(wishlist, userId))
		}

		return c.Render(200, "index", MapWishlistToWishlistFormData(wishlist))
	})

	e.Static("/css", "css")
	e.Static("/js", "js")

	e.Logger.Fatal(e.Start(":1323"))
}
