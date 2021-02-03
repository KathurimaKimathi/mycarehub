package usecases_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"gitlab.slade360emr.com/go/base"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/application/resources"
	"gitlab.slade360emr.com/go/profile/pkg/onboarding/domain"
)

func TestNHIFUseCaseImpl_AddNHIFDetails(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to fake initialize onboarding interactor: %v", err)
		return
	}

	validInput := resources.NHIFDetailsInput{
		MembershipNumber:          "123456",
		Employment:                domain.EmploymentTypeEmployed,
		NHIFCardPhotoID:           uuid.New().String(),
		IDDocType:                 domain.IdentificationDocTypeMilitary,
		IdentificationCardPhotoID: uuid.New().String(),
		IDNumber:                  "11111111",
	}
	type args struct {
		ctx   context.Context
		input resources.NHIFDetailsInput
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NHIFDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully add NHIF Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: false,
		},
		{
			name: "sad:( fail to add NHIF Details",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get a user profile",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get logged in user",
			args: args{
				ctx:   ctx,
				input: validInput,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "happy:) successfully add NHIF Details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.AddNHIFDetailsFn = func(ctx context.Context, input resources.NHIFDetailsInput, profileID string) (*domain.NHIFDetails, error) {
					return &domain.NHIFDetails{
						ID:               uuid.New().String(),
						ProfileID:        profileID,
						MembershipNumber: "12345",
						IDNumber:         "12345",
					}, nil
				}
			}

			if tt.name == "sad:( fail to add NHIF Details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: "93ca42bb-5cfc-4499-b137-2df4d67b4a21",
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.AddNHIFDetailsFn = func(ctx context.Context, input resources.NHIFDetailsInput, profileID string) (*domain.NHIFDetails, error) {
					return nil, fmt.Errorf("failed to add nhif details")
				}
			}

			if tt.name == "sad:( fail to get a user profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad:( fail to get logged in user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}
			got, err := i.NHIF.AddNHIFDetails(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.AddNHIFDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

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

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}

func TestNHIFUseCaseImpl_NHIFDetails(t *testing.T) {
	ctx := context.Background()

	i, err := InitializeFakeOnboaridingInteractor()
	if err != nil {
		t.Errorf("failed to initialize fake onboarding interactor")
		return
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.NHIFDetails
		wantErr bool
	}{
		{
			name: "happy:) successfully return NHIF details",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
		{
			name: "happy:) fail to return NHIF Details",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get user profile",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
		{
			name: "sad:( fail to get logged in user",
			args: args{
				ctx: ctx,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.name == "happy:) successfully return NHIF details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.GetNHIFDetailsByProfileIDFn = func(ctx context.Context, profileID string) (*domain.NHIFDetails, error) {
					return &domain.NHIFDetails{
						ID:               uuid.New().String(),
						ProfileID:        profileID,
						MembershipNumber: "12345",
						IDNumber:         "12345",
					}, nil
				}
			}

			if tt.name == "happy:) fail to return NHIF Details" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return uuid.New().String(), nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return &base.UserProfile{
						ID: uuid.New().String(),
						VerifiedIdentifiers: []base.VerifiedIdentifier{
							{
								UID: uid,
							},
						},
					}, nil
				}

				fakeRepo.GetNHIFDetailsByProfileIDFn = func(ctx context.Context, profileID string) (*domain.NHIFDetails, error) {
					return nil, fmt.Errorf("failed to get the user's nhif details")
				}
			}

			if tt.name == "sad:( fail to get user profile" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "7e2aea-d29f2c", nil
				}

				fakeRepo.GetUserProfileByUIDFn = func(ctx context.Context, uid string, suspended bool) (*base.UserProfile, error) {
					return nil, fmt.Errorf("failed to get user profile")
				}
			}

			if tt.name == "sad:( fail to get logged in user" {
				fakeBaseExt.GetLoggedInUserUIDFn = func(ctx context.Context) (string, error) {
					return "", fmt.Errorf("failed to get logged in user")
				}
			}

			got, err := i.NHIF.NHIFDetails(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NHIFUseCaseImpl.NHIFDetails() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

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

				if got == nil {
					t.Errorf("nil response returned")
					return
				}
			}
		})
	}
}
