package repository

import (
	"creeston/lists/internal/domain"
	"time"
)

type WishlistRepository interface {
	AddWishlist(wishlist *domain.Wishlist, creatorId string) *domain.Wishlist
	GetWishlistById(id int) *domain.Wishlist
	GetRecentWishlistsCreatedByUserCount(userId string, ipAddress string, period time.Duration) int
	UpdateWishlist(id int, wishlist *domain.Wishlist) *domain.Wishlist
	UpdateWishlistItem(id int, item domain.WishlistItem) domain.WishlistItem
}
