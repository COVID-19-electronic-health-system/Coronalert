package models

// Number represents an individual phone number
type Number struct {
	Number string
}

// PhoneNumbers represents all phone numbers received from request
type PhoneNumbers struct {
	PhoneNumbers []Number
}
