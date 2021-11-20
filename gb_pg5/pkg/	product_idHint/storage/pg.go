package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PG struct {
	dbpool *pgxpool.Pool
}

func NewPG(dbpool *pgxpool.Pool) *PG {
	return &PG{dbpool}
}

type (
	Store_id int
	Product_id int
	Quantity int
)

type GetQuantityByProduct_idHint struct {
	Store_id Store_id
	Product_id Product_id
	Quantity Quantity
}

func (s *PG) GetQuantityByProduct_id(ctx context.Context, store_id int, product_id int)([]GetQuantityByProduct_idHint, error) {
	const q = `SELECT quantity FROM quantity where store_id = $1 AND product_id = $2;`
	rows, err := s.dbpool.Query(ctx, q, store_id, product_id)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []GetQuantityByProduct_idHint

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint GetQuantityByProduct_idHint

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Store_id, &hint.Product_id, &hint.Quantity)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		hints = append(hints, hint)
	}

	// Проверка, что во время выборки данных не происходило ошибок
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to read response: %w", rows.Err())
	}
	return hints, nil
}