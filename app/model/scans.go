// DON'T EDIT *** generated by scaneo *** DON'T EDIT //

package model

import "database/sql"

func ScanUser(r *sql.Row) (User, error) {
	var s User
	if err := r.Scan(
		&s.ID,
		&s.Name,
		&s.Email,
		&s.Salt,
		&s.Salted,
		&s.Created,
		&s.Updated,
	); err != nil {
		return User{}, err
	}
	return s, nil
}

func ScanUsers(rs *sql.Rows) ([]User, error) {
	structs := make([]User, 0, 16)
	var err error
	for rs.Next() {
		var s User
		if err = rs.Scan(
			&s.ID,
			&s.Name,
			&s.Email,
			&s.Salt,
			&s.Salted,
			&s.Created,
			&s.Updated,
		); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}

func ScanItems(rs *sql.Rows) ([]Item, error) {
	structs := make([]Item, 0, 16)
	var err error
	for rs.Next() {
		var s Item
		if err = rs.Scan(
			&s.ID,
			&s.UserId,
			&s.Title,
			&s.Body,
			&s.Image,
			&s.Point,
			&s.Created,
			&s.Updated,
		); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}

func ScanItem(r *sql.Row) (Item, error) {
	var s Item
	if err := r.Scan(
		&s.ID,
		&s.UserId,
		&s.Title,
		&s.Body,
		&s.Image,
		&s.Point,
		&s.Created,
		&s.Updated,
	); err != nil {
		return Item{}, err
	}
	return s, nil
}

func ScanItemJoinUser(rs *sql.Rows) ([]ItemsJoinUser, error) {
	structs := make([]ItemsJoinUser, 0, 16)
	var err error
	for rs.Next() {
		var s ItemsJoinUser
		if err = rs.Scan(
			&s.Item_ID,
			&s.UserId,
			&s.Title,
			&s.Body,
			&s.Image,
			&s.Point,
			&s.Item_Created,
			&s.Item_Updated,
			&s.User_ID,
			&s.Name,
			&s.Email,
			&s.Salt,
			&s.Salted,
			&s.User_Created,
			&s.User_Updated,
		); err != nil {
			return nil, err
		}
		structs = append(structs, s)
	}
	if err = rs.Err(); err != nil {
		return nil, err
	}
	return structs, nil
}