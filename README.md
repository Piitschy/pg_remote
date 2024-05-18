# pgrd

## Installation

### CLI
```bash
go install github.com/Piitschy/pgrd/cmd/pgrd@0.3.0
```
### Server side componente
```bash
docker pull piitschy/pgrd:0.3.0
```

## Server

Lad den Server einfach in deinen Docker-Stack neben einer Postgress-Istanz - über das Environment kannst du alle wichtigen Konfigurationen angeben.

Dabei muss ein Port auf 3000 exposed werden.

## CLI

Die CLI nutzt den offenen Port, um mit http-Methoden Daten zu transferieren.

### dump

```bash
pgrd --host <host> -p <post> -k <Key> dump -o <filename>
```

### local dump 
```bash
pgrd -u <db-user> --pw <db-password> --db <database> -p <post> ldump -o <filename>
```

### restore

```bash
pgrd --host <host> -p <post> -k <Key> restore -i <filename>
```

### local restore 
```bash
pgrd -u <db-user> --pw <db-password> --db <database> -p <post> lrestore -o <filename>
```
