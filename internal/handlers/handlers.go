package handlers

import (
	"creeston/lists/internal/domain"
	"creeston/lists/internal/repository"
	"creeston/lists/internal/utils"
	"fmt"
	"net/url"
	"strconv"

	"github.com/labstack/echo/v4"

	"golang.org/x/text/message"

	_ "creeston/lists/internal/translations"
)

func GetPrinter(c echo.Context) *message.Printer {
	return c.Get("i18n").(*message.Printer)
}

func GetClientLanguage(c echo.Context) string {
	return c.Get("clientLanguage").(string)
}

func GetUserId(c echo.Context) string {
	return c.Get("userId").(string)
}

func SetupRoutes(e *echo.Echo, repo repository.WishlistRepository, baseUrl string, validationConfig ValidationConfig) {
	e.GET("/", func(c echo.Context) error {
		i18n := GetPrinter(c)
		language := GetClientLanguage(c)
		return c.Render(
			200,
			"index",
			WishlistFormViewParams{
				HasItems: false,
				HasId:    false,
				BaseUrl:  baseUrl,
				ValidationErrors: ValidationErrors{
					FieldErrors: map[string]string{},
					Errors:      map[string]string{},
				},
				Labels: getLabelsData(i18n, language),
			})
	})

	e.GET("/wishlist/:id", func(c echo.Context) error {
		i18n := GetPrinter(c)
		language := GetClientLanguage(c)
		userId := GetUserId(c)

		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		key := c.QueryParam("key")
		wishlist := repo.GetWishlistByID(id)

		if wishlist == nil {
			return c.Render(200, "not-found", NotFoundViewParams{
				Labels:  getLabelsData(i18n, language),
				BaseUrl: baseUrl,
			})
		}

		if wishlist.CreatorId == userId {
			formData := MapWishlistToWishlistFormData(wishlist)
			formData.Labels = getLabelsData(i18n, language)
			formData.BaseUrl = baseUrl
			return c.Render(200, "index", formData)
		}

		if wishlist.Key != key {
			return c.Render(200, "not-found", NotFoundViewParams{
				Labels:  getLabelsData(i18n, language),
				BaseUrl: baseUrl,
			})
		}

		viewData := MapWishlistToWishlistViewFormData(wishlist, userId)
		viewData.Labels = getLabelsData(i18n, language)
		viewData.BaseUrl = baseUrl
		return c.Render(200, "wishlist", viewData)
	})

	e.POST("/wishlist", func(c echo.Context) error {
		params, error := c.FormParams()
		if error != nil {
			return error
		}

		i18n := GetPrinter(c)
		language := GetClientLanguage(c)
		items := ParseWishlistFormDataToNewWishlistItems(params)
		validationErrors := validateWishlistFormData(items, i18n, validationConfig)
		if validationErrors.AnyErrors() {
			return c.Render(200, "wishlist-form", WishlistFormViewParams{
				HasItems:         false,
				HasId:            false,
				ValidationErrors: validationErrors,
				Labels:           getLabelsData(i18n, language),
			})
		}

		userId := GetUserId(c)
		wishlistKey := utils.GenerateUUID()
		wishlist := repo.AddWishlist(userId, wishlistKey, items)

		// Currently redirection implemented on the client side.
		// If we need to immediately redirect user, we should uncomment it.
		// c.Response().Header().Set("HX-Redirect", "/wishlist/"+strconv.Itoa(wishlist.Id))
		clientEvent := fmt.Sprintf("{\"wishlist-created\": %d}", wishlist.Id)
		c.Response().Header().Set("HX-Trigger", clientEvent)
		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(wishlist))
	})

	e.PUT("/wishlist/:id", func(c echo.Context) error {
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		i18n := GetPrinter(c)

		wishlist := repo.GetWishlistByID(id)
		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		userId := GetUserId(c)
		if userId != wishlist.CreatorId {
			return c.String(403, "Forbidden")
		}

		params, error := c.FormParams()
		if error != nil {
			return error
		}

		items := ParseWishlistFormDataToUpdatedWishlistItems(params)
		wishlist.UpdateWishlistItems(items)
		validationErrors := validateUpdateWishlistFormData(wishlist.Items, i18n, validationConfig)
		if validationErrors.AnyErrors() {
			formData := MapWishlistToWishlistFormData(wishlist)
			formData.ValidationErrors = validationErrors
			formData.Labels = getLabelsData(i18n, GetClientLanguage(c))
			return c.Render(200, "wishlist-form", formData)
		}

		updatedWishlist := repo.UpdateWishlist(id, wishlist)
		return c.Render(200, "wishlist-form", MapWishlistToWishlistFormData(updatedWishlist))
	})

	e.PUT("/wishlist/:id/:itemId", func(c echo.Context) error {
		i18n := GetPrinter(c)
		userId := GetUserId(c)

		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		// take hx-current-url from headers
		currentUrl := c.Request().Header.Get("HX-Current-URL")
		u, err := url.Parse(currentUrl)
		if err != nil {
			panic(err)
		}

		m, _ := url.ParseQuery(u.RawQuery)
		key := m.Get("key")

		itemId, error := strconv.Atoi(c.Param("itemId"))
		if error != nil {
			return error
		}

		wishlist := repo.GetWishlistByID(id)
		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		if (itemId < 0) || (itemId >= len(wishlist.Items)) {
			return c.String(404, "Item not found")
		}

		if wishlist.Key != key {
			return c.String(403, "Forbidden")
		}

		checkRequest := c.FormValue(("flag")) == "on"
		wishlistItem := wishlist.Items[itemId]
		viewParams := WishlistCheckableItemParams{
			Index: wishlistItem.Id,
			Text:  wishlistItem.Text,
			Id:    wishlist.Id,
		}

		if checkRequest && wishlistItem.Checked && wishlistItem.CheckedById == userId {
			return c.Render(200, "wishlist-checked-item", viewParams)
		}

		if checkRequest && wishlistItem.Checked && wishlistItem.CheckedById != userId {
			viewData := WishlistAlredyCheckedItemParams{
				Index:  wishlistItem.Id,
				Text:   wishlistItem.Text,
				Labels: getLabelsData(i18n, GetClientLanguage(c)),
			}
			return c.Render(200, "wishlist-already-checked-item-with-popup", viewData)
		}

		if !checkRequest && !wishlistItem.Checked {
			return c.Render(200, "wishlist-not-checked-item", viewParams)
		}

		if !checkRequest && wishlistItem.Checked && wishlistItem.CheckedById != userId {
			return c.Render(200, "wishlist-already-checked-item", viewParams)
		}

		if !checkRequest && wishlistItem.Checked && wishlistItem.CheckedById == userId {
			wishlistItem.Checked = false
			wishlistItem.CheckedById = ""
			repo.UpdateWishlistItem(id, wishlistItem)
			return c.Render(200, "wishlist-not-checked-item", viewParams)
		}

		wishlistItem.Checked = true
		wishlistItem.CheckedById = userId
		repo.UpdateWishlistItem(id, wishlistItem)
		return c.Render(200, "wishlist-checked-item", viewParams)
	})
}

