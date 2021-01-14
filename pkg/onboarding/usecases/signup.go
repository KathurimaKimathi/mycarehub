package usecases

import (
	"context"
	"fmt"

	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/exceptions"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/extension"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/utils"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/infrastructure/services/otp"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/repository"
)

// SignUpUseCases represents all the business logic involved in setting up a user
type SignUpUseCases interface {
	// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
	VerifyPhoneNumber(ctx context.Context, phone string) (*base.OtpResponse, error)
	// checks whether a phone number has been registred by another user. Checks both primary and
	// secondary phone numbers. If the the phone number is foreign, it send an OTP to that phone number
	CheckPhoneExists(ctx context.Context, phone string) (bool, error)

	// creates an account for the user, setting the provided phone number as the PRIMARY PHONE NUMBER
	CreateUserByPhone(ctx context.Context, input *resources.SignUpInput) (*base.UserResponse, error)

	// updates the user profile of the currently logged in user
	UpdateUserProfile(ctx context.Context, input *resources.UserProfileInput) (*base.UserProfile, error)

	// adds a new push token in the users profile if the push token does not exist
	RegisterPushToken(ctx context.Context, token string) (bool, error)

	// called to create a customer account in the ERP. This API is only valid for `BEWELL CONSUMER`
	// it should be the last call after updating the users bio data. Its should not return an error
	// when it fails due to unreachable errors, rather it should retry
	CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error)

	// removes a push token from the users profile
	RetirePushToken(ctx context.Context, token string) (bool, error)

	// fetches the phone numbers of a user for the purposes of recoverying an account.
	// the returned phone numbers should be masked
	GetUserRecoveryPhoneNumbers(ctx context.Context, phoneNumber string) (*resources.AccountRecoveryPhonesResponse, error)

	// called to set the provided phone number as the PRIMARY PHONE NUMBER in the user profile of the user
	// where the phone number is associated with.
	SetPhoneAsPrimary(ctx context.Context, phone, otp string) (bool, error)

	RemoveUserByPhoneNumber(ctx context.Context, phone string) error
}

// SignUpUseCasesImpl represents usecase implementation object
type SignUpUseCasesImpl struct {
	onboardingRepository repository.OnboardingRepository
	profileUsecase       ProfileUseCase
	pinUsecase           UserPINUseCases
	supplierUsecase      SupplierUseCases
	otpUseCases          otp.ServiceOTP
	baseExt              extension.BaseExtension
}

// NewSignUpUseCases returns a new a onboarding usecase
func NewSignUpUseCases(
	r repository.OnboardingRepository,
	profile ProfileUseCase,
	pin UserPINUseCases,
	supplier SupplierUseCases,
	otp otp.ServiceOTP,
	ext extension.BaseExtension,
) SignUpUseCases {
	return &SignUpUseCasesImpl{
		onboardingRepository: r,
		profileUsecase:       profile,
		pinUsecase:           pin,
		supplierUsecase:      supplier,
		otpUseCases:          otp,
		baseExt:              ext,
	}
}

// CheckPhoneExists checks whether a phone number has been registred by another user.
// Checks both primary and secondary phone numbers.
func (s *SignUpUseCasesImpl) CheckPhoneExists(ctx context.Context, phone string) (bool, error) {
	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}
	exists, err := s.onboardingRepository.CheckIfPhoneNumberExists(ctx, *phoneNumber)
	if err != nil {
		return false, exceptions.CheckPhoneNumberExistError()
	}
	return exists, nil
}

