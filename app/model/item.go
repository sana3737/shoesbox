package model

import (
	"database/sql"
)

func ItemsAll(db *sql.DB) ([]ItemsJoinUser, error) {
	rows, err := db.Query(`select * from items JOIN users ON items.user_id = users.user_id order by items.created DESC`)
	if err != nil {
		return nil, err
	}
	return ScanItemJoinUser(rows)
}

// ArticlesAll returns all articles.
func MyItemsAll(db *sql.DB, uid int64) ([]Item, error) {
	rows, err := db.Query(`select * from items where user_id = ?`, uid)
	if err != nil {
		return nil, err
	}
	return ScanItems(rows)
}

func NewItemsAll(db *sql.DB) ([]ItemsJoinUser, error) {
	rows, err := db.Query(`select * from items JOIN users ON items.user_id = users.user_id order by items.created DESC limit 3`)
	if err != nil {
		return nil, err
	}
	return ScanItemJoinUser(rows)
}

func RankItemsAll(db *sql.DB) ([]ItemsJoinUser, error) {
	rows, err := db.Query(`select * from items JOIN users ON items.user_id = users.user_id order by items.point DESC limit 3`)
	if err != nil {
		return nil, err
	}
	return ScanItemJoinUser(rows)
}

func ItemOne(db *sql.DB, id int64) (Item, error) {
	return ScanItem(db.QueryRow(`select * from items where item_id = ?`, id))
}

// Update updates article by given article.
func (t *Item) Update(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update items
		set title=?,body = ?,image = ?
		where item_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Title, t.Body, t.Image, t.ID)
}

func (t *Item) PointUpdate(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	update items
		set point=?
		where item_id = ?
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.Point, t.ID)
}

func (t *Item) Delete(tx *sql.Tx) (sql.Result, error) {
	stmt, err := tx.Prepare(`delete from items where item_id = ?`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(t.ID)
}

// Insert inserts new article.
func (t *Item) Insert(tx *sql.Tx, uid int64) (sql.Result, error) {
	stmt, err := tx.Prepare(`
	insert into items(user_id,title,body,image)
	values(?,?,?, ?)
	`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(uid, t.Title, t.Body, t.Image)
}
