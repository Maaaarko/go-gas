# Instructions

## Generate crt and key files
```bash
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## Run the server
```bash
make start
```


