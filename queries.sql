SELECT count(*) from analog_prod;

[2021-10-28 01:10:52] 1 row retrieved starting from 1 in 148 ms (execution: 4 ms, fetching: 144 ms)

explain select email, phone from analog_prod where email = 'esiemons6@hao123.com';

[2021-10-28 01:12:17] 2 rows retrieved starting from 1 in 48 ms (execution: 4 ms, fetching: 44 ms)

[2021-10-28 01:15:23] completed in 171 ms

explain analyze select email, phone from analog_prod where email like 'esiemons%' order by email asc;

[2021-10-28 01:16:06] 9 rows retrieved starting from 1 in 38 ms (execution: 22 ms, fetching: 16 ms)

Sort  (cost=8.31..8.31 rows=1 width=41) (actual time=0.018..0.018 rows=0 loops=1)
  Sort Key: email
  Sort Method: quicksort  Memory: 25kB
  ->  Index Only Scan using employee_email_idx on analog_prod  (cost=0.28..8.30 rows=1 width=41) (actual time=0.005..0.005 rows=0 loops=1)
        Index Cond: ((email ~>=~ 'alexan'::text) AND (email ~<~ 'alexao'::text))
        Filter: ((email)::text ~~ 'alexan%'::text)
        Heap Fetches: 0
Planning Time: 16.772 ms
Execution Time: 0.031 ms

