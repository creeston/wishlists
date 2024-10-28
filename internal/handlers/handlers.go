package handlers

import (
	"creeston/lists/internal/repository"
	"sort"
	"strconv"

	"github.com/labstack/echo/v4"

	"golang.org/x/text/message"

	_ "creeston/lists/internal/translations"
)

func getLanguageList(i18n *message.Printer) []LanguageData {
	languages := []LanguageData{
		{
			Language: i18n.Sprintf("English"),
			Code:     "en-GB",
		},
		{
			Language: i18n.Sprintf("Russian"),
			Code:     "ru-RU",
		},
		{
			Language: i18n.Sprintf("Polish"),
			Code:     "pl-PL",
		},
		{
			Language: i18n.Sprintf("Belarusian"),
			Code:     "be-BY",
		},
	}

	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Language < languages[j].Language
	})

	return languages
}

func SetupRoutes(e *echo.Echo, repo *repository.Data, baseUrl string) {
	e.GET("/", func(c echo.Context) error {
		i18n := c.Get("i18n").(*message.Printer)
		language := c.Get("clientLanguage").(string)
		return c.Render(
			200,
			"index",
			WishlistFormData{
				HasItems:                    false,
				HasId:                       false,
				CopyToClipboardTooltipLabel: i18n.Sprintf("Copy to clipboard"),
				WishlistItemPlaceholder:     i18n.Sprintf("Start typing..."),
				SaveButtonTitle:             i18n.Sprintf("Save"),
				EditButtonTitle:             i18n.Sprintf("Edit"),
				Languages:                   getLanguageList(i18n),
				SelectedLanguage:            language,
				BaseUrl:                     baseUrl,
			})
	})

	e.GET("/wishlist/:id", func(c echo.Context) error {
		i18n := c.Get("i18n").(*message.Printer)
		clientLanguage := c.Get("clientLanguage").(string)

		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		if (id < 0) || (id >= len(repo.Wishlists)) {
			return c.Render(200, "not-found", NotFoundData{
				NotFoundTitle:           i18n.Sprintf("Wishlist not found"),
				CreateNewWishlistButton: i18n.Sprintf("Create new wishlist"),
				Languages:               getLanguageList(i18n),
				SelectedLanguage:        clientLanguage,
				BaseUrl:                 baseUrl,
			})
		}

		wishlist := repo.Wishlists[id]
		userId := c.Get("userId").(string)
		if userId != wishlist.CreatorId {
			viewData := MapWishlistToWishlistViewFormData(wishlist, userId)
			viewData.EditButtonTitle = i18n.Sprintf("Edit")
			viewData.SaveButtonTitle = i18n.Sprintf("Save")
			viewData.Languages = getLanguageList(i18n)
			viewData.SelectedLanguage = clientLanguage
			viewData.BaseUrl = baseUrl
			return c.Render(200, "wishlist", viewData)
		}

		formData := MapWishlistToWishlistFormData(wishlist)
		formData.CopyToClipboardTooltipLabel = i18n.Sprintf("Copy to clipboard")
		formData.WishlistItemPlaceholder = i18n.Sprintf("Start typing...")
		formData.SaveButtonTitle = i18n.Sprintf("Save")
		formData.EditButtonTitle = i18n.Sprintf("Edit")
		formData.Languages = getLanguageList(i18n)
		formData.SelectedLanguage = clientLanguage
		formData.BaseUrl = baseUrl
		return c.Render(200, "index", formData)
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
			i18n := c.Get("i18n").(*message.Printer)
			viewData := WishlistAlredyCheckedItemData{
				Index:                           wishlistItem.Index,
				Text:                            wishlistItem.Text,
				ItemWasAlreadyCheckedPopupTitle: i18n.Sprintf("Item was already taken"),
				ItemWasAlreadyCheckedPopupText:  i18n.Sprintf("This item was already taken by another user"),
				ItemWasAlreadyCheckedOkayButton: i18n.Sprintf("Okay"),
			}
			return c.Render(200, "wishlist-already-checked-item-with-popup", viewData)
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

}
