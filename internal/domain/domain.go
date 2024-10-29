package domain

type UpdateWishlistItem struct {
	Id    int
	Text  string
	HasId bool
}

type WishlistItem struct {
	Text        string
	Checked     bool
	CheckedById string
	Id          int
	HasId       bool
}

type Wishlist struct {
	Items     []*WishlistItem
	Id        int
	CreatorId string
	Key       string
}

type Wishlists = []*Wishlist

func NewWishlist(items []*WishlistItem, id int, userId string, wishlistKey string) *Wishlist {
	return &Wishlist{
		Items:     items,
		Id:        id,
		CreatorId: userId,
		Key:       wishlistKey,
	}
}

func (wishlist *Wishlist) GetItemByIndex(index int) *WishlistItem {
	for _, item := range wishlist.Items {
		if item.Id == index {
			return item
		}
	}

	return nil
}

func (wishlist *Wishlist) UpdateWishlistItems(items []UpdateWishlistItem) *Wishlist {
	newItems := make([]*WishlistItem, 0)
	for _, item := range wishlist.Items {
		if item.Checked {
			// If item was already checked by some user, then it becomes immutable
			newItems = append(newItems, item)
			continue
		}

		for _, newItem := range items {
			if newItem.HasId && newItem.Id == item.Id {
				item.Text = newItem.Text
				newItems = append(newItems, item)
				break
			}
		}
	}

	// add new items
	for _, item := range items {
		if item.HasId {
			continue
		}

		newItems = append(newItems, &WishlistItem{
			Text:  item.Text,
			HasId: false,
		})
	}

	wishlist.Items = newItems
	return wishlist
}
