// Package interactor represent reusable chunks of code that abstract
// logic from presenters while simplifying your app and making future changes effortless.
package interactor

import (
	pg "github.com/savannahghi/onboarding-service/pkg/onboarding/infrastructure/database/postgres"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/client"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/facility"
	metrics "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/metric"
	"github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/staff"
	userusecase "github.com/savannahghi/onboarding-service/pkg/onboarding/usecases/user"
	"github.com/savannahghi/onboarding/pkg/onboarding/infrastructure"
	libOnboardingUsecase "github.com/savannahghi/onboarding/pkg/onboarding/usecases"
)

// Interactor represents an assemble of all use cases into a single object that can be instantiated anywhere
type Interactor struct {
	OpenSourceInfra    infrastructure.Infrastructure
	database           pg.OnboardingDb
	OpenSourceUsecases libOnboardingUsecase.Interactor
	FacilityUsecase    facility.UseCasesFacility
	MetricUsecase      metrics.UsecasesMetrics
	UserUsecase        userusecase.UseCasesUser
	StaffUsecase       staff.UsecasesStaffProfile
	ClientUseCase      client.UseCasesClientProfile
}

// NewOnboardingInteractor returns a new onboarding interactor
func NewOnboardingInteractor(
	openSourceInfra infrastructure.Infrastructure,
	database pg.OnboardingDb,
	openSourceUsecases libOnboardingUsecase.Interactor,
	facilityUseCase facility.UseCasesFacility,
	metricUsecase metrics.UsecasesMetrics,
	usecaseUser userusecase.UseCasesUser,
	staffProfileUsecase staff.UsecasesStaffProfile,
	clientUseCase client.UseCasesClientProfile,
) *Interactor {
	return &Interactor{
		OpenSourceInfra:    openSourceInfra,
		database:           database,
		OpenSourceUsecases: openSourceUsecases,
		FacilityUsecase:    facilityUseCase,
		MetricUsecase:      metricUsecase,
		StaffUsecase:       staffProfileUsecase,
		UserUsecase:        usecaseUser,
		ClientUseCase:      clientUseCase,
	}
}
