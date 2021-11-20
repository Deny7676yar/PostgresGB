package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	url2 "net/url"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Name string
	Email string
)

type NameSearchHint struct {
	Name Name
	Email Email
}

// search ищет всех сотрудников, email которых начинается с prefix.
// Из функции возвращается список EmailSearchHint, отсортированный по Email.
// Размер возвращаемого списка ограничен значением limit.
func search(ctx context.Context, dbpool *pgxpool.Pool, prefix string, limit int)([]NameSearchHint, error) {
	const sql = `select email, phone from analog_prod where email like $1 order by email asc limit $2;`
	pattern := prefix + "%"
	rows, err := dbpool.Query(ctx, sql, pattern, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query data: %w", err)
	}
	// Вызов Close нужен, чтобы вернуть соединение в пул
	defer rows.Close()

	// В слайс hints будут собраны все строки, полученные из базы
	var hints []NameSearchHint

	// rows.Next() итерируется по всем строкам, полученным из базы.
	for rows.Next() {
		var hint NameSearchHint

		// Scan записывает значения столбцов в свойства структуры hint
		err = rows.Scan(&hint.Email, &hint.Name)
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

type AttackResults struct {
	Duration time.Duration
	Threads int
	QueriesPerformed uint64
}

func attack(ctx context.Context, duration time.Duration, threads int, dbpool *pgxpool.Pool) AttackResults {
	var queries uint64

	attacker := func (stopAt time.Time) {
		for {
			_, err := search(ctx, dbpool, "jjamary2", 5)
			if err != nil { log.Fatal(err) }

			atomic.AddUint64(&queries, 1)

			if time.Now().After(stopAt) {
				return
			}
		}
	}

	var wg sync.WaitGroup
	wg.Add(threads)

	startAt := time.Now()
	stopAt := startAt.Add(duration)

	for i := 0; i < threads; i++ {
		go func() {
			attacker(stopAt)
			wg.Done()
		}()
	}
	wg.Wait()
	return AttackResults{
		Duration:         time.Now().Sub(startAt),
		Threads:          threads,
		QueriesPerformed: queries,
	}
}

func main() {
	ctx := context.Background()
	url :=
		fmt.Sprintf("%s://%s:%s/%s?user=%s;password=%s",
			"postgresql",
			"localhost",
			"54333",
			"guest",
			url2.QueryEscape("guest"),
			url2.QueryEscape("guest")
		)

		//"postgresql://localhost:54333/guest?user=guest;password=guest"

	cfg, err := pgxpool.ParseConfig(url)
	if err != nil { log.Fatal(err) }

	cfg.MaxConns = 8
	cfg.MinConns = 4

	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil { log.Fatal(err) }
	defer dbpool.Close()

	duration := time.Duration(5 * time.Second)
	threads := 1000
	fmt.Println("start attack")
	res := attack(ctx, duration, threads, dbpool)

	fmt.Println("duration:", res.Duration)
	fmt.Println("threads:", res.Threads)
	fmt.Println("queries:", res.QueriesPerformed)
	qps := res.QueriesPerformed / uint64(res.Duration.Seconds())
	fmt.Println("QPS:", qps)
}
