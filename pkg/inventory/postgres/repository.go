package postgres

import (
	"context"
	"fmt"

	"github.com/Tra-Dew/inventory-write/pkg/inventory"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type repositoryPostgres struct {
	pool *pgxpool.Pool
}

// NewRepository ...
func NewRepository(pool *pgxpool.Pool) inventory.Repository {
	return &repositoryPostgres{
		pool: pool,
	}
}

// InsertBulk ...
func (r *repositoryPostgres) InsertBulk(ctx context.Context, items []*inventory.Item) error {

	batch := &pgx.Batch{}

	stmt := `
		insert into
		items(id, owner_id, name, status, description, total_quantity, locked_quantity, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	for _, i := range items {
		batch.Queue(
			stmt,
			i.ID,
			i.OwnerID,
			i.Name,
			i.Status,
			i.Description,
			i.TotalQuantity,
			i.LockedQuantity,
			i.CreatedAt,
			i.UpdatedAt,
		)
	}

	tx, _ := r.pool.Begin(ctx)
	res := tx.SendBatch(ctx, batch)

	defer res.Close()

	_, err := res.Exec()

	return err
}

// UpdateBulk ...
func (r *repositoryPostgres) UpdateBulk(ctx context.Context, userID *string, items []*inventory.Item) error {

	batch := &pgx.Batch{}

	filter := "id = $1"

	if userID != nil {
		filter = filter + " and and owner_id = $2"
	}

	stmt := fmt.Sprintf(`
		update items
		set
			name = $3,
			status = $4,
			description = $5,
			total_quantity = $6,
			locked_quantity = $7,
			created_at = $8,
			updated_at = $9
		where
			%s
	`, filter)

	for _, i := range items {
		batch.Queue(
			stmt,
			i.ID,
			i.OwnerID,
			i.Name,
			i.Status,
			i.Description,
			i.TotalQuantity,
			i.LockedQuantity,
			i.CreatedAt,
			i.UpdatedAt,
		)
	}

	tx, _ := r.pool.Begin(ctx)
	res := tx.SendBatch(ctx, batch)

	defer res.Close()

	_, err := res.Exec()

	return err
}

// DeleteBulk ...
func (r *repositoryPostgres) DeleteBulk(ctx context.Context, userID string, ids []string) error {

	return nil
}

// Get ...
func (r *repositoryPostgres) Get(ctx context.Context, userID string, ids []string) ([]*inventory.Item, error) {

	return nil, nil
}

// GetByStatus ...
func (r *repositoryPostgres) GetByStatus(ctx context.Context, status inventory.ItemStatus) ([]*inventory.Item, error) {

	return nil, nil
}
