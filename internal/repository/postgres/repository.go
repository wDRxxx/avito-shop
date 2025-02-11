package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wDRxxx/avito-shop/internal/repository"
	rm "github.com/wDRxxx/avito-shop/internal/repository/models"
)

type repo struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

const (
	usersTable        = "users"
	itemsTable        = "items"
	transactionsTable = "transactions"
	inventoryTable    = "inventory"
)

func NewPostgresRepo(db *pgxpool.Pool, timeout time.Duration) repository.Repository {
	return &repo{
		db:      db,
		timeout: timeout,
	}
}

func (r *repo) User(ctx context.Context, username string) (*rm.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"id",
		"username",
		"password",
		"balance",
		"created_at",
		"updated_at",
	).
		From(usersTable).
		Where(sq.Eq{"username": username}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	var user rm.User
	err = row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Balance,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repo) InsertUser(ctx context.Context, username string, password string) (*rm.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Insert(usersTable).Columns(
		"username",
		"password",
	).
		Values(username, password).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var id int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &rm.User{
		ID:        id,
		Username:  username,
		Password:  password,
		Balance:   0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (r *repo) Item(ctx context.Context, title string) (*rm.Item, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	builder := sq.Select(
		"id",
		"title",
		"price",
		"created_at",
		"updated_at",
	).From(itemsTable).
		Where(sq.Eq{"title": title}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRow(ctx, sql, args...)
	var item rm.Item
	err = row.Scan(
		&item.ID,
		&item.Title,
		&item.Price,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *repo) BuyItem(ctx context.Context, userID int, item *rm.Item) (err error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
			return
		}

		err = tx.Commit(ctx)
	}()

	updateBuilder := sq.Update(usersTable).
		Set("balance", sq.Expr("balance - ?", item.Price)).
		Where(sq.Eq{"id": userID}).
		PlaceholderFormat(sq.Dollar)

	sql, args, err := updateBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	insertBuilder := sq.Insert(inventoryTable).
		Columns("user_id", "item_id", "quantity").
		Values(userID, item.ID, 1).
		Suffix(`
	       ON CONFLICT (user_id, item_id)
	       DO UPDATE SET quantity = inventory.quantity + EXCLUDED.quantity
	   `).
		PlaceholderFormat(sq.Dollar)

	sql, args, err = insertBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}
