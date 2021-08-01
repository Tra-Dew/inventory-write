package postgres

import (
	"context"
	"fmt"

	"github.com/d-leme/tradew-inventory-write/pkg/inventory"
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

	sqlItems := `
		insert into
		items(id, owner_id, name, status, description, total_quantity, created_at, updated_at)
		values($1, $2, $3, $4, $5, $6, $7, $8)
	`

	sqlLocks := `
		insert into
		item_locks(item_id, locked_by, quantity)
		values($1, $2, $3)
	`

	for _, i := range items {
		batch.Queue(
			sqlItems,
			i.ID,
			i.OwnerID,
			i.Name,
			i.Status,
			i.Description,
			i.TotalQuantity,
			i.CreatedAt,
			i.UpdatedAt,
		)

		for _, l := range i.Locks {
			batch.Queue(
				sqlLocks,
				i.ID,
				l.LockedBy,
				l.Quantity,
			)
		}
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	res := tx.SendBatch(ctx, batch)

	if _, err := res.Exec(); err != nil {
		return err
	}

	res.Close()

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// UpdateBulk ...
func (r *repositoryPostgres) UpdateBulk(ctx context.Context, items []*inventory.Item) error {

	batch := &pgx.Batch{}

	sqlItems := `
		update items
		set
			name = $1,
			status = $2,
			description = $3,
			total_quantity = $4,
			created_at = $5,
			updated_at = $6
		where
			id = $7
	`
	sqlDeleteLocks := `
		delete from item_locks
		where
			item_id = $1
	`

	sqlInsertLock := `
		insert into item_locks (item_id, locked_by, quantity)
		values($1, $2, $3)
	`

	for _, i := range items {
		batch.Queue(sqlItems,
			i.Name, i.Status, i.Description,
			i.TotalQuantity, i.CreatedAt, i.UpdatedAt, i.ID,
		)

		batch.Queue(sqlDeleteLocks, i.ID)
		for _, l := range i.Locks {
			batch.Queue(sqlInsertLock, i.ID, l.LockedBy, l.Quantity)
		}
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	res := tx.SendBatch(ctx, batch)

	if _, err := res.Exec(); err != nil {
		return err
	}

	res.Close()

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// DeleteBulk ...
func (r *repositoryPostgres) DeleteBulk(ctx context.Context, ids []string) error {

	sqlDeleteLocks := `
		delete from item_locks
		where
			item_id = any($1)
	`

	sqlDeleteItems := `
		delete from items
		where
			id = any($1)
	`

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sqlDeleteLocks, ids); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sqlDeleteItems, ids); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *repositoryPostgres) Get(ctx context.Context, userID *string, ids []string) ([]*inventory.Item, error) {

	filter := "i.id = any($1)"
	args := []interface{}{ids}

	if userID != nil {
		filter = filter + " and owner_id = $2"
		args = append(args, *userID)
	}

	sql := fmt.Sprintf(`
		select * from items i
			left join item_locks l on i.id = l.item_id
		where
			%s
	`, filter)

	return r.getItems(ctx, sql, args...)
}

// GetByStatus ...
func (r *repositoryPostgres) GetByStatus(ctx context.Context, status inventory.ItemStatus) ([]*inventory.Item, error) {

	sql := `select * from items where status = $1`

	return r.getItems(ctx, sql, status)
}

func (r *repositoryPostgres) getItems(ctx context.Context, sql string, args ...interface{}) ([]*inventory.Item, error) {
	itemMap := map[string]*inventory.Item{}

	rows, err := r.pool.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		item := new(inventory.Item)
		var itemID, lockedBy *string
		var quantity *int64

		err := rows.Scan(
			&item.ID, &item.OwnerID, &item.Name, &item.Status,
			&item.Description, &item.TotalQuantity,
			&item.CreatedAt, &item.UpdatedAt,

			&itemID, &lockedBy, &quantity,
		)

		if err != nil {
			return nil, err
		}

		if i, exist := itemMap[item.ID]; exist {
			i.Locks = append(i.Locks, &inventory.ItemLock{LockedBy: *lockedBy, Quantity: inventory.ItemQuantity(*quantity)})
		} else {
			if itemID != nil {
				item.Locks = append(item.Locks, &inventory.ItemLock{LockedBy: *lockedBy, Quantity: inventory.ItemQuantity(*quantity)})
			}
			itemMap[item.ID] = item
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var items []*inventory.Item
	for _, item := range itemMap {
		items = append(items, item)
	}

	return items, nil
}
