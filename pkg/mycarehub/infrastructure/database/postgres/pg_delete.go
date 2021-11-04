package postgres

import "context"

// DeleteFacility does the actual delete of a facility from the database.
func (d *MyCareHubDb) DeleteFacility(ctx context.Context, id string) (bool, error) {
	return d.delete.DeleteFacility(ctx, id)
}
