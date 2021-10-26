package postgres

import (
	"strconv"

	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/application/enums"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres/gorm"
)

// mapFacilityObjectToDomain maps the db facility to a domain model.
// It fetches the database to fetch items specific to the facility
func (d *OnboardingDb) mapFacilityObjectToDomain(facilityObject *gorm.Facility) *domain.Facility {
	if facilityObject == nil {
		return nil
	}

	active, err := strconv.ParseBool(facilityObject.Active)
	if err != nil {
		return nil
	}

	return &domain.Facility{
		ID:          facilityObject.FacilityID,
		Name:        facilityObject.Name,
		Code:        facilityObject.Code,
		Active:      active,
		County:      facilityObject.County,
		Description: facilityObject.Description,
	}
}

// mapMetricObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapMetricObjectToDomain(metricObject *gorm.Metric) *domain.Metric {
	if metricObject == nil {
		return nil
	}

	return &domain.Metric{
		MetricID:  metricObject.MetricID,
		Type:      metricObject.Type,
		Payload:   metricObject.Payload,
		Timestamp: metricObject.Timestamp,
		UID:       metricObject.UID,
	}
}

// mapProfileObjectToDomain maps the db metrics to a domain model.
// It searches the database to fetch items specific to the metrics
func (d *OnboardingDb) mapProfileObjectToDomain(profileObject *gorm.User) *domain.User {
	if profileObject == nil {
		return nil
	}

	user := createMapUser(profileObject)

	return user
}

// mapPINObjectToDomain maps the db pin data to a domain model.
// It searches the database to fetch items specific to the PIN
func (d *OnboardingDb) mapPINObjectToDomain(pinObj *gorm.PINData) *domain.UserPIN {
	if pinObj == nil {
		return nil
	}

	return &domain.UserPIN{
		UserID:    pinObj.UserID,
		HashedPIN: pinObj.HashedPIN,
		ValidFrom: pinObj.ValidFrom,
		ValidTo:   pinObj.ValidTo,
		Flavour:   pinObj.Flavour,
		IsValid:   pinObj.IsValid,
		Salt:      pinObj.Salt,
	}
}

// mapUserObjectToDomain maps the db user to a domain model.
// It searches the database to fetch items specific to the user
func (d *OnboardingDb) mapRegisterStaffObjectToDomain(userStaffObject *gorm.StaffUserProfile) *domain.StaffUserProfile {

	userObject := userStaffObject.User
	staffObject := userStaffObject.Staff

	facilities := []*domain.Facility{}

	for _, f := range staffObject.Facilities {
		active, err := strconv.ParseBool(f.Active)
		if err != nil {
			return nil
		}
		facility := &domain.Facility{
			ID:          f.FacilityID,
			Name:        f.Name,
			Code:        f.Code,
			Active:      active,
			County:      f.County,
			Description: f.Description,
		}
		facilities = append(facilities, facility)
	}

	user := createMapUser(userObject)

	addresses := []*domain.Addresses{}
	for _, a := range staffObject.Addresses {
		address := &domain.Addresses{
			Type:       a.Type,
			Text:       a.Text,
			Country:    a.Country,
			PostalCode: a.PostalCode,
			County:     a.County,
			Active:     a.Active,
		}
		addresses = append(addresses, address)
	}

	roles := []enums.RolesType{}
	for _, r := range staffObject.Roles {
		roles = append(roles, enums.RolesType(r))

	}

	staffProfile := &domain.StaffProfile{
		ID:                staffObject.StaffProfileID,
		UserID:            userObject.UserID,
		StaffNumber:       staffObject.StaffNumber,
		Facilities:        facilities,
		DefaultFacilityID: staffObject.DefaultFacilityID,
		Addresses:         addresses,
		Roles:             roles,
	}
	return &domain.StaffUserProfile{
		User:  user,
		Staff: staffProfile,
	}
}

// mapRegisterClientObjectToDomain maps the database client object to our custom domain type
func (d *OnboardingDb) mapRegisterClientObjectToDomain(clientObject *gorm.ClientUserProfile) *domain.ClientUserProfile {
	userObject := clientObject.User
	client := clientObject.Client

	user := createMapUser(userObject)

	clientProfile := &domain.ClientProfile{
		ID:         client.ID,
		UserID:     client.UserID,
		ClientType: client.ClientType,
	}

	return &domain.ClientUserProfile{
		User:   user,
		Client: clientProfile,
	}
}

