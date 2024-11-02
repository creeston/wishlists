package handlers

import (
	"creeston/lists/internal/domain"
	"creeston/lists/internal/repository"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
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
		wishlist := repo.GetWishlistById(id)

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
		ipAddress := c.RealIP()
		userId := GetUserId(c)

		period := time.Hour * 24
		recentWishlistsCount := repo.GetRecentWishlistsCreatedByUserCount(userId, ipAddress, period)
		if recentWishlistsCount >= validationConfig.MaxWishlistsPerDay {
			return c.String(429, "Too many wishlists created")
		}

		params, error := c.FormParams()
		if error != nil {
			return error
		}

		i18n := GetPrinter(c)
		language := GetClientLanguage(c)
		items := ParseWishlistFormDataToNewWishlistItems(params)
		wishlist := domain.NewWishlist(items).CreatedBy(userId)
		validationErrors := validateWishlistItems(wishlist.Items, i18n, validationConfig)
		if validationErrors.AnyErrors() {
			return c.Render(200, "wishlist-form", WishlistFormViewParams{
				HasItems:         false,
				HasId:            false,
				ValidationErrors: validationErrors,
				Labels:           getLabelsData(i18n, language),
			})
		}

		wishlist = repo.AddWishlist(wishlist, ipAddress)

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

		wishlist := repo.GetWishlistById(id)
		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		userId := GetUserId(c)
		if !wishlist.IsAllowedToEdit(userId) {
			return c.String(403, "Forbidden")
		}

		params, error := c.FormParams()
		if error != nil {
			return error
		}

		updateItemCommands := ParseWishlistFormDataToUpdatedWishlistItems(params)
		wishlist.UpdateItems(updateItemCommands)
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
		id, error := strconv.Atoi(c.Param("id"))
		if error != nil {
			return error
		}

		itemId, error := strconv.Atoi(c.Param("itemId"))
		if error != nil {
			return error
		}

		currentUrl := c.Request().Header.Get("HX-Current-URL")
		key, error := getKeyFromUrl(currentUrl)
		if error != nil {
			panic(error)
		}

		wishlist := repo.GetWishlistById(id)
		if wishlist == nil {
			return c.String(404, "Wishlist not found")
		}

		if !wishlist.IsAllowedToView(key) {
			return c.String(403, "Forbidden")
		}

		wishlistItem := wishlist.GetItemByIndex(itemId)

		if wishlistItem == nil {
			return c.String(404, "Item not found")
		}

		checkRequest := strings.ToLower(c.FormValue(("flag"))) == "on"
		viewParams := WishlistCheckableItemParams{
			Index: wishlistItem.Id,
			Text:  wishlistItem.Text,
			Id:    wishlist.Id,
		}

		i18n := GetPrinter(c)
		userId := GetUserId(c)

		if checkRequest {
			err := wishlistItem.Take(userId)
			if err == domain.ErrAlreadyTakenItem {
				return c.Render(200, "wishlist-checked-item", viewParams)
			} else if err == domain.ErrAlreadyTakenItemByAnotherUser {
				viewData := WishlistAlredyCheckedItemParams{
					Index:  wishlistItem.Id,
					Text:   wishlistItem.Text,
					Labels: getLabelsData(i18n, GetClientLanguage(c)),
				}
				return c.Render(200, "wishlist-already-checked-item-with-popup", viewData)
			}
		} else {
			err := wishlistItem.Untake(userId)
			if err == domain.ErrAlreadyUntakenItem {
				return c.Render(200, "wishlist-not-checked-item", viewParams)
			} else if err == domain.ErrAlreadyTakenItemByAnotherUser {
				return c.Render(200, "wishlist-already-checked-item", viewParams)
			}
		}

		repo.UpdateWishlistItem(id, *wishlistItem)
		if wishlistItem.IsTaken() {
			return c.Render(200, "wishlist-checked-item", viewParams)
		} else {
			return c.Render(200, "wishlist-not-checked-item", viewParams)
		}
	})
}

func getKeyFromUrl(currentUrl string) (string, error) {
	parsedUrl, err := url.Parse(currentUrl)
	if err != nil {
		return "", err
	}

	queryValues, _ := url.ParseQuery(parsedUrl.RawQuery)
	key := queryValues.Get("key")
	return key, nil
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
