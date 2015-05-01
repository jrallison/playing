postgres benchmarks (all on local macbook air with ssd)
-------------------------------------------------------

Simple inserting with 0 indexes
-------------------------------
Inserting a million people with 50 attributes and 200 segments
(transactional batching of 100 inserts per transaction)

time: 13:10 (~1265 inserts per second)
7.7 GB disk space (du -sh /usr/local/var/postgres)

Inserting a million people with 20 attributes and 50 segments
(transactional batching of 100 inserts per transaction)

time: 4:38 (~3597 inserts per second)
2.2 GB disk space (du -sh /usr/local/var/postgres)

Simple inserting with 0 indexes & 1000 schemas
----------------------------------------------
Inserting a million people with 20 attributes and 50 segments
(transactional batching of 100 inserts per transaction)

time: 4:44 (~3521 inserts per second)
2.5 GB disk space (du -sh /usr/local/var/postgres)

Inserting with GIN index on attributes
--------------------------------------
Inserting a million people with 50 attributes and 200 segments
(transactional batching of 100 inserts per transaction)

time: 43:20 (~384 inserts per second)
9.7 GB disk space (du -sh /usr/local/var/postgres)

Inserting a million people with 20 attributes and 50 segments
(transactional batching of 100 inserts per transaction)

time: 14:10 (~1176 inserts per second)
~3.2 GB disk space (du -sh /usr/local/var/postgres)

Adding segment membership with GIN index on memberships
-------------------------------------------------------
Adding key to memberships hstore to 3 million people
(transactional batching of 100 updates per transaction)

time: 1:18:00 (~641 inserts per second)

Querying examples:
------------------

COUNTS ARE SLOW:

test=# explain analyze select id, attributes -> 'attr0' from s0.people where (attributes ? 'attr19' AND attributes ? 'attr18') OR (attributes ? 'attr17') OR (attributes ? 'attr16') order by attributes -> 'attr0';
                                                                         QUERY PLAN
-------------------------------------------------------------------------------------------------------------------------------------------------------------
 Sort  (cost=18652.94..18657.94 rows=2000 width=568) (actual time=47553.187..49160.574 rows=1000000 loops=1)
   Sort Key: ((attributes -> 'attr0'::text))
   Sort Method: external merge  Disk: 28680kB
   ->  Bitmap Heap Scan on people  (cost=11056.51..18543.29 rows=2000 width=568) (actual time=958.637..7126.340 rows=1000000 loops=1)
         Recheck Cond: (((attributes ? 'attr19'::text) AND (attributes ? 'attr18'::text)) OR (attributes ? 'attr17'::text) OR (attributes ? 'attr16'::text))
         Heap Blocks: exact=118155 lossy=131845
         ->  BitmapOr  (cost=11056.51..11056.51 rows=2001 width=0) (actual time=907.024..907.024 rows=0 loops=1)
               ->  Bitmap Index Scan on attrs_index  (cost=0.00..3696.01 rows=1 width=0) (actual time=460.311..460.311 rows=1000000 loops=1)
                     Index Cond: ((attributes ? 'attr19'::text) AND (attributes ? 'attr18'::text))
               ->  Bitmap Index Scan on attrs_index  (cost=0.00..3679.50 rows=1000 width=0) (actual time=216.265..216.265 rows=1000000 loops=1)
                     Index Cond: (attributes ? 'attr17'::text)
               ->  Bitmap Index Scan on attrs_index  (cost=0.00..3679.50 rows=1000 width=0) (actual time=230.447..230.447 rows=1000000 loops=1)
                     Index Cond: (attributes ? 'attr16'::text)
 Planning time: 0.166 ms
 Execution time: 49242.350 ms
(15 rows)

LIMITED SCANS ARE FAST:

test=# explain analyze select id, attributes -> 'attr0' from s0.people where (attributes ? 'attr19' AND attributes ? 'attr18') OR (attributes ? 'attr17') OR (attributes ? 'attr16') order by attributes -> 'attr0' limit 20;
                                                                      QUERY PLAN
-------------------------------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=0.42..6000.98 rows=20 width=568) (actual time=1.408..3.284 rows=20 loops=1)
   ->  Index Scan using attrs0_index on people  (cost=0.42..600056.42 rows=2000 width=568) (actual time=1.407..3.280 rows=20 loops=1)
         Filter: (((attributes ? 'attr19'::text) AND (attributes ? 'attr18'::text)) OR (attributes ? 'attr17'::text) OR (attributes ? 'attr16'::text))
 Planning time: 0.622 ms
 Execution time: 3.328 ms
(5 rows)

test=# select id, attributes -> 'attr0' from s0.people where (attributes ? 'attr19' AND attributes ? 'attr18') OR (attributes ? 'attr17') OR (attributes ? 'attr16') order by attributes -> 'attr0' limit 20;
   id   |    ?column?
--------+-----------------
      1 | valuee081
 100001 | valuee100000143
 100002 | valuee10000119
  10001 | valuee10000126
 100003 | valuee100002124
 100004 | valuee10000388
 100005 | valuee1000048
 100006 | valuee10000565
 100007 | valuee10000682
 100008 | valuee10000767
 100009 | valuee10000835
 100010 | valuee10000937
 100011 | valuee100010180
 100012 | valuee100011159
  10002 | valuee10001183
 100013 | valuee100012119
 100014 | valuee100013133
 100015 | valuee100014113
   1001 | valuee1000149
 100016 | valuee10001528
(20 rows)

LIMITED SCAN ON EQUALITY IS FAST (may be a different value in your generated dataset):

test=# explain analyze select id, attributes -> 'attr10' from s0.people where attributes @> 'attr10 => valuee094' limit 20;
                                                          QUERY PLAN
------------------------------------------------------------------------------------------------------------------------------
 Limit  (cost=59.75..136.18 rows=20 width=568) (actual time=1.487..1.489 rows=1 loops=1)
   ->  Bitmap Heap Scan on people  (cost=59.75..3881.30 rows=1000 width=568) (actual time=1.485..1.487 rows=1 loops=1)
         Recheck Cond: (attributes @> '"attr10"=>"valuee094"'::hstore)
         Heap Blocks: exact=1
         ->  Bitmap Index Scan on attrs_index  (cost=0.00..59.50 rows=1000 width=0) (actual time=1.395..1.395 rows=1 loops=1)
               Index Cond: (attributes @> '"attr10"=>"valuee094"'::hstore)
 Planning time: 24.400 ms
 Execution time: 1.674 ms
(8 rows)

test=# select id, attributes -> 'attr10' from s0.people where attributes @> 'attr10 => valuee094' limit 20;
 id | ?column?
----+-----------
  1 | valuee094
