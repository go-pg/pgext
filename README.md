# go-pg extensions

## Faster JSON encoding by segmentio

```go
import (
    "github.com/go-pg/pg/v10/pgjson"
    "github.com/go-pg/pgext"
)

func init() {
    pgjson.SetProvider(pgext.SegmentJSONProvider{})
}
```

## Tracing using OpenTelemetryHook

For more details see [documentation](https://pg.uptrace.dev/tracing/):

```go
db := pg.Connect(&pg.Options{...})
db.AddQueryHook(&pgext.OpenTelemetryHook{})
```

## Print failed queries using DebugHook

```go
db := pg.Connect(&pg.Options{...})

if debug {
    db.AddQueryHook(&pgext.DebugHook{
        //Verbose: true,
    })
}
```