// CreateUserByPhone creates an account for the user, setting the provided phone number as the PRIMARY PHONE NUMBER
func (s *SignUpUseCasesImpl) CreateUserByPhone(
	ctx context.Context,
	input *resources.SignUpInput,
) (*base.UserResponse, error) {
	userData, err := utils.ValidateSignUpInput(input)
	if err != nil {
		return nil, err
	}
	verified, err := s.otpUseCases.VerifyOTP(
		ctx,
		*userData.PhoneNumber,
		*userData.OTP,
	)
	if err != nil {
		return nil, exceptions.VerifyOTPError(err)
	}

	if !verified {
		return nil, exceptions.VerifyOTPError(nil)
	}
	// check if phone number is registered to another user
	exists, err := s.CheckPhoneExists(ctx, *userData.PhoneNumber)
	if err != nil {
		return nil, err
	}
	// if phone exists return early
	if exists {
		return nil, exceptions.CheckPhoneNumberExistError()
	}
	// get or create user via their phone number
	user, err := s.onboardingRepository.GetOrCreatePhoneNumberUser(ctx, *userData.PhoneNumber)
	if err != nil {
		return nil, err
	}
	// create a user profile
	profile, err := s.onboardingRepository.CreateUserProfile(
		ctx,
		*userData.PhoneNumber,
		user.UID,
	)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}
	// generate auth credentials
	auth, err := s.onboardingRepository.GenerateAuthCredentials(
		ctx,
		*userData.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}
	// save the user pin
	_, err = s.pinUsecase.SetUserPIN(
		ctx,
		*userData.PIN,
		*userData.PhoneNumber,
	)
	if err != nil {
		return nil, err
	}

	var supplier *base.Supplier
	var customer *base.Customer

	supplier, err = s.onboardingRepository.CreateEmptySupplierProfile(ctx, profile.ID)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	customer, err = s.onboardingRepository.CreateEmptyCustomerProfile(ctx, profile.ID)
	if err != nil {
		return nil, exceptions.InternalServerError(err)
	}

	return &base.UserResponse{
		Profile:         profile,
		SupplierProfile: supplier,
		CustomerProfile: customer,
		Auth:            *auth,
	}, nil
}

// UpdateUserProfile  updates the user profile of the currently logged in user
func (s *SignUpUseCasesImpl) UpdateUserProfile(ctx context.Context, input *resources.UserProfileInput) (*base.UserProfile, error) {

	// get the old user profile
	pr, err := s.profileUsecase.UserProfile(ctx)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}

	if input.PhotoUploadID != nil {
		if err := s.profileUsecase.UpdatePhotoUploadID(ctx, *input.PhotoUploadID); err != nil {
			return nil, err
		}
	}

	if err := s.profileUsecase.UpdateBioData(ctx, base.BioData{
		FirstName: func(n *string) *string {
			if n != nil {
				return n
			}
			return pr.UserBioData.FirstName
		}(input.FirstName),
		LastName: func(n *string) *string {
			if n != nil {
				return n
			}
			return pr.UserBioData.LastName
		}(input.LastName),
		DateOfBirth: func(n *base.Date) *base.Date {
			if n != nil {
				return n
			}
			return pr.UserBioData.DateOfBirth
		}(input.DateOfBirth),
		Gender: func(n *base.Gender) base.Gender {
			if n != nil {
				return *n
			}
			return pr.UserBioData.Gender
		}(input.Gender),
	}); err != nil {
		return nil, err
	}
	return s.profileUsecase.UserProfile(ctx)
}

// RegisterPushToken adds a new push token in the users profile if the push token does not exist
func (s *SignUpUseCasesImpl) RegisterPushToken(ctx context.Context, token string) (bool, error) {
	if len(token) < 5 {
		return false, exceptions.InValidPushTokenLengthError()
	}
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, false); err != nil {
		return false, err
	}
	return true, nil
}

