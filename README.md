## Prerequisite
```bash
go install github.com/qiuzhanghua/autotag@latest
```

# Build and Run
```bash
make && BUNDEBUG=error ./qraft
```

# Redis
on my m1x, redis cluster nodes=6
```bash
wrk -t12 -c200 -d30s http://localhost:9999/rdb
```
1. 不设置PoolSize，33k qps
2. PoolSize = 1, 38k qps
3. PoolSize = 2, 63k qps
4. PoolSize = 3, 72k qps, with few socket errors
5. PoolSize = 4, 76k qps, with few socket errors

所以PoolSize建议设置3或者2(生产环境中也需要验证测试一下).


# MySQL
on my m1x, 8.0.30
```bash
wrk -t12 -c200 -d30s http://localhost:9999/db
```

```go
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sql.SetMaxOpenConns(maxOpenConns)
	sql.SetMaxIdleConns(maxOpenConns)
	sql.SetConnMaxIdleTime(2 * time.Second)
	sql.SetConnMaxLifetime(30 * time.Second)

	...
	
var ver string
err = db.NewRaw("SELECT version() as v").Scan(ctx, &ver)

```
约为32k qps.

