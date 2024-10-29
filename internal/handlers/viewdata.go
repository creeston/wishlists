package handlers

import (
	"creeston/lists/internal/domain"
	"strconv"
)

type ValidationConfig struct {
	MaxItemsCount int
	MaxItemLength int
}

type LanguageData struct {
	Language string
	Code     string
}

type ValidationErrors struct {
	FieldErrors map[string]string
	Errors      map[string]string
}

func (v *ValidationErrors) AnyErrors() bool {
	return len(v.Errors) > 0 || len(v.FieldErrors) > 0
}

type WishlistFormData struct {
	Items                       []WishlistFormItem
	HasItems                    bool
	HasId                       bool
	ValidationErrors            ValidationErrors
	Id                          int
	Key                         string
	CopyToClipboardTooltipLabel string
	WishlistItemPlaceholder     string
	SaveButtonTitle             string
	EditButtonTitle             string
	Languages                   []LanguageData
	SelectedLanguage            string
	BaseUrl                     string
}

type WishlistFormItem struct {
	Id             int
	HasId          bool
	Text           string
	AlreadyChecked bool
}

type WishlistViewFormData struct {
	Items            []WishlistCheckedItemData
	Id               int
	SaveButtonTitle  string
	EditButtonTitle  string
	Languages        []LanguageData
	SelectedLanguage string
	BaseUrl          string
}

type WishlistCheckedItemData struct {
	Index                int
	Text                 string
	Id                   int
	Checked              bool
	CheckedByAnotherUser bool
}

type WishlistAlredyCheckedItemData struct {
	Text                            string
	Index                           int
	ItemWasAlreadyCheckedPopupTitle string
	ItemWasAlreadyCheckedPopupText  string
	ItemWasAlreadyCheckedOkayButton string
}

type NotFoundData struct {
	NotFoundTitle           string
	CreateNewWishlistButton string
	Languages               []LanguageData
	SelectedLanguage        string
	BaseUrl                 string
}

func MapWishlistToWishlistFormData(wishlist *domain.Wishlist) WishlistFormData {
	items := []WishlistFormItem{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistFormItem{
			Id:             item.Id,
			HasId:          item.HasId,
			Text:           item.Text,
			AlreadyChecked: item.Checked,
		})
	}

	return WishlistFormData{
		Items:    items,
		HasItems: true,
		HasId:    true,
		Id:       wishlist.Id,
		Key:      wishlist.Key,
		ValidationErrors: ValidationErrors{
			FieldErrors: map[string]string{},
			Errors:      map[string]string{},
		},
	}
}

func MapWishlistToWishlistViewFormData(wishlist *domain.Wishlist, userId string) WishlistViewFormData {
	items := []WishlistCheckedItemData{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistCheckedItemData{
			Index:                item.Id,
			Text:                 item.Text,
			Id:                   wishlist.Id,
			Checked:              item.Checked,
			CheckedByAnotherUser: item.CheckedById != "" && item.CheckedById != userId,
		})
	}

	return WishlistViewFormData{
		Items: items,
		Id:    wishlist.Id,
	}
}

func ParseWishlistFormDataToNewWishlistItems(data map[string][]string) []string {
	items := []string{}
	formValues := data["item"]
	for _, value := range formValues {
		if len(value) == 0 {
			continue
		}

		if value == "" {
			continue
		}

		items = append(items, value)
	}

	return items
}

func ParseWishlistFormDataToUpdatedWishlistItems(data map[string][]string) []domain.UpdateWishlistItem {
	items := []domain.UpdateWishlistItem{}

	for key, value := range data {
		if key == "item" {
			values := []string{}
			for _, v := range value {
				if v == "" {
					continue
				}
				values = append(values, v)
			}
			for _, v := range values {
				items = append(items, domain.UpdateWishlistItem{
					Text:  v,
					HasId: false,
				})
			}
		} else if key[:4] == "item" {
			idValue := key[5:]
			id, err := strconv.Atoi(idValue)
			if err != nil {
				continue
			}
			value := value[0]
			println(id)
			println(value)
			items = append(items, domain.UpdateWishlistItem{
				Id:    id,
				Text:  value,
				HasId: true,
			})
		} else {
			continue
		}

	}

	return items
}
