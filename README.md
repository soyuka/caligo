# Caligo 🧿

🧿 [Caligo meaning](https://fr.wiktionary.org/wiki/caligo#la).

## Stack

- Go http api
- [etcd](https://etcd.io/) database
- [nanoid](https://github.com/matoous/go-nanoid) for id generation

Go HTTP and etcd should scale well if you have money.

## How

- /?google.com

## Configuration

- `ETCD_URL` (default=localhost:2379) coma-separated list of etcd hosts
- `ETCD_DIAL_TIMEOUT` (default=5s) dial timeout for etcd
- `CALIGO_HOSTNAME` (default=localhost:8080) hostname
- `CALIGO_ID_LENGTH` (default=12) nanoid length (see [collision calculator](https://zelark.github.io/nano-id-cc/))
- `CALIGO_ID_ALPHABET` (default=0123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNOPQRSTUVWXYZ) nanoid alphabet)
- `CALIGO_PORT` (default=5376)
