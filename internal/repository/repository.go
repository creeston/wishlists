package repository

import (
	"creeston/lists/internal/domain"
)

type WishlistRepository interface {
	AddWishlist(userID string, wishlistKey string, items []string) *domain.Wishlist
	GetWishlistByID(id int) *domain.Wishlist
	UpdateWishlist(id int, wishlist *domain.Wishlist) *domain.Wishlist
	UpdateWishlistItem(id int, item *domain.WishlistItem) *domain.WishlistItem
}

type InMemoryRepository struct {
	wishlists []*domain.Wishlist
}

func NewInMemoryRepository() *InMemoryRepository {
	repository := &InMemoryRepository{
		wishlists: make([]*domain.Wishlist, 0),
	}

	repository.AddWishlist("default", "public", []string{
		"Cake",
		"Candles",
		"Balloons",
		"Presents. A lot a lot a lof a very long list of presents please!",
	})

	return repository
}

func (r *InMemoryRepository) AddWishlist(userID string, wishlistKey string, items []string) *domain.Wishlist {
	id := len(r.wishlists)
	wishlistItems := make([]*domain.WishlistItem, 0)
	for i, item := range items {
		wishlistItems = append(wishlistItems, &domain.WishlistItem{
			Text:  item,
			Id:    i,
			HasId: true,
		})
	}
	wishlist := domain.NewWishlist(wishlistItems, id, userID, wishlistKey)
	r.wishlists = append(r.wishlists, wishlist)
	return wishlist
}

func (r *InMemoryRepository) GetWishlistByID(id int) *domain.Wishlist {
	if id < 0 || id >= len(r.wishlists) {
		return nil
	}

	wishlist := r.wishlists[id]
	items := make([]*domain.WishlistItem, 0)
	for _, item := range wishlist.Items {
		items = append(items, &domain.WishlistItem{
			Text:        item.Text,
			Id:          item.Id,
			HasId:       item.HasId,
			Checked:     item.Checked,
			CheckedById: item.CheckedById,
		})
	}
	return &domain.Wishlist{
		Items:     items,
		Id:        wishlist.Id,
		CreatorId: wishlist.CreatorId,
		Key:       wishlist.Key,
	}
}

func (r *InMemoryRepository) UpdateWishlist(id int, wishlist *domain.Wishlist) *domain.Wishlist {
	if id < 0 || id >= len(r.wishlists) {
		return nil
	}

	maxWishlistItemId := 0
	for _, item := range wishlist.Items {
		if item.Id > maxWishlistItemId {
			maxWishlistItemId = item.Id
		}
	}

	for _, item := range wishlist.Items {
		if !item.HasId {
			item.Id = maxWishlistItemId + 1
			item.HasId = true
			maxWishlistItemId++
		}
	}

	r.wishlists[id] = &domain.Wishlist{
		Items:     wishlist.Items,
		Id:        wishlist.Id,
		CreatorId: wishlist.CreatorId,
		Key:       wishlist.Key,
	}
	return r.GetWishlistByID(id)
}

func (r *InMemoryRepository) UpdateWishlistItem(id int, item *domain.WishlistItem) *domain.WishlistItem {
	wishlist := r.wishlists[id]
	if wishlist == nil {
		return nil
	}

	wishlistItem := wishlist.GetItemByIndex(item.Id)
	wishlistItem.Checked = item.Checked
	wishlistItem.CheckedById = item.CheckedById

	return &domain.WishlistItem{
		Text:        wishlistItem.Text,
		Id:          wishlistItem.Id,
		Checked:     wishlistItem.Checked,
		CheckedById: wishlistItem.CheckedById,
	}
}
