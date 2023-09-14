package data

import (
	"database/sql"
	"errors"
	"fmt"
)

const (
	PENDING  = "PENDING"
	ACCEPTED = "ACCEPTED"
	DECLINES = "DECLINES"
)

type RequestDBStore struct {
	db *sql.DB
}

func NewRequestStore(db *sql.DB) RequestStore {
	return &RequestDBStore{db}
}

//func (u *RequestDBStore) InsertRequest(requestDetail Request) (Request, error) {
//	query := `INSERT INTO request_details
//    (id,requestedDate,requesterId,requesterName,requesteeName,requestStatus,expirationDate,amount,requesteeId,invoiceNumber,message)
//VALUES (?,?,?)`
//
//	_, err := u.db.Exec(query, requestDetail.ID, requestDetail.RequestedDate,
//		requestDetail.RequesterId, requestDetail.RequesterName, requestDetail.RequesteeName, requestDetail.RequestStatus,
//		requestDetail.ExpirationDate, requestDetail.Amount, requestDetail.RequesteeId, requestDetail.InvoiceNumber, requestDetail.Message)
//	if err != nil {
//		return Request{}, err
//	}
//
//	return requestDetail, nil
//}

// This should be all request recieved by
func (u *RequestDBStore) GetRequestsByClientId(clientId string) ([]Request, error) {

	requests := make([]Request, 0)

	query := "SELECT * FROM request_details WHERE requesteeId=?"

	results, err := u.db.Query(query, clientId)
	if err != nil {
		return nil, err
	}

	for results.Next() {
		request := Request{}
		err := results.Scan(&request.ID, &request.RequestedDate, &request.RequesterId, &request.RequesterName,
			&request.RequesteeName, &request.RequestStatus, &request.ExpirationDate, &request.Amount, &request.RequesteeId,
			&request.InvoiceNumber, &request.Message)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (u *RequestDBStore) InsertMoneyRequest(request *Request) (*Request, error) {
	query := `INSERT INTO request_details (id,requestedDate,requesterId,requesterName,requesteeName,
                             requestStatus,expirationDate,amount,requesteeId,invoiceNumber,message)
VALUES (?,?,?,?,?,?,?,?,?,?,?)`

	_, err := u.db.Exec(query, request.ID, request.RequestedDate, request.RequesterId, request.RequesterName, request.RequesteeName, request.RequestStatus,
		request.ExpirationDate, request.Amount, request.RequesteeId, request.InvoiceNumber, request.Message)
	if err != nil {
		return &Request{}, err
	}

	return request, nil
}

func (u *RequestDBStore) GetMoneyRequestById(id string) (*Request, error) {
	query := "SELECT * FROM request_details WHERE id=?"

	request := Request{}

	err := u.db.QueryRow(query, id).Scan(&request.ID, &request.RequestedDate, &request.RequesterId, &request.RequesterName,
		&request.RequesteeName, &request.RequestStatus, &request.ExpirationDate, &request.Amount, &request.RequesteeId, &request.InvoiceNumber, &request.Message)

	if err == sql.ErrNoRows {
		return &Request{}, errors.New(fmt.Sprintf("no record for the request id %s", id))
	}

	return &request, nil

}

func (u *RequestDBStore) UpdateMoneyRequestStatusById(id string, status string) error {
	query := "UPDATE request_details SET requestStatus = ? WHERE id = ? AND requestStatus = ?"
	_, err := u.db.Exec(query, status, id, PENDING)
	if err != nil {
		return err
	}
	return nil
}
