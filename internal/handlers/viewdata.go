package handlers

import (
	"creeston/lists/internal/domain"
	"sort"
	"strconv"

	"golang.org/x/text/message"
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

func getLabelsData(i18n *message.Printer, selectedLanguage string) LabelsData {
	labels := getLanguageList(i18n)
	return LabelsData{
		SelectedLanguage: selectedLanguage,
		Languages:        labels,

		CopyToClipboardTooltipLabel:     i18n.Sprintf("Copy link to share"),
		WishlistItemPlaceholder:         i18n.Sprintf("Add a wish"),
		SaveButtonTitle:                 i18n.Sprintf("Save"),
		EditButtonTitle:                 i18n.Sprintf("Edit"),
		NotFoundTitle:                   i18n.Sprintf("Wishlist not found"),
		CreateNewWishlistButton:         i18n.Sprintf("Create a new wishlist"),
		ItemWasAlreadyCheckedPopupTitle: i18n.Sprintf("This item is already reserved"),
		ItemWasAlreadyCheckedPopupText:  i18n.Sprintf("This item is already taken by another user"),
		ItemWasAlreadyCheckedOkayButton: i18n.Sprintf("Got it"),
	}
}

type LabelsData struct {
	SelectedLanguage string
	Languages        []LanguageData

	CopyToClipboardTooltipLabel     string
	WishlistItemPlaceholder         string
	SaveButtonTitle                 string
	EditButtonTitle                 string
	NotFoundTitle                   string
	CreateNewWishlistButton         string
	ItemWasAlreadyCheckedPopupTitle string
	ItemWasAlreadyCheckedPopupText  string
	ItemWasAlreadyCheckedOkayButton string
}

type WishlistFormViewParams struct {
	Items            []WishlistFormItemParams
	HasItems         bool
	Id               int
	HasId            bool
	Key              string
	ValidationErrors ValidationErrors

	BaseUrl string
	Labels  LabelsData
}

type WishlistFormItemParams struct {
	Id             int
	HasId          bool
	Text           string
	AlreadyChecked bool
}

type WishlistViewParams struct {
	Items   []WishlistCheckableItemParams
	Id      int
	BaseUrl string
	Labels  LabelsData
}

type WishlistCheckableItemParams struct {
	Index                int
	Text                 string
	Id                   int
	Checked              bool
	CheckedByAnotherUser bool
}

type WishlistAlredyCheckedItemParams struct {
	Text   string
	Index  int
	Labels LabelsData
}

type NotFoundViewParams struct {
	BaseUrl string
	Labels  LabelsData
}

func MapWishlistToWishlistFormData(wishlist *domain.Wishlist) WishlistFormViewParams {
	items := []WishlistFormItemParams{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistFormItemParams{
			Id:             item.Id,
			HasId:          item.HasId,
			Text:           item.Text,
			AlreadyChecked: item.IsTaken(),
		})
	}

	return WishlistFormViewParams{
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

func MapWishlistToWishlistViewFormData(wishlist *domain.Wishlist, userId string) WishlistViewParams {
	items := []WishlistCheckableItemParams{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistCheckableItemParams{
			Index:                item.Id,
			Text:                 item.Text,
			Id:                   wishlist.Id,
			Checked:              item.IsTaken(),
			CheckedByAnotherUser: item.IsTaken() && item.TakenById != userId,
		})
	}

	return WishlistViewParams{
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

func ParseWishlistFormDataToUpdatedWishlistItems(data map[string][]string) []domain.UpdateWishlistItemCommand {
	items := []domain.UpdateWishlistItemCommand{}

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
				items = append(items, domain.UpdateWishlistItemCommand{
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
			items = append(items, domain.UpdateWishlistItemCommand{
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
