package main

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

func MapWishlistToWishlistFormData(wishlist *Wishlist) WishlistFormData {
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

func MapWishlistToWishlistViewFormData(wishlist *Wishlist, userId string) WishlistViewFormData {
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
