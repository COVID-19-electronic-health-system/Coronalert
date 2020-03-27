package models

// Number represents an individual phone number
type Number struct {
	Number string `json:"number"`
}

// PhoneNumbers represents all phone numbers received from request
type PhoneNumbers struct {
	PhoneNumbers []Number `json:"phoneNumbers"`
}
