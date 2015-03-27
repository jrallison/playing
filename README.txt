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

Inserting with GIN index on attributes
--------------------------------------
Inserting a million people with 50 attributes and 200 segments
(transactional batching of 100 inserts per transaction)

time: ~40:00 (~384 inserts per second)
~10 GB disk space (du -sh /usr/local/var/postgres)

Querying:

Count number of people with existance of (attr49 AND attr48) OR attr47 OR attr46
--------------------------------------------------------------------------------
select count(*) from test.people where (attributes ? 'attr49' AND attributes ? 'attr48') OR (attributes ? 'attr47') OR (attributes ? 'attr46');
