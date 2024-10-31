package repository

import (
	"creeston/lists/internal/domain"
)

type WishlistRepository interface {
	AddWishlist(wishlist *domain.Wishlist) *domain.Wishlist
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

	wishlist := domain.NewWishlist("default", "public", []string{"Cake",
		"Candles",
		"Balloons",
		"Presents. A lot a lot a lof a very long list of presents please!",
	})

	repository.AddWishlist(wishlist)
	return repository
}

func (r *InMemoryRepository) AddWishlist(wishlist *domain.Wishlist) *domain.Wishlist {
	id := len(r.wishlists)
	items := make([]*domain.WishlistItem, 0)

	for i, item := range wishlist.Items {
		items = append(items, &domain.WishlistItem{
			Text:  item.Text,
			Id:    i,
			HasId: true,
		})
	}

	newWishlist := &domain.Wishlist{
		Items:     items,
		Id:        id,
		CreatorId: wishlist.CreatorId,
		Key:       wishlist.Key,
	}

	r.wishlists = append(r.wishlists, newWishlist)
	return r.GetWishlistByID(newWishlist.Id)
}

func (r *InMemoryRepository) GetWishlistByID(id int) *domain.Wishlist {
	if id < 0 || id >= len(r.wishlists) {
		return nil
	}

	wishlist := r.wishlists[id]
	items := make([]*domain.WishlistItem, 0)
	for _, item := range wishlist.Items {
		items = append(items, &domain.WishlistItem{
			Text:      item.Text,
			Id:        item.Id,
			HasId:     item.HasId,
			TakenById: item.TakenById,
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
	wishlistItem.TakenById = item.TakenById

	return &domain.WishlistItem{
		Text:      wishlistItem.Text,
		Id:        wishlistItem.Id,
		TakenById: wishlistItem.TakenById,
	}
}
