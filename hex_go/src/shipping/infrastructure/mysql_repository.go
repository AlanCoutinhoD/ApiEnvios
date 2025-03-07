package infrastructure

import (
	"context"
	"database/sql"
	"demo/src/shipping/domain"
	"fmt"
	"time"
)

type MySQLRepository struct {
	db *sql.DB
}

func NewMySQLRepository(db *sql.DB) *MySQLRepository {
	return &MySQLRepository{
		db: db,
	}
}

func (r *MySQLRepository) Save(shipping *domain.Shipping) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `INSERT INTO shipping (idUser, idProduct, quantity, created_at) VALUES (?, ?, ?, ?)`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("error preparando statement: %v", err)
	}
	defer stmt.Close()

	now := time.Now()
	result, err := stmt.ExecContext(ctx, shipping.IdUser, shipping.IdProduct, shipping.Quantity, now)
	if err != nil {
		return fmt.Errorf("error guardando shipping: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error obteniendo ID: %v", err)
	}

	shipping.ID = id
	shipping.CreatedAt = now
	return nil
}

func (r *MySQLRepository) GetByID(id int64) (*domain.Shipping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, idUser, idProduct, quantity, created_at FROM shipping WHERE id = ?`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparando statement: %v", err)
	}
	defer stmt.Close()

	shipping := &domain.Shipping{}
	err = stmt.QueryRowContext(ctx, id).Scan(
		&shipping.ID,
		&shipping.IdUser,
		&shipping.IdProduct,
		&shipping.Quantity,
		&shipping.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error consultando shipping: %v", err)
	}

	return shipping, nil
}

func (r *MySQLRepository) GetByUserID(userID int64) ([]*domain.Shipping, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `SELECT id, idUser, idProduct, quantity, created_at FROM shipping WHERE idUser = ? ORDER BY created_at DESC`
	
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error preparando statement: %v", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error consultando shippings: %v", err)
	}
	defer rows.Close()

	var shippings []*domain.Shipping
	for rows.Next() {
		shipping := &domain.Shipping{}
		err := rows.Scan(
			&shipping.ID,
			&shipping.IdUser,
			&shipping.IdProduct,
			&shipping.Quantity,
			&shipping.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error escaneando shipping: %v", err)
		}
		shippings = append(shippings, shipping)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando resultados: %v", err)
	}

	return shippings, nil
}
