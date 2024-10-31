package repository

import (
	"creeston/lists/internal/domain"
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteRepository struct {
	db *sql.DB
}

func NewSqliteRepository(filepath string) *SqliteRepository {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		panic(err)
	}

	repository := &SqliteRepository{
		db: db,
	}

	repository.Init()
	return repository
}

func NewInMemorySqliteRepository() *SqliteRepository {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	repository := &SqliteRepository{
		db: db,
	}

	repository.Init()
	return repository
}

func (r *SqliteRepository) Init() {
	_, err := r.db.Exec("CREATE TABLE IF NOT EXISTS wishlists (id INTEGER PRIMARY KEY, creator_id TEXT, key TEXT, created_at TEXT, updated_at TEXT)")
	if err != nil {
		panic(err)
	}

	_, err = r.db.Exec("CREATE TABLE IF NOT EXISTS wishlist_items (id INTEGER PRIMARY KEY, wishlist_id INTEGER, item_id INTEGER, text TEXT, taken_by TEXT, taken_at TEXT)")
	if err != nil {
		panic(err)
	}
}

func (r *SqliteRepository) AddWishlist(wishlist *domain.Wishlist) *domain.Wishlist {
	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}

	currentTime := time.Now()
	wishlist.CreatedAt = currentTime
	wishlist.UpdatedAt = currentTime
	result, err := tx.Exec(
		"INSERT INTO wishlists (creator_id, key, created_at, updated_at) VALUES (?, ?, ?, ?)",
		wishlist.CreatorId,
		wishlist.Key,
		wishlist.CreatedAt.Format(time.RFC3339),
		wishlist.UpdatedAt.Format(time.RFC3339),
	)

	if err != nil {
		panic(err)
	}

	wishlistId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	wishlist.Id = int(wishlistId)

	for index, item := range wishlist.Items {
		_, err = tx.Exec("INSERT INTO wishlist_items (wishlist_id, text, item_id, taken_by, taken_at) VALUES (?, ?, ?, ?, ?)", wishlistId, item.Text, index, "", "")
		if err != nil {
			panic(err)
		}

		item.Id = index
		item.HasId = true
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return wishlist
}

func (r *SqliteRepository) GetWishlistByID(id int) *domain.Wishlist {
	wishlist := &domain.Wishlist{}
	createdAt := ""
	updatedAt := ""
	err := r.db.QueryRow("SELECT id, creator_id, key, created_at, updated_at FROM wishlists WHERE id = ?", id).Scan(&wishlist.Id, &wishlist.CreatorId, &wishlist.Key, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}

		panic(err)
	}

	wishlist.CreatedAt, err = time.Parse(time.RFC3339, createdAt)
	if err != nil {
		panic(err)
	}

	wishlist.UpdatedAt, err = time.Parse(time.RFC3339, updatedAt)
	if err != nil {
		panic(err)
	}

	rows, err := r.db.Query("SELECT item_id, text, taken_by, taken_at FROM wishlist_items WHERE wishlist_id = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		item := &domain.WishlistItem{}
		takenAt := ""
		err := rows.Scan(&item.Id, &item.Text, &item.TakenById, &takenAt)
		if err != nil {
			panic(err)
		}

		if takenAt != "" {
			item.TakenAt, err = time.Parse(time.RFC3339, takenAt)
			if err != nil {
				panic(err)
			}
		}

		item.HasId = true
		wishlist.Items = append(wishlist.Items, item)
	}

	return wishlist
}

func (r *SqliteRepository) UpdateWishlist(id int, wishlist *domain.Wishlist) *domain.Wishlist {
	existingWishlist := r.GetWishlistByID(id)
	if existingWishlist == nil {
		return nil
	}

	maxWishlistItemId := 0
	for _, item := range existingWishlist.Items {
		if item.Id > maxWishlistItemId {
			maxWishlistItemId = item.Id
		}
	}

	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("DELETE FROM wishlist_items WHERE wishlist_id = ?", id)
	if err != nil {
		panic(err)
	}

	for _, item := range wishlist.Items {
		if !item.HasId {
			item.Id = maxWishlistItemId + 1
			item.HasId = true
			maxWishlistItemId++
		}

		takenAtValue := ""
		if item.IsTaken() {
			takenAtValue = item.TakenAt.Format(time.RFC3339)
		}

		_, err = tx.Exec("INSERT INTO wishlist_items (wishlist_id, item_id, text, taken_by, taken_at) VALUES (?, ?, ?, ?, ?)", id, item.Id, item.Text, item.TakenById, takenAtValue)
		if err != nil {
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}

	return r.GetWishlistByID(id)
}

func (r *SqliteRepository) UpdateWishlistItem(id int, item *domain.WishlistItem) *domain.WishlistItem {
	takenAtValue := ""
	if item.IsTaken() {
		item.TakenAt = time.Now()
		takenAtValue = item.TakenAt.Format(time.RFC3339)
	}

	internalId := r.getWishlistItemInternalId(id, item.Id)
	_, err := r.db.Exec("UPDATE wishlist_items SET taken_by = ?, taken_at = ? WHERE id = ?", item.TakenById, takenAtValue, internalId)
	if err != nil {
		panic(err)
	}

	return item
}

func (r *SqliteRepository) getWishlistItemInternalId(wishlistId int, itemId int) int {
	internalId := 0
	err := r.db.QueryRow("SELECT id FROM wishlist_items WHERE wishlist_id = ? AND item_id = ?", wishlistId, itemId).Scan(&internalId)
	if err != nil {
		panic(err)
	}

	return internalId
}
