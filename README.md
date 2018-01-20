# Sharded-Key-Value-Service
A key/value storage system that "shards" or partitions the keys over a set of replica groups.
- Characteristic:
  - Master-Slave Architecture 
  - Fault tolerant( Raft protocol for replication set)
  - Horizontal scalable
  
Sharded key/value store will have two main components. First, a set of "Replica Groups". 
Each replica group is responsible for a subset of the shards. A replica consists of a handful of servers that use Raft to 
replicate the group's shard. The second component is the "Shard Master". The shard master decides which replica group 
should serve each shard. RPC may be used for interaction between clients and servers, between different servers, 
and between different clients.

This architecture is patterned at a high level on a number of systems: BigTable, Spanner, Apache HBase and many others.

The above was the part of lab assignment for class ##MIT 6.824 - Distributed Systems - Spring 2015
https://pdos.csail.mit.edu/archive/6.824-2013/labs/lab-4.html
