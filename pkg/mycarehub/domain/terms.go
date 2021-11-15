package domain

// TermsOfService contains the struct field to hold the required display data for the terms of service.
type TermsOfService struct {
	TermsID int     `json:"terms_id"`
	Text    *string `json:"text"`
}
