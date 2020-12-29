package domain

import (
	"net/url"

	"gitlab.slade360emr.com/go/base"
)

// UserProfileInput is used to create or update a user's profile.
type UserProfileInput struct {
	PhotoUploadID string       `json:"photoUploadID"`
	DateOfBirth   *base.Date   `json:"dateOfBirth,omitempty"`
	Gender        *base.Gender `json:"gender,omitempty"`
	FirstName     *string      `json:"lastName"`
	LastName      *string      `json:"firstName"`
}

// PostVisitSurveyInput is used to send the results of post-visit surveys to the
// server.
type PostVisitSurveyInput struct {
	LikelyToRecommend int    `json:"likelyToRecommend" firestore:"likelyToRecommend"`
	Criticism         string `json:"criticism" firestore:"criticism"`
	Suggestions       string `json:"suggestions" firestore:"suggestions"`
}

// PayablesAccount stores a supplier's payables account info
type PayablesAccount struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	IsActive    bool   `json:"is_active"`
	Number      string `json:"number"`
	Tag         string `json:"tag"`
	Description string `json:"description"`
}

// Supplier used to create a supplier request payload
type Supplier struct {
	ProfileID              *string                `json:"profileID" firestore:"profileID"`
	SupplierID             string                 `json:"id" firestore:"supplierid"`
	PayablesAccount        *PayablesAccount       `json:"payablesAccount"`
	SupplierKYC            map[string]interface{} `json:"supplierKYC"`
	Active                 bool                   `json:"active" firestore:"active"`
	AccountType            AccountType            `json:"accountType"`
	UnderOrganization      bool                   `json:"underOrganization"`
	IsOrganizationVerified bool                   `json:"isOrganizationVerified"`
	SladeCode              string                 `json:"sladeCode"`
	ParentOrganizationID   string                 `json:"parentOrganizationID"`
	HasBranches            bool                   `json:"hasBranches,omitempty"`
	Location               *Location              `json:"location,omitempty"`
	PartnerType            PartnerType            `json:"partnerType"`
	EDIUserProfile         *base.EDIUserProfile   `json:"ediuserprofile" firestore:"ediuserprofile"`
	PartnerSetupComplete   bool                   `json:"partnerSetupComplete" firestore:"partnerSetupComplete"`
	KYCSubmitted           bool                   `json:"kycSubmitted" firestore:"kycSubmitted"`
}

// Location is used to store a user's branch or organisation
type Location struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	BranchSladeCode *string `json:"branchSladeCode"`
}

// Branch represents a Slade 360 Charge Master branch
type Branch struct {
	base.Model

	ID                    string `json:"id"`
	Name                  string `json:"name"`
	OrganizationSladeCode string `json:"organizationSladeCode"`
	BranchSladeCode       string `json:"branchSladeCode"`
}

// BranchEdge is used to serialize GraphQL Relay edges for locations
type BranchEdge struct {
	Cursor *string `json:"cursor"`
	Node   *Branch `json:"node"`
}

// BranchConnection is used tu serialize GraphQL Relay connections for locations
type BranchConnection struct {
	Edges    []*BranchEdge  `json:"edges"`
	PageInfo *base.PageInfo `json:"pageInfo"`
}

// BranchFilterInput is used to supply filter parameters for locatioon list queries
type BranchFilterInput struct {
	Search               *string `json:"search"`
	SladeCode            *string `json:"sladeCode"`
	ParentOrganizationID *string `json:"parentOrganizationID"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BranchFilterInput) ToURLValues() url.Values {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.SladeCode != nil {
		vals.Add("slade_code", *i.SladeCode)
	}
	if i.ParentOrganizationID != nil {
		vals.Add("parent", *i.ParentOrganizationID)
	}
	return vals
}

// KYCRequest the to payload stage kyc processing request
type KYCRequest struct {
	ID                  string                 `json:"id" firestore:"id"`
	ReqPartnerType      PartnerType            `json:"reqPartnerType" firestore:"reqPartnerType"`
	ReqOrganizationType OrganizationType       `json:"reqOrganizationType" firestore:"reqOrganizationType"`
	ReqRaw              map[string]interface{} `json:"reqRaw" firestore:"reqRaw"`
	Proceseed           bool                   `json:"proceseed" firestore:"proceseed"`
	SupplierRecord      *Supplier              `json:"supplierRecord" firestore:"supplierRecord"`
	Status              KYCProcessStatus       `json:"status" firestore:"status"`
	RejectionReason     *string                `json:"rejectionRejection" firestore:"rejectionRejection"`
}

// UserResponse ...
type UserResponse struct {
}

// AuthCredentialResponse ...
type AuthCredentialResponse struct {
}

// BusinessPartner represents a Slade 360 Charge Master business partner
type BusinessPartner struct {
	base.Model

	ID        string  `json:"id"`
	Name      string  `json:"name"`
	SladeCode string  `json:"slade_code"`
	Parent    *string `json:"parent"`
}

// BusinessPartnerEdge is used to serialize GraphQL Relay edges for organization
type BusinessPartnerEdge struct {
	Cursor *string          `json:"cursor"`
	Node   *BusinessPartner `json:"node"`
}

// BusinessPartnerConnection is used to serialize GraphQL Relay connections for organizations
type BusinessPartnerConnection struct {
	Edges    []*BusinessPartnerEdge `json:"edges"`
	PageInfo *base.PageInfo         `json:"pageInfo"`
}

// BusinessPartnerFilterInput is used to supply filter parameters for organizatiom filter inputs
type BusinessPartnerFilterInput struct {
	Search    *string `json:"search"`
	Name      *string `json:"name"`
	SladeCode *string `json:"slade_code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BusinessPartnerFilterInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Search != nil {
		vals.Add("search", *i.Search)
	}
	if i.Name != nil {
		vals.Add("name", *i.Name)
	}
	if i.SladeCode != nil {
		vals.Add("slade_code", *i.SladeCode)
	}
	return vals
}

