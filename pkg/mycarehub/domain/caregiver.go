package domain

import "github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"

// Caregiver maps the schema for the table that stores the caregiver
type Caregiver struct {
	ID            string
	FirstName     string
	LastName      string
	PhoneNumber   string
	CaregiverType enums.CaregiverType
	Active        bool
}
