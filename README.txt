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
