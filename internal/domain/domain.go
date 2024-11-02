package domain

import (
	"creeston/lists/internal/utils"
	"errors"
	"time"
)

var ErrAlreadyTakenItemByAnotherUser = errors.New("domain: wishlist item is already taken by another user")
var ErrAlreadyTakenItem = errors.New("domain: wishlist item is already taken by this user")
var ErrAlreadyUntakenItem = errors.New("domain: wishlist item is already untaken")

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

func (item WishlistItem) IsTakenBy(userId string) bool {
	return item.TakenById == userId
}

func (item *WishlistItem) Take(userId string) error {
	if item.IsTaken() {
		if item.IsTakenBy(userId) {
			return ErrAlreadyTakenItem
		} else {
			return ErrAlreadyTakenItemByAnotherUser
		}
	}

	item.TakenById = userId
	item.TakenAt = time.Now()
	return nil
}

func (item *WishlistItem) Untake(userId string) error {
	if !item.IsTaken() {
		return ErrAlreadyUntakenItem
	}

	if !item.IsTakenBy(userId) {
		return ErrAlreadyTakenItemByAnotherUser
	}

	item.TakenById = ""
	return nil
}

type Wishlist struct {
	Id        int
	Key       string
	Items     []*WishlistItem
	CreatorId string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWishlist(items []string) *Wishlist {
	wishlistItems := make([]*WishlistItem, 0)
	for _, item := range items {
		wishlistItems = append(wishlistItems, &WishlistItem{
			Text:  item,
			HasId: false,
		})
	}

	return &Wishlist{
		Items: wishlistItems,
		Key:   utils.GenerateUUID(),
	}
}

func (wishlist *Wishlist) CreatedBy(userId string) *Wishlist {
	wishlist.CreatorId = userId
	return wishlist
}

func (wishlist *Wishlist) IsAllowedToEdit(userId string) bool {
	return wishlist.CreatorId == userId
}

func (wishlist *Wishlist) IsAllowedToView(key string) bool {
	return wishlist.Key == key
}

func (wishlist *Wishlist) GetItemByIndex(index int) *WishlistItem {
	for _, item := range wishlist.Items {
		if item.Id == index {
			return item
		}
	}

	return nil
}

func (wishlist *Wishlist) UpdateItems(items []UpdateWishlistItemCommand) *Wishlist {
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
