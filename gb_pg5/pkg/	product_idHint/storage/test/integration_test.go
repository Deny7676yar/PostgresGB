// +build integration_tests

package storage

import (
	"Deny7676yar/Go_level2/gb_pg5/pkg/	product_idHint/storage"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestIntegrationSearch(t *testing.T) {
	ctx := context.Background()
	dbpool := connect(ctx)
	defer dbpool.Close()

	tests := []struct {
		store *storage.PG
		ctx context.Context
		store_id int
		product_id int
		prepare func(*pgxpool.Pool)
		check func(*testing.T, []storage.GetQuantityByProduct_idHint, error)
	}{
		{
			store: storage.NewPG(dbpool),
			ctx: context.Background(),
			store_id: 1,
			product_id: 1,
			prepare: func(dbpool *pgxpool.Pool) {
				dbpool.Exec(context.Background(),`insert into quantity ...`)
			},
			check: func(t *testing.T, hints []storage.GetQuantityByProduct_idHint, err error) {
			require.NoError(t, err)
			require.NotEmpty(t, hints)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.product_id, func(t *testing.T) {
			tt.prepare(dbpool)
			hints, err := tt.store.GetQuantityByProduct_id(tt.ctx, tt.store_id, tt.product_id)
			tt.check(t, hints, err)
		})
	}
}

// Соединение с экземпляром Postgres
func connect(ctx context.Context) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	return dbpool
}