## Go redis-lite

Go redis-lite is a lightweight and efficient database system designed to store key-value pairs. 
It provides six main commands: SET, GET, DELETE, EXPIRE, TTL AND KEYS allowing you to store, retrieve, delete, expire, ttl, keys command from the database. 


## Features

- **SET**: Store a key-value pair in the database, with an optional time-to-live (TTL) value.
  * Example Command : ```SET <key> <value> <ttl_in_secs>```
- **GET**: Retrieve the value associated with a given key from the database.
  * Example Command : ```GET <key>```
- **DELETE**: Removes a key-value pair(s) from the database.
  * Example Command : ```DEL <key1> <key2> ...```
- **EXPIRE**: Set TTL for given key from the database.
  * Example Command : ```EXPIRE <key> <ttl_in_secs>```
- **TTL**: Retrieve expired time left for a given key from the database.
  * Example Command : ```TTL <key>```
- **KEYS**: Returns all keys matching pattern.
  * Example Command : ```KEYS <pattern>```

## Getting Started

To get started with Go redis-lite, follow these steps:

1. Clone the repository: `https://github.com/evilzone/redis-lite.git`
2. Start database server: `make run`


## Usage

Once the database server is up and running, you can interact with it using any tcp client like telnet or nc. Here are some examples:

```shell
%  nc localhost 8000 
SET key1 val1
OK
SET key2 val2
OK
KEYS key*
1) key1
2) key2

GET key1
val1
GET j
(nil)
TTL key1
-1
TTL ke
-2
TTL key1 100
-1
DEL key1
1
GET key1
(nil)
GET key2
val2
EXPIRE key2 100
1
EXPIRE ku 100
0
TTL key2
83
GET key2
val2
TTL key2
-2
GET key2
(nil)
SET key1 val1
OK
GET key1
val1
TTL key1
-1
KEYS key*
1) key1

SET key2 val2
OK
KEYS key*
1) key2
2) key1

TTL key2
-1
TTL key1
-1
EXPIRE key1 1000
1
EXPIRE m 100
0
TTL m
-2
TTL key1
987
```
