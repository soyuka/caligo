# Caligo 🧿

🧿 [Caligo meaning](https://fr.wiktionary.org/wiki/caligo#la).

## Stack

- Go http api
- [boltdb](https://github.com/etcd-io/bbolt) database
- [nanoid](https://github.com/matoous/go-nanoid) for id generation

Go HTTP and boltdb should scale well if you have money.

## How

- /?google.com

## Configuration

- `CALIGO_DB` (default=bolt://data.bolt) db path
- `CALIGO_HOSTNAME` (default=localhost:5376) hostname
- `CALIGO_ID_LENGTH` (default=12) nanoid length (see [collision calculator](https://zelark.github.io/nano-id-cc/))
- `CALIGO_ID_ALPHABET` (default=0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ) nanoid alphabet)
- `CALIGO_PORT` (default=5376)

## Docker

```
docker pull soyuka/caligo
docker run -d --name caligo -p 5376:5376 -e CALIGO_DB=bolt:///bolt/caligo.bolt -e CALIGO_HOSTNAME=https://caligo.space -v "/home/soyuka/caligodb:/bolt" soyuka/caligo
```

## Kubernetes

Get redis password and create a config map.

```
kubectl get secret --namespace default caligo-redis -o jsonpath="{.data.redis-password}" | base64 --decode
kubectl create configmap caligo-config --from-literal caligo-db=redis://:Y3Jv8HO2Yc@localhost:6379/1 
```

## TODO

- allow custom url
- statistics click