// BusinessPartnerSortInput is used to supply sort input for organization list queries
type BusinessPartnerSortInput struct {
	Name      *base.SortOrder `json:"name"`
	SladeCode *base.SortOrder `json:"slade_code"`
}

// ToURLValues transforms the filter input to `url.Values`
func (i *BusinessPartnerSortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.SladeCode != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("slade_code", "number")
		} else {
			vals.Add("slade_code", "-number")
		}
	}
	return vals
}

// BranchSortInput is used to supply sorting input for location list queries
type BranchSortInput struct {
	Name      *base.SortOrder `json:"name"`
	SladeCode *base.SortOrder `json:"slade_code"`
}

// ToURLValues transforms the sort input to `url.Values`
func (i *BranchSortInput) ToURLValues() (values url.Values) {
	vals := url.Values{}
	if i.Name != nil {
		if *i.Name == base.SortOrderAsc {
			vals.Add("order_by", "name")
		} else {
			vals.Add("order_by", "-name")
		}
	}
	if i.SladeCode != nil {
		if *i.SladeCode == base.SortOrderAsc {
			vals.Add("slade_code", "number")
		} else {
			vals.Add("slade_code", "-number")
		}
	}
	return vals
}

//TODO: restore commented structs when implementing profile missing methods

// // PostVisitSurvey is used to record and retrieve post visit surveys from Firebase
// type PostVisitSurvey struct {
// 	LikelyToRecommend int    `json:"likelyToRecommend" firestore:"likelyToRecommend"`
// 	Criticism         string `json:"criticism" firestore:"criticism"`
// 	Suggestions       string `json:"suggestions" firestore:"suggestions"`

// 	UID       string    `json:"uid" firestore:"uid"`
// 	Timestamp time.Time `json:"timestamp" firestore:"timestamp"`
// }

// // PIN is used to store a PIN (Personal Identifiation Number) associated
// // to a phone number sign up to Firebase
// type PIN struct {
// 	ProfileID string `json:"profile_id" firestore:"profileID"`
// 	MSISDN    string `json:"msisdn,omitempty" firestore:"msisdn"`
// 	PINNumber string `json:"pin_number" firestore:"pin"`
// 	IsValid   bool   `json:"isValid,omitempty" firestore:"isValid"`
// }

// // PinRecovery stores information required in resetting and updating a forgotten pin
// type PinRecovery struct {
// 	MSISDN    string `json:"msisdn" firestore:"msisdn"`
// 	PINNumber string `json:"pin_number" firestore:"PINNumber"`
// 	OTP       string `json:"otp" firestore:"otp"`
// }

// // OtpResponse returns an otp
// type OtpResponse struct {
// 	OTP string `json:"otp"`
// }

// // Customer used to create a customer request payload
// type Customer struct {
// 	UserProfile        base.UserProfile   `json:"userprofile,omitempty" firestore:"userprofile"`
// 	CustomerID         string             `json:"id,omitempty" firestore:"customerid"`
// 	ReceivablesAccount ReceivablesAccount `json:"receivables_account,omitempty"`
// 	Active             bool               `json:"active" firestore:"active"`
// }

// // Beneficiary stores a customer's beneficiary details
// type Beneficiary struct {
// 	Name         string                  `json:"name"`
// 	Msisdns      []string                `json:"msisdns"`
// 	Emails       []string                `json:"emails"`
// 	Relationship BeneficiaryRelationship `json:"relationship"`
// 	DateOfBirth  base.Date               `json:"dateOfBirth"`
// }

