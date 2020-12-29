package errors

const (
	// UserNotFoundErrMsg is the error message displayed when a user is not found
	UserNotFoundErrMsg = "failed to get a Firebase user"

	// ProfileNotFoundErrMsg is the error message displayed when a user is not found
	ProfileNotFoundErrMsg = "failed to get a Firebase user profile"

	// PINNotFoundErrMsg is the error message displayed when a pin is not found
	PINNotFoundErrMsg = "failed to get a Firebase user pin"

	// CustomTokenErrMsg is the error message displayed when a
	// custom firebase token is not created
	CustomTokenErrMsg = "failed to create custom token"

	// AuthenticateTokenErrMsg is the error message displayed when a
	// custom firebase token is not authenticated
	AuthenticateTokenErrMsg = "failed to authenticate custom token"

	// UpdateProfileErrMsg is the error message displayed when a
	// user profile is not found
	UpdateProfileErrMsg = "failed to update a user profile"

	// AddRecordErrMsg is the error message displayed when a
	// record fails to be added to the database
	AddRecordErrMsg = "failed to add the record to the database"

	// LikelyToRecommendErrMsg is the error message displayed that
	// occurs when the recommendation threshold is crossed
	LikelyToRecommendErrMsg = "the likelihood of recommending should be an int between 0 and 10"

	// ValidatePINLengthErrMsg  is the error message displayed when
	// an invalid Pin length is given
	ValidatePINLengthErrMsg = "pin should be of 4,5, or six digits"

	// ValidatePINDigitsErrMsg  is the error message displayed when
	// an invalid  pin digits are given
	ValidatePINDigitsErrMsg = "pin should be a valid number"

	// UsePinExistErrMsg  is the error message displayed when
	// user has pin already during set pin
	UsePinExistErrMsg = "the user has PIN already"

	// EncryptPINErrMsg  is the error message displayed when
	// pin encryption failed
	EncryptPINErrMsg = "unable to encrypt PIN"

	// RetrieveRecordErrMsg is the error message displayed when a
	// record fails to be retrieved from database
	RetrieveRecordErrMsg = "unable to retrieve newly created record"
)
