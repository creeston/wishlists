package domain

type WishlistItem struct {
	Text        string
	Checked     bool
	CheckedById string
	Index       int
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
		if item.Index == index {
			return item
		}
	}

	return nil
}