// CompleteSignup called to create a customer account in the ERP. This API is only valid for `BEWELL CONSUMER`
func (s *SignUpUseCasesImpl) CompleteSignup(ctx context.Context, flavour base.Flavour) (bool, error) {

	if flavour == base.FlavourConsumer {
		pr, err := s.profileUsecase.UserProfile(ctx)
		if err != nil {
			return false, exceptions.ProfileNotFoundError(err)
		}
		if pr.UserBioData.FirstName == nil || pr.UserBioData.LastName == nil {
			return false, exceptions.CompleteSignUpError(nil)
		}
		fullName := fmt.Sprintf("%v %v", *pr.UserBioData.FirstName, *pr.UserBioData.LastName)
		_, _ = s.supplierUsecase.AddCustomerSupplierERPAccount(ctx, fullName, base.PartnerTypeConsumer)
		return true, nil
	}
	return false, exceptions.InvalidFlavourDefinedError()
}

// RetirePushToken removes a push token from the users profile
func (s *SignUpUseCasesImpl) RetirePushToken(ctx context.Context, token string) (bool, error) {
	if len(token) < 5 {
		return false, exceptions.InValidPushTokenLengthError()
	}
	if err := s.profileUsecase.UpdatePushTokens(ctx, token, true); err != nil {
		return false, exceptions.InternalServerError(err)
	}
	return true, nil
}

// GetUserRecoveryPhoneNumbers fetches the phone numbers of a user for the purposes of recoverying an account.
func (s *SignUpUseCasesImpl) GetUserRecoveryPhoneNumbers(ctx context.Context, phone string) (*resources.AccountRecoveryPhonesResponse, error) {
	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}

	pr, err := s.onboardingRepository.GetUserProfileByPhoneNumber(ctx, *phoneNumber)
	if err != nil {
		return nil, exceptions.ProfileNotFoundError(err)
	}

	// cherrypick the phone numbers and mask them
	phones := func(p *base.UserProfile) []string {
		phs := []string{}
		phs = append(phs, *p.PrimaryPhone)
		phs = append(phs, p.SecondaryPhoneNumbers...)
		return phs

	}(pr)

	masked := s.profileUsecase.MaskPhoneNumbers(phones)

	return &resources.AccountRecoveryPhonesResponse{
		MaskedPhoneNumbers:   masked,
		UnMaskedPhoneNumbers: phones,
	}, nil
}

// SetPhoneAsPrimary called to set the provided phone number as the PRIMARY PHONE NUMBER in the user profile of the user
// where the phone number is associated with.
func (s *SignUpUseCasesImpl) SetPhoneAsPrimary(ctx context.Context, phone, otp string) (bool, error) {
	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return false, exceptions.NormalizeMSISDNError(err)
	}

	if err := s.profileUsecase.SetPrimaryPhoneNumber(ctx, *phoneNumber, otp, false); err != nil {
		return false, err
	}
	return true, nil
}

// RemoveUserByPhoneNumber removes the record of a user using the provided phone number. This method will ONLY be called
// in testing environment.
func (s *SignUpUseCasesImpl) RemoveUserByPhoneNumber(ctx context.Context, phone string) error {
	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return exceptions.NormalizeMSISDNError(err)
	}
	return s.onboardingRepository.PurgeUserByPhoneNumber(ctx, *phoneNumber)
}

// VerifyPhoneNumber checks validity of a phone number by sending an OTP to it
func (s *SignUpUseCasesImpl) VerifyPhoneNumber(ctx context.Context, phone string) (*base.OtpResponse, error) {
	phoneNumber, err := s.baseExt.NormalizeMSISDN(phone)
	if err != nil {
		return nil, exceptions.NormalizeMSISDNError(err)
	}
	// check if phone number exists
	exists, err := s.CheckPhoneExists(ctx, *phoneNumber)
	if err != nil {
		return nil, err
	}
	// if phone exists return early
	if exists {
		return nil, exceptions.CheckPhoneNumberExistError()
	}
	// generate and send otp to the phone number
	otpResp, err := s.otpUseCases.GenerateAndSendOTP(ctx, *phoneNumber)

	if err != nil {
		return nil, exceptions.GenerateAndSendOTPError(err)
	}
	// return the generated otp
	return otpResp, nil
}
