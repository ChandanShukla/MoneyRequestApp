package model

import (
	"errors"
	"strings"
)

type MoneyRequest struct {
	ExpirationDate string `json:"expirationDate"`
	Amount         string `json:"amount"`
	RequesteeId    string `json:"requesteeId"`
	InvoiceNumber  string `json:"invoiceNumber"`
	Message        string `json:"message"`
}

type UpdateMoneyRequest struct {
	MessageToRequester string `json:"messageToRequester"`
	Action             Action `json:"action"`
}

type Action string

func (a Action) ToString() (string, error) {
	action := strings.ToLower(string(a))

	switch action {
	case "accept":
		return "ACCEPT", nil
	case "decline":
		return "DECLINE", nil
	default:
		return "UNKNOWN", errors.New("invalid value of action")

	}
}