// a helper method to create mapped user
func createMapUser(userObject *gorm.User) *domain.User {
	contacts := []*domain.Contact{}
	if len(userObject.Contacts) > 0 {
		for _, u := range userObject.Contacts {
			contact := &domain.Contact{
				ID:      u.ContactID,
				Type:    u.Type,
				Contact: u.Contact,
				Active:  u.Active,
				OptedIn: u.OptedIn,
			}
			contacts = append(contacts, contact)
		}
	}

	languages := []enumutils.Language{}
	for _, l := range userObject.Languages {
		languages = append(languages, enumutils.Language(l))

	}

	user := &domain.User{
		ID:                  userObject.UserID,
		Username:            userObject.Username,
		DisplayName:         userObject.DisplayName,
		FirstName:           userObject.FirstName,
		MiddleName:          userObject.MiddleName,
		LastName:            userObject.LastName,
		Gender:              userObject.Gender,
		UserType:            userObject.UserType,
		Contacts:            contacts,
		Languages:           languages,
		Active:              userObject.Active,
		LastSuccessfulLogin: userObject.LastSuccessfulLogin,
		LastFailedLogin:     userObject.LastFailedLogin,
		FailedLoginCount:    userObject.FailedLoginCount,
		NextAllowedLogin:    userObject.NextAllowedLogin,
		TermsAccepted:       userObject.TermsAccepted,
		AcceptedTermsID:     userObject.AcceptedTermsID,
		Flavour:             userObject.Flavour,
	}
	return user
}

// mapIdentifierObjectToDomain maps the identifier object to our domain defined type
func (d *OnboardingDb) mapIdentifierObjectToDomain(identifierObject *gorm.Identifier) *domain.Identifier {
	if identifierObject == nil {
		return nil
	}

	return &domain.Identifier{
		ID:                  identifierObject.ID,
		ClientID:            identifierObject.ClientID,
		IdentifierType:      identifierObject.IdentifierType,
		IdentifierUse:       identifierObject.IdentifierUse,
		IdentifierValue:     identifierObject.IdentifierValue,
		Description:         identifierObject.Description,
		ValidFrom:           identifierObject.ValidFrom,
		ValidTo:             identifierObject.ValidTo,
		Active:              identifierObject.Active,
		IsPrimaryIdentifier: identifierObject.IsPrimaryIdentifier,
	}
}

// mapClientObjectToDomain maps the client object to the domain defined type
func (d *OnboardingDb) mapClientObjectToDomain(client *gorm.ClientProfile) *domain.ClientProfile {
	if client == nil {
		return nil
	}

	return &domain.ClientProfile{
		ID:                      client.ID,
		UserID:                  client.UserID,
		TreatmentEnrollmentDate: client.TreatmentEnrollmentDate,
		ClientType:              client.ClientType,
		Active:                  client.Active,
		HealthRecordID:          client.HealthRecordID,
		// Identifiers:             client.Identifiers,
		FacilityID:           client.FacilityID,
		TreatmentBuddyUserID: client.TreatmentBuddy,
		CHVUserID:            client.CHVUserID,
		ClientCounselled:     client.ClientCounselled,
	}
}

//mapStaffObjectToDomain maps the staff object to the domain defined type
func (d *OnboardingDb) mapStaffObjectToDomain(staff *gorm.StaffProfile) *domain.StaffProfile {
	if staff == nil {
		return nil
	}

	addresses := []*domain.Addresses{}
	if len(staff.Addresses) > 0 {
		for _, a := range staff.Addresses {
			address := &domain.Addresses{
				ID:         *a.AddressesID,
				Type:       a.Type,
				Text:       a.Text,
				Country:    a.Country,
				PostalCode: a.PostalCode,
				County:     a.County,
				Active:     a.Active,
			}
			addresses = append(addresses, address)
		}
	}

	return &domain.StaffProfile{
		ID:                staff.StaffProfileID,
		UserID:            staff.UserID,
		StaffNumber:       staff.StaffNumber,
		DefaultFacilityID: staff.DefaultFacilityID,
		Addresses:         addresses,
	}
}
