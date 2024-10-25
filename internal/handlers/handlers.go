package handlers

import (
	"creeston/lists/internal/repository"
	"strconv"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, repo *repository.Data, baseUrl string) {
	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", WishlistFormData{HasItems: false, HasId: false})
	})

	e.POST("/wishlist", func(c echo.Context) error {
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		items := ParseWishlistFormDataToWishlistItems(params)

		if len(items) == 0 {
			return c.String(400, "No items provided")
		}
		userId := c.Get("userId").(string)
		wishlist := repo.AddWishlist(userId, items)

		c.Response().Header().Set("HX-Redirect", "/wishlist/"+strconv.Itoa(wishlist.Id))
		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(wishlist))
	})

	e.PUT("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		wishlist := repo.GetWishlistByIdOrNull(id)

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
		repo.UpdateWishlistWithItems(id, items)
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

		wishlist := repo.GetWishlistByIdOrNull(id)
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

		if (id < 0) || (id >= len(repo.Wishlists)) {
			return c.Render(200, "not-found", nil)
		}

		wishlist := repo.Wishlists[id]
		userId := c.Get("userId").(string)
		if userId != wishlist.CreatorId {
			return c.Render(200, "wishlist", MapWishlistToWishlistViewFormData(wishlist, userId))
		}

		return c.Render(200, "index", MapWishlistToWishlistFormData(wishlist))
	})
}
