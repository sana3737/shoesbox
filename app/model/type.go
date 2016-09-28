//go:generate scaneo $GOFILE

package model

import "time"

// User returns model object for user.
type User struct {
	ID      int64      `json:"id"`
	Name    string     `json:"name"`
	Email   string     `json:"email"`
	Salt    string     `json:"salt"`
	Salted  string     `json:"salted"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

// Article returns model object for article.
type Item struct {
	ID      int64      `json:"id"`
	UserId  int64      `json:"userId"`
	Title   string     `json:"title"`
	Body    string     `json:"body"`
	Image   string     `json:"image"`
	Point   int64      `json:"point"`
	Created *time.Time `json:"created"`
	Updated *time.Time `json:"updated"`
}

type ItemsJoinUser struct {
	Item_ID      int64      `json:"id"`
	UserId       int64      `json:"userId"`
	Title        string     `json:"title"`
	Body         string     `json:"body"`
	Image        string     `json:"image"`
	Point        int64      `json:"point"`
	Item_Created *time.Time `json:"created"`
	Item_Updated *time.Time `json:"updated"`
	User_ID      int64      `json:"item_id"`
	Name         string     `json:"item_name"`
	Email        string     `json:"email"`
	Salt         string     `json:"salt"`
	Salted       string     `json:"salted"`
	User_Created *time.Time `json:"created"`
	User_Updated *time.Time `json:"updated"`
}
