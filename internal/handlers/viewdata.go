package handlers

import (
	"creeston/lists/internal/domain"
	"sort"
	"strconv"
)

type WishlistFormData struct {
	Items    []WishlistFormItem
	HasItems bool
	HasId    bool
	Id       int
}

type WishlistFormItem struct {
	Index          int
	Text           string
	AlreadyChecked bool
}

type WishlistViewFormData struct {
	Items []WishlistCheckedItemData
	Id    int
}

type WishlistCheckedItemData struct {
	Index                int
	Text                 string
	Id                   int
	Checked              bool
	CheckedByAnotherUser bool
}

func MapWishlistToWishlistFormData(wishlist *domain.Wishlist) WishlistFormData {
	items := []WishlistFormItem{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistFormItem{
			Index:          item.Index,
			Text:           item.Text,
			AlreadyChecked: item.Checked,
		})
	}

	return WishlistFormData{
		Items:    items,
		HasItems: true,
		HasId:    true,
		Id:       wishlist.Id,
	}
}

func MapWishlistToWishlistViewFormData(wishlist *domain.Wishlist, userId string) WishlistViewFormData {
	items := []WishlistCheckedItemData{}
	for _, item := range wishlist.Items {
		items = append(items, WishlistCheckedItemData{
			Index:                item.Index,
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

func ParseWishlistFormDataToWishlistItems(data map[string][]string) []*domain.WishlistItem {
	items := []*domain.WishlistItem{}
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

		items = append(items, &domain.WishlistItem{Text: value[0], Index: index})
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Index < items[j].Index
	})

	return items
}
