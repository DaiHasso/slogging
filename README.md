# slogging
Simple logging framework for go.

## Usage

Using the slogging package directly you can just call Info, Debug, Warn, or
Error to use the singleton default logger. Like so:
``` go
slogging.Info("HTTP request made.").Send()
```

### Getting a new logger
If you want a special instance of a logger you can get one like so:

``` go
newLogger := slogging.GetNewLogger(
    "MyCustomLogger",
    slogging.JSON,
    slogging.Stdout,
    []slogging.LogLevel{slogging.ERROR},
```

### Logging formats
Two formats are currently supported:
+ JSON
+ ELF

**JSON example:**
``` json
{"message":"HTTP request made","log_level":"INFO","timestamp":1530763651}
```

**ELF example:**
``` text
#Version: 1.0
#Date: 2018-07-06 19:02:56.0311886 +0000 UTC m=+132925.628801801
#Fields: level | message
INFO | HTTP request made.
```

### Logging a simple message
**Code:**
``` go
slogging.Info("HTTP request made.").Send()
```

**Output:**
``` json
{"message":"HTTP request made","log_level":"INFO","timestamp":1530763651}
```

### Logging a message with extras
**Code:**
``` go
slogging.Info("HTTP request made.").
    With("path", r.URL.Path).
    And("requestor", r.RemoteAddr).
    And("method", r.Method).
    Send()
```

**Output:**
``` json
{"message":"HTTP request made","log_level":"INFO","timestamp":1530763651,"path":
"/test","requestor":"127.0.0.1","method":"GET"}
```
