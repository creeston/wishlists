package handlers

import (
	"creeston/lists/internal/domain"
	"creeston/lists/internal/repository"
	"creeston/lists/internal/utils"
	"fmt"
	"net/url"
	"strconv"
	"unicode/utf8"

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
		userId := GetUserId(c)
		wishlistKey := utils.GenerateUUID()
		items := ParseWishlistFormDataToNewWishlistItems(params)
		wishlist := domain.NewWishlist(userId, wishlistKey, items)
		validationErrors := validateWishlistItems(wishlist.Items, i18n, validationConfig)
		if validationErrors.AnyErrors() {
			return c.Render(200, "wishlist-form", WishlistFormViewParams{
				HasItems:         false,
				HasId:            false,
				ValidationErrors: validationErrors,
				Labels:           getLabelsData(i18n, language),
			})
		}

		wishlist = repo.AddWishlist(wishlist)

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
		validationErrors := validateWishlistItems(wishlist.Items, i18n, validationConfig)
		if validationErrors.AnyErrors() {
			formProps := MapWishlistToWishlistFormData(wishlist)
			formProps.ValidationErrors = validationErrors
			formProps.Labels = getLabelsData(i18n, GetClientLanguage(c))
			return c.Render(200, "wishlist-form", formProps)
		}

		updatedWishlist := repo.UpdateWishlist(id, wishlist)
		formProps := MapWishlistToWishlistFormData(updatedWishlist)
		formProps.Labels = getLabelsData(i18n, GetClientLanguage(c))
		return c.Render(200, "wishlist-form", formProps)
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

		if wishlist.Key != key {
			return c.String(403, "Forbidden")
		}

		wishlistItem := wishlist.GetItemByIndex(itemId)

		if wishlistItem == nil {
			return c.String(404, "Item not found")
		}

		checkRequest := c.FormValue(("flag")) == "on"
		viewParams := WishlistCheckableItemParams{
			Index: wishlistItem.Id,
			Text:  wishlistItem.Text,
			Id:    wishlist.Id,
		}

		if checkRequest && wishlistItem.IsTaken() && wishlistItem.TakenById == userId {
			return c.Render(200, "wishlist-checked-item", viewParams)
		}

		if checkRequest && wishlistItem.IsTaken() && wishlistItem.TakenById != userId {
			viewData := WishlistAlredyCheckedItemParams{
				Index:  wishlistItem.Id,
				Text:   wishlistItem.Text,
				Labels: getLabelsData(i18n, GetClientLanguage(c)),
			}
			return c.Render(200, "wishlist-already-checked-item-with-popup", viewData)
		}

		if !checkRequest && !wishlistItem.IsTaken() {
			return c.Render(200, "wishlist-not-checked-item", viewParams)
		}

		if !checkRequest && wishlistItem.IsTaken() && wishlistItem.TakenById != userId {
			return c.Render(200, "wishlist-already-checked-item", viewParams)
		}

		if !checkRequest && wishlistItem.IsTaken() && wishlistItem.TakenById == userId {
			wishlistItem.TakenById = ""
			repo.UpdateWishlistItem(id, wishlistItem)
			return c.Render(200, "wishlist-not-checked-item", viewParams)
		}

		wishlistItem.TakenById = userId
		repo.UpdateWishlistItem(id, wishlistItem)
		return c.Render(200, "wishlist-checked-item", viewParams)
	})
}

func validateWishlistItems(items []*domain.WishlistItem, i18n *message.Printer, validationConfig ValidationConfig) ValidationErrors {
	var validationErrors = ValidationErrors{
		FieldErrors: map[string]string{},
		Errors:      map[string]string{},
	}
	for _, item := range items {
		textLength := utf8.RuneCountInString(item.Text)
		if textLength > validationConfig.MaxItemLength {
			validationErrors.FieldErrors[item.Text] = i18n.Sprintf("Your text is %d characters; max is %d", textLength, validationConfig.MaxItemLength)
		}
	}

	if len(items) > validationConfig.MaxItemsCount {
		validationErrors.Errors["maxItemsCount"] = i18n.Sprintf("Max items: %d. You added %d.", validationConfig.MaxItemsCount, len(items))
	}

	if len(items) == 0 {
		validationErrors.Errors["noItems"] = i18n.Sprintf("No items provided")
	}

	return validationErrors
}
