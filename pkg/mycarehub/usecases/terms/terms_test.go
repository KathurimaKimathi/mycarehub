package terms_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/savannahghi/enumutils"
	"github.com/savannahghi/feedlib"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/common/testutils"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/enums"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/infrastructure/database/postgres/gorm"
	"github.com/segmentio/ksuid"
)

var (
	termsID   = 50005
	orgID     = uuid.New().String()
	testingDB *gorm.PGInstance
)

func TestMain(m *testing.M) {
	log.Println("setting up test database")
	var err error
	testingDB, err = gorm.NewPGInstance()
	if err != nil {
		os.Exit(1)
	}
	// add organization
	createOrganization()

	//create terms
	createTermsOfService()

	log.Printf("Running tests ...")
	os.Exit(m.Run())

	// teardown
	// remove organization
	log.Println("tearing down test database")

	testingDB.DB.Unscoped().Delete(gorm.Organisation{OrganisationID: &orgID})
	testingDB.DB.Unscoped().Delete(gorm.TermsOfService{TermsID: &termsID})
}

func createOrganization() {
	organisation := &gorm.Organisation{
		OrganisationID:   &orgID,
		Active:           true,
		OrgCode:          gofakeit.Name(),
		Code:             gofakeit.Number(100, 344),
		OrganisationName: gofakeit.Name(),
		EmailAddress:     gofakeit.Email(),
		PhoneNumber:      gofakeit.Phone(),
		PostalAddress:    gofakeit.Address().Address,
		PhysicalAddress:  gofakeit.Address().City,
		DefaultCountry:   "KEN",
	}

	testingDB.DB.Create(organisation)
}

func createTermsOfService() {
	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	terms := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
	}

	testingDB.DB.Create(terms)
}

func TestServiceTermsImpl_GetCurrentTerms_Integration(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)

	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	termsID := 2000002
	termsInput := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
		Active:    true,
	}

	err = pg.DB.Create(&termsInput).Error
	if err != nil {
		t.Errorf("failed to create terms: %v", err)
	}

	flavour := feedlib.FlavourConsumer
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     termsInput.TermsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         true,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}
	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *domain.TermsOfService
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx: ctx,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Terms.GetCurrentTerms(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTermsImpl.GetCurrentTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && got != nil {
				t.Errorf("expected facilities to be nil for %v", tt.name)
				return
			}

			if !tt.wantErr && got == nil {
				t.Errorf("expected facilities not to be nil for %v", tt.name)
				return
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", termsInput.TermsID).Unscoped().Delete(&gorm.TermsOfService{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}

func TestServiceTermsImpl_AcceptTerms_Integration_test(t *testing.T) {
	ctx := context.Background()

	i, _ := testutils.InitializeTestService(ctx)
	pg, err := gorm.NewPGInstance()
	if err != nil {
		t.Errorf("pgInstance.Teardown() = %v", err)
	}

	validFrom := time.Now()
	validTo := time.Now().AddDate(0, 0, 50)
	txt := gofakeit.HipsterSentence(15)
	termsID := 2010001
	termsInput := &gorm.TermsOfService{
		TermsID:   &termsID,
		Text:      &txt,
		ValidFrom: &validFrom,
		ValidTo:   &validTo,
	}

	err = pg.DB.Create(&termsInput).Error
	if err != nil {
		t.Errorf("failed to create terms: %v", err)
	}

	flavour := feedlib.FlavourConsumer
	currentTime := time.Now()
	nextTime := time.Now().AddDate(0, 0, 2)

	// Setup test user
	userInput := &gorm.User{
		Username:            uuid.New().String(),
		FirstName:           gofakeit.FirstName(),
		MiddleName:          gofakeit.FirstName(),
		LastName:            gofakeit.LastName(),
		UserType:            enums.ClientUser,
		Gender:              enumutils.GenderMale,
		Active:              false,
		PushTokens:          []string{},
		LastSuccessfulLogin: &currentTime,
		LastFailedLogin:     &currentTime,
		FailedLoginCount:    0,
		NextAllowedLogin:    &nextTime,
		TermsAccepted:       true,
		AcceptedTermsID:     termsInput.TermsID,
		Flavour:             flavour,
		Avatar:              "",
		IsSuspended:         false,
		OrganisationID:      orgID,
		Password:            "",
		IsSuperuser:         false,
		IsStaff:             false,
		Email:               "",
		DateJoined:          "",
		Name:                "",
		IsApproved:          false,
		ApprovalNotified:    false,
		Handle:              "",
	}

	err = pg.DB.Create(&userInput).Error
	if err != nil {
		t.Errorf("failed to create user: %v", err)
	}

	invalidUserID := ksuid.New().String()
	edgecaseUserID := gofakeit.HipsterParagraph(1, 50, 100, "")
	invalidTermsID := 0

	type args struct {
		ctx     context.Context
		userID  *string
		termsID *int
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Happy case",
			args: args{
				ctx:     ctx,
				userID:  userInput.UserID,
				termsID: termsInput.TermsID,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Sad case - invalid user ID",
			args: args{
				ctx:     ctx,
				userID:  &invalidUserID,
				termsID: termsInput.TermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid terms ID",
			args: args{
				ctx:     ctx,
				userID:  userInput.UserID,
				termsID: &invalidTermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - invalid user ID and terms ID",
			args: args{
				ctx:     ctx,
				userID:  &invalidUserID,
				termsID: &invalidTermsID,
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Sad case - edge user ID",
			args: args{
				ctx:     ctx,
				userID:  &edgecaseUserID,
				termsID: &invalidTermsID,
			},
			want:    false,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := i.Terms.AcceptTerms(tt.args.ctx, tt.args.userID, tt.args.termsID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ServiceTermsImpl.AcceptTerms() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ServiceTermsImpl.AcceptTerms() = %v, want %v", got, tt.want)
			}
		})
	}

	//TearDown
	if err = pg.DB.Where("id", userInput.UserID).Unscoped().Delete(&gorm.User{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
	if err = pg.DB.Where("id", termsInput.TermsID).Unscoped().Delete(&gorm.TermsOfService{}).Error; err != nil {
		t.Errorf("failed to delete record = %v", err)
	}
}