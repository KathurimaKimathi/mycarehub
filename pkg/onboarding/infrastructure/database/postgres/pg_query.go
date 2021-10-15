package postgres

import (
	"context"
	"fmt"
	"strconv"

	"github.com/google/uuid"

	"github.com/savannahghi/onboarding-service/pkg/onboarding/domain"
)

//GetFacilities returns a slice of healthcare facilities in the platform.
func (d *OnboardingDb) GetFacilities(ctx context.Context) ([]*domain.Facility, error) {
	var facility []*domain.Facility
	facilities, err := d.query.GetFacilities(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get facilities: %v", err)
	}

	if len(facilities) == 0 {
		return facility, nil
	}
	for _, m := range facilities {
		active, err := strconv.ParseBool(m.Active)
		if err != nil {
			return nil, fmt.Errorf("failed to parse facility.Active to boolean")
		}
		singleFacility := domain.Facility{
			ID:          *m.FacilityID,
			Name:        m.Name,
			Code:        m.Code,
			Active:      active,
			County:      m.County,
			Description: m.Description,
		}

		facility = append(facility, &singleFacility)
	}

	return facility, nil
}

// RetrieveFacility gets a facility by ID from the database
func (d *OnboardingDb) RetrieveFacility(ctx context.Context, id *uuid.UUID, isActive bool) (*domain.Facility, error) {
	if id == nil {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacility(ctx, id, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve one facility: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}

// RetrieveByFacilityMFLCode gets a facility by ID from the database
func (d *OnboardingDb) RetrieveByFacilityMFLCode(ctx context.Context, MFLCode string, isActive bool) (*domain.Facility, error) {
	if MFLCode == "" {
		return nil, fmt.Errorf("facility ID should be defined")
	}
	facilitySession, err := d.query.RetrieveFacilityByMFLCode(ctx, MFLCode, isActive)
	if err != nil {
		return nil, fmt.Errorf("failed query and retrieve facility by MFLCode: %s", err)
	}

	return d.mapFacilityObjectToDomain(facilitySession), nil
}