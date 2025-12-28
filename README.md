# example docker compose for valkey cluster connection

Turns out dynamically specifying `--cluster-announce-hostname` was the key to connect valkey cluster with hostname like `valkey-cluster:6379`.

```sh
docker compose up --build

docker compose down --volumes --remove-orphans
```
