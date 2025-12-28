# example docker compose for valkey cluster connection

```sh
docker compose --profile valkey-cluster up --build

docker compose --profile valkey-cluster down --volumes --remove-orphans
```
