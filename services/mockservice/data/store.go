package data

type ClientStore interface {
	InsertClientDetails(details ClientDetail) (ClientDetail, error)
	GetClientDetailsById(id string) (ClientDetail, error)
	GetClientDetailsByEmail(emailAddress string) (ClientDetail, error)
}

type RequestStore interface {
	GetRequestsByClientId(clientId string) ([]Request, error)
	InsertMoneyRequest(request *Request) (*Request, error)
	GetMoneyRequestById(id string) (*Request, error)
	UpdateMoneyRequestStatusById(id string, status string) error
}
