package data

type Request struct {
	ID             string `json:"id"`
	RequestedDate  string `json:"requestedDate"`
	RequesterId    string `json:"requesterId"`
	RequesterName  string `json:"requesterName"`
	RequesteeName  string `json:"requesteeName"`
	RequestStatus  string `json:"requestStatus"`
	ExpirationDate string `json:"expirationDate"`
	Amount         string `json:"amount"`
	RequesteeId    string `json:"requesteeId"`
	InvoiceNumber  string `json:"invoiceNumber"`
	Message        string `json:"message"`
}
