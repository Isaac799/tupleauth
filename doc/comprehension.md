# Notes

below are notes (several quotes too) from [article](https://authzed.com/zanzibar) to help me understand and set goals

the engine

- high-concurrency graph engine that lives on top of a database
  - graphs created on the fly, following relationships of shards
  - stateless
- designed for traversal (not just search)
  - leaf always either a user (deadend) or a userset
  - edges is relation

a tuple

- External consistency and snapshot reads with staleness bounded by zookie prevent:
  - Neglecting ACL update order
  - Misapplying old ACL to new content
- Before clients can store... namespace configuration specifies its relations as well as its storage parameters
- rules
  - IO
    - in: obj ID
    - out: userset expression tree
  - `this` (default): all users including indirect referenced by usersets
  - `computed_userset`: new userset like: viewer relation to refer to the editor userset on the same object
    - allow inferring an object’s owner ID from the object ID prefix, which reduces space requirements for clients that manage many private objects
  - `tuple_to_userset`: enables: look up the parent folder of the document and inherit its viewers
    - represent object hierarchy with only one relation tuple per hop

## zookie

- the lower bound for a read snapshot
- the upper bound for config usage
- on write get latest zookie from the write
- on read/check get either same or default stale (whatever used)
- no zookie = default

zookie is an opaque byte sequence encoding a globally meaningful timestamp

- bookmarks a requests perception of world at request
- lets server know its stale
- prevents out of order execution
- not a lifetime

reflects:
- an ACL write
- a client content version
- a read snapshot

use an opaque cookie instead of the actual timestamp 
- discourage our clients from choosing arbitrary timestamps
- allow future extensions

may be a b64 encoding of
- timestamp
- sequence (breaks tie of timestamps)
- clusterID (breaks tie of sequence, since sequence generation can be per regional-database)

clients can provide zookies to ensure a minimum snapshot timestamp for request evaluation. When a zookie is not provided, the server uses a default staleness 

### server default staleness


- chosen to ensure that all transactions are evaluated at a timestamp that is as recent as possible without impacting latency
- (performance optimization) strategy to calculate more accurate default over time
  - identify out-of-zone reads
  - z-test to see if increasing the default helps

## api

### read

request specifies

- one or multiple tuplesets
  - set can include a single tuple key, or all tuples with a given object ID or userset in a namespace
- and an optional zookie

look up w/ tuplesets

- a specific membership entry
- read all entries in an ACL or group
- look up all groups with a given user as a direct member

look up w/ zookie

- a read snapshot no earlier than a previous write (if the zookie from the write response is given in the read request)
- at the same snapshot as a previous read (if the zookie from the earlier read response is given in the subsequent request)

(else) doesn’t contain a zookie choose a reasonably recent snapshot

### write

read modify write (rmw) at database level

1. Read all relation tuples of an object

- include per object lock

2. writes, along with a touch on the lock tuple
3. fail, go to 1

- Optimistic Concurrency Control (OCC) with a 'Lock Tuple'.
- One 'leader' no matter shard distribution that locks an object
- Saves to cluster for shard range based on shard map (so shard is routing for read and write)

### watch

- request specifies:
  - one or more namespaces
  - a zookie representing the time to start watching
- response contains:
  - all tuple modification events in ascending timestamp order
  - heartbeat zookie (can use to resume watching)

### check

- request specifies:
  - userset
  - user (some token)
  - zookie corresponding to the desired object version

Like reads, a check is always evaluated at a consistent snapshot no earlier than the given zookie

- content-change check (to authorize application content modifications)
  - request: does not carry a zookie and is evaluated at the latest snapshot
  - response:
    - zookie for clients to store
    - object contents and use for subsequent checks of the content version

### expand

- request specifies:
  - object#relation
- response contains:
  - effective userset
    - userset tree whose leaf nodes are user IDs or usersets
    - (concept) intermediate nodes represent union, intersection, or exclusion operators
  - optional zookie

enables clients to:
- reason about the complete set of object access for users and groups
- build efficient search indices for access-controlled content

## architecture

- database 
  - store relation tuples for each client namespace
    - pk: shard ID, object ID, relation, user, commit timestamp
      - shard ID: computed from both object ID and user 
        - shards used for performance optimizations
        - break up giant groups (high concurrency for check) and keeps permissions together (low check latency)
        - acts as routing (tells the system where in the world the data lives)
          - there is a source of truth map saying shards 0-100 on server 1, 101-200 on server 2...
    - ordering of primary keys 
      - allows us to look up all relation tuples for a given object ID or (object ID, relation) pair
  - hold all namespace configurations
    - 2 tables
      - 1: ID, configs
      - 2: changelog of config updates and is keyed by commit timestamps
  - one changelog database shared across all namespaces
    - pk: changelog shard ID, timestamp, unique update ID
    - history of tuple updates for the Watch API
- cluster that respond to Check, Read, Expand, and Write requests
  - fans out the work (as necessary) compute intermediate results
- watchservers are a specialized server type that respond to Watch requests
  - stream of namespace changes to clients
  - subscribes to shard leader
  - the 'realized' graph - some sliding window to read all zookie on startup
  - tail the changelog
  - denormalization
    - merges on the fly based on zookie
    - saved to separate table (an index)
    - watch server manages the index, deleting records based on reads against shards
- background jobs
  - backup (produce dumps of the relation tuples in each namespace at a known snapshot)
  - garbage-collect (tuple versions older than a threshold configured per namespace)
  - optimize large and deeply nested sets (transformations on that data, such as denormalization)

cluster config sync

- problem: cluster configs out of sync
- solution: ? provider tracks who has access to what finds latest common ground

Check Evaluation

- evaluate all leaf nodes of the boolean expression tree concurrently. 
- When the outcome of one node determines the result of a subtree, evaluation of other nodes in the subtree is cancelled

## to consider

- namespace config management and usage (zookie can represent schema? too)
- write must be committed to both the tuple storage and the changelog shard in a single transaction (rollback handling)
- recursive cycles
- collapsing requests ('thundering herd')

## observations

- can startup 
  - instantly to solve acl needs
  - only delay is watcher init, and if a zookie says a server is stale and that server waiting for watcher to update it
- the engine that deals with the zookie, and the engine will send a zookie request to the correct shard leader and that leader will always have the most updated info for a shard (stall if waiting for inter cluster sync). So the watcher is watching the shard leaders
