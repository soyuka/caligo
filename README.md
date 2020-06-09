# Caligo 🧿

🧿 [Caligo meaning](https://fr.wiktionary.org/wiki/caligo#la).

## Stack

- Go http api
- [boltdb](https://github.com/etcd-io/bbolt) database
- [nanoid](https://github.com/matoous/go-nanoid) for id generation

Go HTTP and etcd should scale well if you have money.

## How

- /?google.com

## Configuration

- `CALIGO_DB_PATH` (default=data.bolt) db path
- `CALIGO_HOSTNAME` (default=localhost:8080) hostname
- `CALIGO_ID_LENGTH` (default=12) nanoid length (see [collision calculator](https://zelark.github.io/nano-id-cc/))
- `CALIGO_ID_ALPHABET` (default=0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ) nanoid alphabet)
- `CALIGO_PORT` (default=5376)

## Docker

```
docker pull soyuka/caligo
docker run -d --name caligo -p 5376:5376 -e CALIGO_DB_PATH=/bolt/caligo.bolt -e CALIGO_HOSTNAME=https://caligo.space -v "/home/soyuka/caligodb:/bolt" soyuka/caligo
```
