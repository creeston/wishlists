package main

type Data struct {
	Wishlists Wishlists
}

func NewData() *Data {
	data := &Data{
		Wishlists: Wishlists{},
	}

	data.AddWishlist("default", []*WishlistItem{
		{Text: "Cake", Index: 0},
		{Text: "Candles", Index: 1},
		{Text: "Balloons", Index: 2},
		{Text: "Presents. A lot a lot a lof a very long list of presents please!", Index: 3},
	})

	return data
}

func (data *Data) UpdateWishlistWithItems(wishlistId int, items []*WishlistItem) {
	wishlist := data.GetWishlistByIdOrNull(wishlistId)
	if wishlist == nil {
		return
	}

	for _, item := range items {
		existingItem := wishlist.GetItemByIndex(item.Index)
		if existingItem == nil {
			wishlist.Items = append(wishlist.Items, item)
			continue
		}

		if existingItem.Checked {
			continue
		}

		existingItem.Text = item.Text
	}
}

func (data *Data) AddWishlist(userId string, items []*WishlistItem) *Wishlist {
	wishlistId := len(data.Wishlists)
	wishlist := NewWishlist(items, wishlistId, userId)
	data.Wishlists = append(data.Wishlists, wishlist)

	return wishlist
}

func (data *Data) GetWishlistByIdOrNull(id int) *Wishlist {
	if (id < 0) || (id >= len(data.Wishlists)) {
		return nil
	}

	return data.Wishlists[id]
}

func (wishlist *Wishlist) GetItemByIndex(index int) *WishlistItem {
	for _, item := range wishlist.Items {
		if item.Index == index {
			return item
		}
	}

	return nil
}
