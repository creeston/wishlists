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

type Count struct {
	Count int
}

type WishlistItem struct {
	Text    string
	Checked bool
	Index   int
}

type Wishlist struct {
	Items  []*WishlistItem
	Id     int
	UserId string
}

type WishlistFormData struct {
	Items    []WishlistFormItem
	HasItems bool
	HasId    bool
	Id       int
}

type WishlistFormItem struct {
	Index int
	Text  string
}

type WishlistViewFormData struct {
	Items []WishlistCheckedItemData
	Id    int
}

type WishlistCheckedItemData struct {
	Index   int
	Text    string
	Id      int
	Checked bool
}

func MapWishlistToWishlistFormData(wishlist *Wishlist) WishlistFormData {
	items := []WishlistFormItem{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistFormItem{
			Index: item.Index,
			Text:  item.Text,
		})
	}

	return WishlistFormData{
		Items:    items,
		HasItems: true,
		HasId:    true,
		Id:       wishlist.Id,
	}
}

func MapWishlistToWishlistViewFormData(wishlist *Wishlist) WishlistViewFormData {
	items := []WishlistCheckedItemData{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistCheckedItemData{
			Index:   item.Index,
			Text:    item.Text,
			Id:      wishlist.Id,
			Checked: item.Checked,
		})
	}

	return WishlistViewFormData{
		Items: items,
		Id:    wishlist.Id,
	}
}

func NewWishlist(items []*WishlistItem, id int, userId string) *Wishlist {
	return &Wishlist{
		Items:  items,
		Id:     id,
		UserId: userId,
	}
}

type Wishlists = []*Wishlist

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

	data.Wishlists = append(data.Wishlists, NewWishlist(
		[]*WishlistItem{
			{Text: "Cake", Index: 0},
			{Text: "Candles", Index: 1},
			{Text: "Balloons", Index: 2},
			{Text: "Presents. A lot a lot a lof a very long list of presents please!", Index: 3},
		},
		0,
		"default"))

	e.POST("/wishlist", func(c echo.Context) error {
		items := []*WishlistItem{}
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		userId := ""
		userIdCookie, error := c.Cookie("wishlist_uid")
		if error != nil {
			if error != http.ErrNoCookie {
				return error
			} else {
				userId = uuid.New().String()
			}
		} else {
			userId = userIdCookie.Value
		}

		c.SetCookie(&http.Cookie{
			Name:  "wishlist_uid",
			Value: userId,
		})

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

			items = append(items, &WishlistItem{Text: value[0], Index: index})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Index < items[j].Index
		})

		wishlistId := len(data.Wishlists)
		wishlist := NewWishlist(items, wishlistId, userId)
		data.Wishlists = append(data.Wishlists, wishlist)

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

		if (id < 0) || (id >= len(data.Wishlists)) {
			return c.String(404, "Wishlist not found")
		}

		wishlist := data.Wishlists[id]
		if (itemId < 0) || (itemId >= len(wishlist.Items)) {
			return c.String(404, "Item not found")
		}

		wishlistItem := wishlist.Items[itemId]

		isChecked := c.FormValue(("flag")) == "on"
		wishlistItem.Checked = isChecked
		formData := WishlistCheckedItemData{
			Index: wishlistItem.Index,
			Text:  wishlistItem.Text,
			Id:    wishlist.Id,
		}

		if isChecked {
			return c.Render(200, "wishlist-checked-item", formData)
		} else {
			return c.Render(200, "wishlist-not-checked-item", formData)
		}
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

		items := []*WishlistItem{}
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

			items = append(items, &WishlistItem{Text: value[0], Index: index})
		}

		sort.Slice(items, func(i, j int) bool {
			return items[i].Index < items[j].Index
		})

		wishlist := data.Wishlists[id]

		data.Wishlists[id].Items = items

		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(wishlist))
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
		userIdCookie, error := c.Cookie("wishlist_uid")
		if error != nil {
			if error == http.ErrNoCookie {
				return c.Render(200, "wishlist", MapWishlistToWishlistViewFormData(wishlist))
			}
			return error
		}

		if userIdCookie.Value != data.Wishlists[id].UserId {
			return c.Render(200, "wishlist", MapWishlistToWishlistViewFormData(wishlist))
		}

		return c.Render(200, "index", MapWishlistToWishlistFormData(wishlist))
	})

	e.Static("/css", "css")
	e.Static("/js", "js")

	e.Logger.Fatal(e.Start(":1323"))
}
