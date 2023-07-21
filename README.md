# postgress-dump-tool
 
## Server

Lad den Server einfach in deinen Docker-Stack neben einer Postgress-Istanz - über das Environment kannst du alle wichtigen Konfigurationen angeben.

Dabei muss ein Port auf 3000 exposed werden.

## CLI

Die CLI nutzt den offenen Port, um mit http-Methoden Daten zu transferieren.

### dump

```bash
pg_remote -H <host> -p <post> dump > dump.sql
```
oder
```bash
pg_remote -H <host> -p <post> dump -f dump.sql
```

### restore

Beim Restoring wird immer auch ein Dump erzeugt.

```bash
pg_remote -H <host> -p <post> restore < dump.sql
```
oder
```bash
pg_remote -H <host> -p <post> restore -f dump.sql
```

### list

Gibt eine Liste aller gespecherten Dumps

```bash
pg_remote -H <host> -p <post> list
```

### load

```bash
pg_remote -H <host> -p <post> load [dumpname] > dump.sql
```
oder
```bash
pg_remote -H <host> -p <post> load -f dump.sql [dumpname]
```
