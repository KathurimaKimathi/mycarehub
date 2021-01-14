package usecases_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"firebase.google.com/go/auth"
	"github.com/stretchr/testify/assert"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
)

func TestMaskPhoneNumbers(t *testing.T) {

	ctx := context.Background()
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		phones []string
	}

	tests := []struct {
		name string
		arg  args
		want []string
	}{
		{
			name: "valid case",
			arg: args{
				phones: []string{"+254789874267"},
			},
			want: []string{"+254789***267"},
		},
		{
			name: "valid case < 10 digits",
			arg: args{
				phones: []string{"+2547898742"},
			},
			want: []string{"+2547***742"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			maskedPhone := s.Onboarding.MaskPhoneNumbers(tt.arg.phones)
			if len(maskedPhone) != len(tt.want) {
				t.Errorf("returned masked phone number not the expected one, wanted: %v got: %v", tt.want, maskedPhone)
				return
			}

			for i, number := range maskedPhone {
				if tt.want[i] != number {
					t.Errorf("wanted: %v, got: %v", tt.want[i], number)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_GetUserProfileByUID(t *testing.T) {
	ctx, auth, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("failed to get test authenticated context: %v", err)
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}
	type args struct {
		ctx context.Context
		UID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sucess: get a user profile given their UID",
			args: args{
				ctx: ctx,
				UID: auth.UID,
			},
			wantErr: false,
		},
		{
			name: "failure: fail get a user profile given a bad UID",
			args: args{
				ctx: ctx,
				UID: "not-a-valid-uid",
			},
			wantErr: true,
		},
		{
			name: "failure: fail get a user profile given an empty UID",
			args: args{
				ctx: ctx,
				UID: "",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := s.Onboarding.GetUserProfileByUID(tt.args.ctx, tt.args.UID)
			if tt.wantErr && profile != nil {
				t.Errorf("expected nil but got %v, since the error %v occurred",
					profile,
					err,
				)
				return
			}

			if !tt.wantErr && profile == nil {
				t.Errorf("expected a profile but got nil, since no error occurred")
				return
			}

		})
	}
}

func TestProfileUseCaseImpl_UserProfile(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}
	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *base.UserProfile
		wantErr bool
	}{
		{
			name: "valid: user profile retrieved",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "invalid: unauthenticated context supplied",
			args: args{
				ctx: context.Background(),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.UserProfile(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UserProfile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_GetProfileByID(t *testing.T) {

	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	profile, err := s.Onboarding.UserProfile(ctx)
	if err != nil {
		t.Errorf("could not retreive user profile")
		return
	}

	type args struct {
		ctx context.Context
		id  *string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid: user profile retreived",
			args: args{
				ctx: ctx,
				id:  &profile.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Onboarding.GetProfileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.GetProfileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantErr {
				t.Errorf("nil user profile returned")
				return
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateBioData(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	dateOfBirth := base.Date{
		Day:   12,
		Year:  2000,
		Month: 2,
	}

	firstName := "Jatelo"
	lastName := "Omera"
	bioData := base.BioData{
		FirstName:   &firstName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
	}

	var gender base.Gender = "female"
	updateGender := base.BioData{
		Gender: gender,
	}

	updateDOB := base.BioData{
		DateOfBirth: &dateOfBirth,
	}

	updateFirstName := base.BioData{
		FirstName: &firstName,
	}

	updateLastName := base.BioData{
		LastName: &lastName,
	}

	type args struct {
		ctx  context.Context
		data base.BioData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update biodata",
			args: args{
				ctx:  ctx,
				data: bioData,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the firstname",
			args: args{
				ctx:  ctx,
				data: updateFirstName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the lastname",
			args: args{
				ctx:  ctx,
				data: updateLastName,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the date of birth",
			args: args{
				ctx:  ctx,
				data: updateDOB,
			},
			wantErr: false,
		},
		{
			name: "Happy case - Successfully update the gender",
			args: args{
				ctx:  ctx,
				data: updateGender,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Unauthenticated context",
			args: args{
				ctx:  context.Background(),
				data: bioData,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdateBioData(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdateBioData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePhotoUploadID(t *testing.T) {
	ctx, _, err := GetTestAuthenticatedContext(t)
	if err != nil {
		t.Errorf("could not get test authenticated context")
		return
	}

	s, err := InitializeTestService(ctx)
	if err != nil {
		t.Errorf("unable to initialize test service")
		return
	}

	uid, err := base.GetLoggedInUserUID(ctx)
	if err != nil {
		t.Errorf("could not get the logged in user")
		return
	}

	profile, err := s.Onboarding.GetUserProfileByUID(ctx, uid)
	if err != nil {
		t.Errorf("could not retrieve user profile")
		return
	}

	uploadID := "some-photo-upload-id"
	log.Printf("THE UPLOAD ID IS %v", profile.PhotoUploadID)

	type args struct {
		ctx      context.Context
		uploadID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Happy Case - Successfully update the PhotoUploadID",
			args: args{
				ctx:      ctx,
				uploadID: uploadID,
			},
			wantErr: false,
		},
		{
			name: "Sad Case - Use an unauthenticated context",
			args: args{
				ctx:      context.Background(),
				uploadID: uploadID,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := s.Onboarding.UpdatePhotoUploadID(tt.args.ctx, tt.args.uploadID); (err != nil) != tt.wantErr {
				t.Errorf("ProfileUseCaseImpl.UpdatePhotoUploadID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetPhoneAsPrimary(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	secondaryPhone := base.TestUserPhoneNumberWithPin
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), secondaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	// add a secondary phone number to the user
	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))

	// login to add assert the secondary phone number has been added
	login2, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login2)
	assert.Equal(t, 1, len(login2.Profile.SecondaryPhoneNumbers))

	// send otp to the secondary phone number we intend to make primary
	otpResp, err := s.Otp.GenerateAndSendOTP(context.Background(), secondaryPhone)
	assert.Nil(t, err)
	assert.NotNil(t, otpResp)

	// set the old secondary phone number as the new primary phone number
	setResp, err := s.Signup.SetPhoneAsPrimary(context.Background(), secondaryPhone, otpResp.OTP)
	assert.Nil(t, err)
	assert.NotNil(t, setResp)

	// login with the old primary phone number. This should fail
	login3, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login3)

	// login with the new primary phone number. This should not fail
	login4, err := s.Login.LoginByPhone(context.Background(), secondaryPhone, pin, base.FlavourConsumer)
	assert.NotNil(t, err)
	assert.Nil(t, login4)

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), secondaryPhone)
}

func TestAddSecondaryPhoneNumbers(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	secondaryPhone1 := base.TestUserPhoneNumberWithPin
	secondaryPhone2 := "+25712345690"
	secondaryPhone3 := "+25710375600"

	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone1})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.SecondaryPhoneNumbers))

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone2})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.SecondaryPhoneNumbers))

	err = s.Onboarding.UpdateSecondaryPhoneNumbers(authenticatedContext, []string{secondaryPhone3})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 3, len(pr.SecondaryPhoneNumbers))
}

func TestUserProfileCovers(t *testing.T) {
	s, err := InitializeTestService(context.Background())
	if err != nil {
		t.Error("failed to setup signup usecase")
	}
	primaryPhone := base.TestUserPhoneNumber
	// clean up
	_ = s.Signup.RemoveUserByPhoneNumber(context.Background(), primaryPhone)

	otp, err := generateTestOTP(t, primaryPhone)
	if err != nil {
		t.Errorf("failed to generate test OTP: %v", err)
		return
	}
	pin := "1234"
	resp, err := s.Signup.CreateUserByPhone(
		context.Background(),
		&resources.SignUpInput{
			PhoneNumber: &primaryPhone,
			PIN:         &pin,
			Flavour:     base.FlavourConsumer,
			OTP:         &otp.OTP,
		},
	)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Profile)
	assert.NotNil(t, resp.CustomerProfile)
	assert.NotNil(t, resp.SupplierProfile)

	login1, err := s.Login.LoginByPhone(context.Background(), primaryPhone, pin, base.FlavourConsumer)
	assert.Nil(t, err)
	assert.NotNil(t, login1)

	// create authenticated context
	ctx := context.Background()
	authCred := &auth.Token{UID: login1.Auth.UID}
	authenticatedContext := context.WithValue(
		ctx,
		base.AuthTokenContextKey,
		authCred,
	)
	s, _ = InitializeTestService(authenticatedContext)

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer1", PayerSladeCode: 1, MemberName: "name1", MemberNumber: "mem1"}})
	assert.Nil(t, err)

	pr, err := s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 1, len(pr.Covers))

	err = s.Onboarding.UpdateCovers(authenticatedContext, []base.Cover{{PayerName: "payer2", PayerSladeCode: 2, MemberName: "name2", MemberNumber: "mem2"}})
	assert.Nil(t, err)

	pr, err = s.Onboarding.UserProfile(authenticatedContext)
	assert.Nil(t, err)
	assert.NotNil(t, pr)
	assert.Equal(t, 2, len(pr.Covers))

}

func TestProfileUseCaseImpl_UpdateVerifiedUIDS(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx  context.Context
		uids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_profile_uids",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: false,
		},

		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:  ctx,
				uids: []string{"f4f39af7-5b64-4c2f-91bd-42b3af315a4e", "5d46d3bd-a482-4787-9b87-3c94510c8b53"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_profile_uids" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeBaseExt.NormalizeMSISDNFn = func(msisdn string) (*string, error) {
					phone := "+254721123123"
					return &phone, nil
				}
				fakeRepo.UpdateVerifiedUIDSFn = func(ctx context.Context, id string, uids []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateVerifiedUIDS(tt.args.ctx, tt.args.uids)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_UpdateSecondaryEmailAddresses(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	type args struct {
		ctx            context.Context
		emailAddresses []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_profile_secondary_email",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:            ctx,
				emailAddresses: []string{"me4@gmail.com", "kalulu@gmail.com"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_profile_secondary_email" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, uids []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateSecondaryEmailAddresses(tt.args.ctx, tt.args.emailAddresses)

			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateUserName(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx      context.Context
		userName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_name_suceeds",
			args: args{
				ctx:      ctx,
				userName: "kamau",
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:      ctx,
				userName: "mwas",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_name_suceeds" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateUserNameFn = func(ctx context.Context, id string, phoneNumber string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}
			err := i.Onboarding.UpdateUserName(tt.args.ctx, tt.args.userName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdateVerifiedIdentifiers(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	type args struct {
		ctx         context.Context
		identifiers []base.VerifiedIdentifier
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_name_suceeds",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "a4f39af7-5b64-4c2f-91bd-42b3af315a5h",
						LoginProvider: "Facebook",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "j4f39af7-5b64-4c2f-91bd-42b3af225a5c",
						LoginProvider: "Phone",
					},
				},
			},
			wantErr: true,
		},

		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx: ctx,
				identifiers: []base.VerifiedIdentifier{
					{
						UID:           "p4f39af7-5b64-4c2f-91bd-42b3af315a5c",
						LoginProvider: "Google",
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_update_name_suceeds" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
					}, nil
				}
				fakeRepo.UpdateVerifiedIdentifiersFn = func(ctx context.Context, id string, identifiers []base.VerifiedIdentifier) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdateVerifiedIdentifiers(tt.args.ctx, tt.args.identifiers)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}
		})
	}
}

func TestProfileUseCaseImpl_UpdatePrimaryEmailAddress(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	primaryEmail := "me@gmail.com"

	type args struct {
		ctx          context.Context
		emailAddress string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_update_email_suceeds",
			args: args{
				ctx:          ctx,
				emailAddress: primaryEmail,
			},
			wantErr: false,
		},
		{
			name: "invalid:_unable_to_get_logged_in_user",
			args: args{
				ctx:          ctx,
				emailAddress: "kalulu@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_get_profile_of_logged_in_user",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_primary_email_address",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
		{
			name: "invalid:_unable_to_update_secondary_email_address",
			args: args{
				ctx:          ctx,
				emailAddress: "juha@gmail.com",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "valid:_update_email_suceeds" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
			}

			if tt.name == "invalid:_unable_to_update_primary_email_address" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return fmt.Errorf("unable to update primary address")
				}
			}

			if tt.name == "invalid:_unable_to_update_secondary_email_address" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
						SecondaryEmailAddresses: []string{
							"", "lulu@gmail.com",
						},
					}, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return fmt.Errorf("unable to update secondary email")
				}
			}

			if tt.name == "invalid:_unable_to_get_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get logged user")
				}
			}

			if tt.name == "invalid:_unable_to_get_profile_of_logged_in_user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return nil, fmt.Errorf("unable to get profile")
				}
			}

			err := i.Onboarding.UpdatePrimaryEmailAddress(tt.args.ctx, tt.args.emailAddress)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}