// // ReceivablesAccount stores a customer's receivables account info
// type ReceivablesAccount struct {
// 	ID          string `json:"id,omitempty"`
// 	Name        string `json:"name,omitempty"`
// 	IsActive    bool   `json:"is_active,omitempty"`
// 	Number      string `json:"number,omitempty"`
// 	Tag         string `json:"tag,omitempty"`
// 	Description string `json:"description,omitempty"`
// }

// // BeneficiaryInput stores beneficiary input details
// type BeneficiaryInput struct {
// 	Name         string                  `json:"name"`
// 	Msisdns      []string                `json:"msisdns"`
// 	Emails       []string                `json:"emails"`
// 	Relationship BeneficiaryRelationship `json:"relationship"`
// 	DateOfBirth  base.Date               `json:"dateOfBirth"`
// }

// // OtherPractitionerServiceInput ..
// type OtherPractitionerServiceInput struct {
// 	OtherServices []string `json:"otherServices"`
// }

// // PractitionerServiceInput ..
// type PractitionerServiceInput struct {
// 	Services []PractitionerService `json:"services"`
// }

// // ServicesOffered ..
// type ServicesOffered struct {
// 	Services      []PractitionerService `json:"services"`
// 	OtherServices []string              `json:"otherServices"`
// }

// // StatusResponse creates a status response for requests
// type StatusResponse struct {
// 	Status string `json:"status"`
// }

// // BusinessPartnerUID is the user ID used in some inter-service requests
// type BusinessPartnerUID struct {
// 	UID string `json:"uid"`
// }

// // SendRetryOTP is an input struct for generating and
// // sending fallback otp
// type SendRetryOTP struct {
// 	Msisdn    string `json:"msisdn"`
// 	RetryStep int    `json:"retryStep"`
// }

// // UserUIDs is an input of a list of user uids for isc requests
// type UserUIDs struct {
// 	UIDs []string `json:"uids"`
// }

// // CreatedUserResponse represents payload returned after creating a user
// type CreatedUserResponse struct {
// 	UserProfile *base.UserProfile `json:"user_profile"`
// 	CustomToken *string           `json:"custom_token"`
// }

// // CreateUserViaPhoneInput represents input required to create a user via phoneNumber
// type CreateUserViaPhoneInput struct {
// 	MSISDN string `json:"msisdn"`
// }

// // PhoneSignInInput represents input required to sign in a user via phoneNumber
// type PhoneSignInInput struct {
// 	PhoneNumber string `json:"phonenumber"`
// 	Pin         string `json:"pin"`
// }

// // PhoneSignInResponse is a thin payload returned when a user signs in
// // with their phone number
// type PhoneSignInResponse struct {
// 	CustomToken  string `json:"custom_token"`
// 	IDToken      string `json:"id_token"`
// 	RefreshToken string `json:"refresh_token"`
// }

// // OKResp is used to return OK responses in inter-service calls
// type OKResp struct {
// 	Status string `json:"status"`
// }

// // SaveMemberCoverPayload deserializes inter-service requests to save
// // member covers
// type SaveMemberCoverPayload struct {
// 	PayerName      string `json:"payerName"`
// 	MemberName     string `json:"memberName"`
// 	MemberNumber   string `json:"memberNumber"`
// 	PayerSladeCode int    `json:"payerSladeCode"`
// 	UID            string `json:"uid"`
// }

// // SaveResponsePayload is used to return successful save feedback for
// // inter-service calls
// type SaveResponsePayload struct {
// 	SuccessfullySaved bool `json:"successfullySaved"`
// }

// // OTPResponse is used to return the results of requesting an OTP
// // or OTP retry.
// type OTPResponse struct {
// 	OTP string `json:"otp"`
// }

// // PhoneNumberInput is used to deserialize phone numbers sent to
// // inter-service APIs e.g phone number sign up or verify
// type PhoneNumberInput struct {
// 	PhoneNumber string `json:"phoneNumber"`
// }

// // SupplierAccountInput is used when setting up basic/"key" supplier
// // account during onboarding
// type SupplierAccountInput struct {
// 	AccountType            AccountType    `json:"accountType"`
// 	UnderOrganization      bool           `json:"underOrganization"`
// 	IsOrganizationVerified *bool          `json:"isOrganizationVerified"`
// 	SladeCode              *string        `json:"sladeCode"`
// 	ParentOrganizationID   *string        `json:"parentOrganizationID"`
// 	Location               *LocationInput `json:"location,omitempty"`
// }

// // LocationInput is used when setting up a location (branch or parent) for a user
// type LocationInput struct {
// 	ID              string  `json:"id"`
// 	Name            string  `json:"name"`
// 	BranchSladeCode *string `json:"branchSladeCode"`
// }
