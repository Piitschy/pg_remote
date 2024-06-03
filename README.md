# pgrd

## Installation

### CLI
```bash
go install github.com/Piitschy/pgrd/cmd/pgrd@latest
```
### Server side component
```bash
docker pull piitschy/pgrd:0.3.1
```

## Server

Simply add the server to your Docker stack alongside a Postgres instance - you can specify all important configurations through the environment.

A port must be exposed on 3000.

## CLI

The CLI uses the open port to transfer data with HTTP methods.

### dump

```bash
pgrd --host <host> -p <port> -k <key> dump -o <filename>
```

### local dump 
```bash
pgrd -u <db-user> --pw <db-password> --db <database> -p <port> ldump -o <filename>
```

### restore

```bash
pgrd --host <host> -p <port> -k <key> restore -i <filename>
```

### local restore 
```bash
pgrd -u <db-user> --pw <db-password> --db <database> -p <port> lrestore -i <filename>
```
