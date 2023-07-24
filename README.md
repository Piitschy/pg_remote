# postgress-dump-tool
 
## Server

Lad den Server einfach in deinen Docker-Stack neben einer Postgress-Istanz - über das Environment kannst du alle wichtigen Konfigurationen angeben.

Dabei muss ein Port auf 3000 exposed werden.

## CLI

Die CLI nutzt den offenen Port, um mit http-Methoden Daten zu transferieren.

### dump

```bash
pg_remote -H <host> -p <post> -k <Key> dump [dump.tar]
```

### restore

Beim Restoring wird immer auch ein Dump erzeugt.

```bash
pg_remote -H <host> -p <post> -k <Key> restore [dump.tar]
```
