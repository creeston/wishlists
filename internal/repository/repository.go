package repository

import (
	"creeston/lists/internal/domain"
)

type WishlistRepository interface {
	AddWishlist(userID string, wishlistKey string, items []*domain.WishlistItem) *domain.Wishlist
	GetWishlistByID(id int) *domain.Wishlist
	UpdateWishlistWithItems(id int, items []*domain.WishlistItem)
}

type InMemoryRepository struct {
	wishlists []*domain.Wishlist
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{
		wishlists: make([]*domain.Wishlist, 0),
	}

	repository.AddWishlist("default", "public", []*domain.WishlistItem{
		{Text: "Cake", Index: 0},
		{Text: "Candles", Index: 1},
		{Text: "Balloons", Index: 2},
		{Text: "Presents. A lot a lot a lof a very long list of presents please!", Index: 3},
	})

	return repository
}

func (r *InMemoryRepository) AddWishlist(userID string, wishlistKey string, items []*domain.WishlistItem) *domain.Wishlist {
	id := len(r.wishlists)
	wishlist := domain.NewWishlist(items, id, userID, wishlistKey)
	r.wishlists = append(r.wishlists, wishlist)
	return wishlist
}

func (r *InMemoryRepository) GetWishlistByID(id int) *domain.Wishlist {
	if id < 0 || id >= len(r.wishlists) {
		return nil
	}
	return r.wishlists[id]
}

func (r *InMemoryRepository) UpdateWishlistWithItems(id int, items []*domain.WishlistItem) {
	wishlist := r.GetWishlistByID(id)
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