func validateWishlistFormData(items []string, i18n *message.Printer, validationConfig ValidationConfig) ValidationErrors {
	var validationErrors = ValidationErrors{
		FieldErrors: map[string]string{},
		Errors:      map[string]string{},
	}
	for _, item := range items {
		if len(item) > validationConfig.MaxItemLength {
			validationErrors.FieldErrors[item] = i18n.Sprintf("Item text is too long. Maximum length is %d characters", validationConfig.MaxItemLength)
		}
	}

	if len(items) > validationConfig.MaxItemsCount {
		validationErrors.Errors["maxItemsCount"] = i18n.Sprintf("Maximum number of items is %d", validationConfig.MaxItemsCount)
	}

	if len(items) == 0 {
		validationErrors.Errors["noItems"] = i18n.Sprintf("No items provided")
	}

	return validationErrors
}

func validateUpdateWishlistFormData(items []*domain.WishlistItem, i18n *message.Printer, validationConfig ValidationConfig) ValidationErrors {
	var validationErrors = ValidationErrors{
		FieldErrors: map[string]string{},
		Errors:      map[string]string{},
	}
	for _, item := range items {
		if len(item.Text) > validationConfig.MaxItemLength {
			validationErrors.FieldErrors[item.Text] = i18n.Sprintf("Item text is too long. Maximum length is %d characters", validationConfig.MaxItemLength)
		}
	}

	if len(items) > validationConfig.MaxItemsCount {
		validationErrors.Errors["maxItemsCount"] = i18n.Sprintf("Maximum number of items is %d", validationConfig.MaxItemsCount)
	}

	if len(items) == 0 {
		validationErrors.Errors["noItems"] = i18n.Sprintf("No items provided")
	}

	return validationErrors
}
