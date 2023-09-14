package data

import (
	"database/sql"
	"fmt"
)

type ClientDBStore struct {
	db *sql.DB
}

func NewClientDetailsStore(db *sql.DB) ClientStore {
	return &ClientDBStore{db}
}

func (u *ClientDBStore) InsertClientDetails(details ClientDetail) (ClientDetail, error) {
	query := "INSERT INTO client_details (id,name,emailAddress) VALUES (?,?,?)"
	fmt.Println("==============> Details: ", details)
	_, err := u.db.Exec(query, details.ID, details.Name, details.EmailAddress)
	if err != nil {
		return ClientDetail{}, err
	}

	return details, nil
}

func (u *ClientDBStore) GetClientDetailsById(id string) (ClientDetail, error) {
	query := "SELECT * FROM client_details WHERE id=?"

	res := ClientDetail{}

	err := u.db.QueryRow(query, id).Scan(&res.ID, &res.Name, &res.EmailAddress)

	if err == sql.ErrNoRows {
		return ClientDetail{}, err
	}

	return res, nil

}

func (u *ClientDBStore) GetClientDetailsByEmail(emailAddress string) (ClientDetail, error) {
	query := "SELECT * FROM client_details WHERE emailAddress = ?"

	res := ClientDetail{}

	err := u.db.QueryRow(query, emailAddress).Scan(&res.ID, &res.Name, &res.EmailAddress)

	if err == sql.ErrNoRows {
		return ClientDetail{}, err
	}

	return res, nil
}
