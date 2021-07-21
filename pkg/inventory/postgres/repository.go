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
	_, err := res.Exec()
	res.Close()

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// UpdateBulk ...
func (r *repositoryPostgres) UpdateBulk(ctx context.Context, userID *string, items []*inventory.Item) error {

	batch := &pgx.Batch{}

	filter := "id = $1"

	if userID != nil {
		filter = filter + " and owner_id = $2"
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
			userID,
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
	_, err := res.Exec()
	res.Close()

	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteBulk ...
func (r *repositoryPostgres) DeleteBulk(ctx context.Context, userID string, ids []string) error {

	stmt := `
		delete from items
		where
			id = any($1) and owner_id = $2
	`
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, stmt, ids, userID); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *repositoryPostgres) Get(ctx context.Context, userID string, ids []string) ([]*inventory.Item, error) {
	var items []*inventory.Item

	rows, err := r.pool.Query(ctx, `select * from items where id = any($1)`, ids)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := new(inventory.Item)
		rows.Scan(
			&item.ID,
			&item.OwnerID,
			&item.Name,
			&item.Status,
			&item.Description,
			&item.TotalQuantity,
			&item.LockedQuantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		items = append(items, item)
	}

	return items, nil
}

// GetByStatus ...
func (r *repositoryPostgres) GetByStatus(ctx context.Context, status inventory.ItemStatus) ([]*inventory.Item, error) {
	var items []*inventory.Item

	rows, err := r.pool.Query(ctx, `select * from items where status = $1`, status)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := new(inventory.Item)
		rows.Scan(
			&item.ID,
			&item.OwnerID,
			&item.Name,
			&item.Status,
			&item.Description,
			&item.TotalQuantity,
			&item.LockedQuantity,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		items = append(items, item)
	}

	return items, nil
}
