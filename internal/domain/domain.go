package domain

import "time"

type WishlistItem struct {
	Id        int
	HasId     bool
	Text      string
	TakenById string
	TakenAt   time.Time
}

func (item WishlistItem) IsTaken() bool {
	return item.TakenById != ""
}

type Wishlist struct {
	Id        int
	Key       string
	Items     []*WishlistItem
	CreatorId string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWishlist(userId string, wishlistKey string, items []string) *Wishlist {
	wishlistItems := make([]*WishlistItem, 0)
	for _, item := range items {
		wishlistItems = append(wishlistItems, &WishlistItem{
			Text:  item,
			HasId: false,
		})
	}

	return &Wishlist{
		Items:     wishlistItems,
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

func (wishlist *Wishlist) UpdateWishlistItems(items []UpdateWishlistItemCommand) *Wishlist {
	newItems := make([]*WishlistItem, 0)
	for _, item := range wishlist.Items {
		if item.IsTaken() {
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

type UpdateWishlistItemCommand struct {
	Id    int
	Text  string
	HasId bool
}
