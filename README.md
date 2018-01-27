## redis-executor
Setup a HA redis cluster by redis-executor.
## Build
```sh
make build-redis-executor
```
Built though by docker container.
## Run
```sh
redis-executor server -l ${DEBUG_LEVEL} \
  --configFile /etc/redis/redis-executor.conf
```
## Result
We can using `redis-executor` to lanuch a HA redis cluster, `redis-executor` will setup as a Appropriate roles:

- redis master
- redis slave
- redis sentinel (using environmental variable SENTINEL=true)

It is more robust than https://github.com/kubernetes/kubernetes/blob/master/examples/storage/redis/image/run.sh.
