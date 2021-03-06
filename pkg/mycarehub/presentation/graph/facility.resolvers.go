package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/savannahghi/mycarehub/pkg/mycarehub/application/dto"
	"github.com/savannahghi/mycarehub/pkg/mycarehub/domain"
)

func (r *mutationResolver) CreateFacility(ctx context.Context, input dto.FacilityInput) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.GetOrCreateFacility(ctx, &input)
}

func (r *mutationResolver) DeleteFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.DeleteFacility(ctx, mflCode)
}

func (r *mutationResolver) ReactivateFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.ReactivateFacility(ctx, &mflCode)
}

func (r *mutationResolver) InactivateFacility(ctx context.Context, mflCode int) (bool, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.InactivateFacility(ctx, &mflCode)
}

func (r *queryResolver) FetchFacilities(ctx context.Context) ([]*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.FetchFacilities(ctx)
}

func (r *queryResolver) RetrieveFacility(ctx context.Context, id string, active bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacility(ctx, &id, active)
}

func (r *queryResolver) RetrieveFacilityByMFLCode(ctx context.Context, mflCode int, isActive bool) (*domain.Facility, error) {
	r.checkPreconditions()
	return r.mycarehub.Facility.RetrieveFacilityByMFLCode(ctx, mflCode, isActive)
}

func (r *queryResolver) ListFacilities(ctx context.Context, searchTerm *string, filterInput []*dto.FiltersInput, paginationInput dto.PaginationsInput) (*domain.FacilityPage, error) {
	return r.mycarehub.Facility.ListFacilities(ctx, searchTerm, filterInput, &paginationInput)
}
