
## Start Redis
```
docker run --rm -p 6379:6379 -d redis
```

## Start golinks
```
./golinks --backend redis --redis-addr localhost:6379 --DB 0
```