func TestProfileUseCaseImpl_SetPrimaryEmailAddress(t *testing.T) {
	ctx := context.Background()
	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}
	primaryEmail := "juha@gmail.com"

	type args struct {
		ctx          context.Context
		emailAddress string
		otp          string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid:_set_primary_address_suceeds",
			args: args{
				ctx:          ctx,
				emailAddress: primaryEmail,
				otp:          "689552",
			},
			wantErr: false,
		},
		{
			name: "valid:_verify_otp_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "kichwa@gmail.com",
				otp:          "453852",
			},
			wantErr: true,
		},
		{
			name: "valid:_verify_otp_returns_false",
			args: args{
				ctx:          ctx,
				emailAddress: "kalu@gmail.com",
				otp:          "235789",
			},
			wantErr: true,
		},
		{
			name: "valid:_update_primary_address_fails",
			args: args{
				ctx:          ctx,
				emailAddress: "mwendwapole@gmail.com",
				otp:          "897523",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "valid:_set_primary_address_suceeds" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return nil
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "5cf354a2-1d3e-400d-8716-7e2aead29f2c", nil
				}
				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID:                  "f4f39af7-5b64-4c2f-91bd-42b3af315a4e",
						PrimaryEmailAddress: &primaryEmail,
					}, nil
				}
				fakeRepo.UpdateSecondaryEmailAddressesFn = func(ctx context.Context, id string, emailAddresses []string) error {
					return nil
				}
			}

			if tt.name == "valid:_verify_otp_fails" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, fmt.Errorf("unable to verify email otp")
				}
			}

			if tt.name == "valid:_verify_otp_returns_false" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return false, nil
				}
			}

			if tt.name == "valid:_update_primary_address_fails" {
				fakeOtp.VerifyEmailOTPFn = func(ctx context.Context, phone, OTP string) (bool, error) {
					return true, nil
				}
				fakeRepo.UpdatePrimaryEmailAddressFn = func(ctx context.Context, id string, emailAddress string) error {
					return fmt.Errorf("unable to update primary email")
				}
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("unable to get loggedin user")
				}
			}

			err := i.Onboarding.SetPrimaryEmailAddress(tt.args.ctx, tt.args.emailAddress, tt.args.otp)
			if tt.wantErr {
				if err == nil {
					t.Errorf("error expected got %v", err)
					return
				}
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("error not expected got %v", err)
					return
				}
			}

		})
	}
}
