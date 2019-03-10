# Slogging
A simple logging framework for go that strives to be familiar while still being
powerful.

## Table Of Contents
* [Basic Usage](#basic-usage)
* [Creating a new logger](#creating-a-new-logger)
* [Retrieving Loggers By Identifier](#retrieving-loggers-by-identifier)
* [Logging Extras](#logging-extras)
* [Default Extras](#default-extras)
  + [Global Default Extras](#global-default-extras)
  + [Advanced Usage](#advanced-usage)
* [Logging formats](#logging-formats)
  + [JSON Example](#json-example)
  + [Standard Example](#standard-example)
  + [Standard Extended Example](#standard-extended-example)
* [Compatibility](#compatibility)


## Basic Usage
Using the slogging package directly you can just call Info, Debug, Warn, or
Error to use the singleton root logger. Like so:
``` go
package main

import (
    "github.com/daihasso/slogging"
)

func main() {
    logging.Info("Hello world!")
}
```

This would result in something like following log:
``` json
{"log_level":"INFO","message":"Hello world!","timestamp":1552169765}
```

The default format for the **root** logger is `JSON`. This is chosen with the
idea of deployments to production in mind. You may want to change this to
something a little more readable for local testing. You can that like so:
``` go
package main

import (
    "github.com/daihasso/slogging"
)

func main() {
    logging.GetRootLogger().SetFormat(logging.Standard)
    logging.Info("Hello world!")
}
```

This would result in the much more straightforward log-line:
``` text
2019-03-09T14:59:50 INFO Hello world!
```

\* **NOTE**: Changing the root logger directly can be dangerous if you're using
multiple logger instances as new loggers are based on the root logger. Keep this
in mind when changing the **root** logger.


## Creating a new logger
If you want a specialized instance of a Logger you can get one like so:
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    myBuf := new(strings.Builder)
    newLogger := logging.NewLogger(
        "MyCustomLogger",
        logging.WithFormat(logging.StandardExtended),
        logging.WithWriters(os.Stdout, myBuf),
        logging.WithLogLevel(logging.ERROR),
    )
    
    newLogger.Error("Just kidding, no error!")
}
```

This would result in something like the following being outputted to stdout and
to myBuf:
``` text
timestamp           | log_level | message
2019-03-09T14:59:50 | ERROR     | Just kidding, no error!
```

## Retrieving Loggers By Identifier
Every logger has an identifier (accessable via `logger.Identifier()`) which is
entered into a global registry in the slogging framework. This means if you want
to retrieve a given logger somewhere else in your code.
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    myLogger, ok := logging.GetLogger(MyLogger)
    if myLogger == nil {
        panic("MyLogger didn't exist!!!")
    }
    
    myLogger.Info(
}

func init() {
    logging.NewLogger(
        "MyLogger",
        logging.WithFormat(logging.Standard),
        logging.WithLogLevel(logging.DEBUG),
    )
}
```

## Logging Extras
Sometimes you don't just want to log a message, you also want to log some extra
data. With slogging, that's relatively straightforward:
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    logging.GetRootLogger().SetFormat(logging.Standard)
    logging.Info("Started app.", logging.Extras{
        "app_name": "Logging Test App"
    })
}
```

This would result in a log like the following:
``` text
2019-03-09T14:59:50 INFO app_name="Logging Test App" Started app.
```

\* **NOTE**: Generally the provided keys are used as-is however; for 
**Standard** format variants `\r`, `\n` and whitespace is removed and replaced
with `_`.

Your value for your extras doesn't even have to be a string! The provided value
will be marshalled into the log format (For **Standard** variants fmt.Sprint is
used to format the value).
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

type MyStruct struct {
    Name string
}

func main() {
    logging.GetRootLogger().SetFormat(logging.Standard)
    
    myStruct := MyStruct{
        Name: "Structicus",
    }
    logging.Info("Started app.", logging.Extras{
        "app_name": "Logging Test App",
        "my_struct": myStruct,
    })
}
```

``` text
2019-03-09T15:42:20 INFO app_name="Logging Test App" test="{Structicus}" Started app.
```

## Default Extras
Sometimes you want all logs for a logger to have a set of default `Extras` that
they log along with your message. This is where default extras come in.
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    newLogger := logging.NewLogger(
        "MyLoggerWithDefaultExtras",
        logging.WithFormat(logging.Standard),
        logging.WithDefaultExtras(
            StaticExtras(Extras{
                "app_name": "MyApp",
            }),
        ),
    )
    newLogger.Info("Started app.")
}
```

This would result in the following log line.
``` text
2019-03-09T15:52:13 INFO app_name="MyApp" Started app.
```

### Global Default Extras
Default `Extras` can also be set on a global scale--meaning every single logger
will evaluate these before logging. This can be done like so:
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    logging.AddGlobalExtras(
        StaticExtras(Extras{
            "app_name": "MyGlobalApp",
        }),
    ),
    newLogger := logging.NewLogger(
        "MyLogger",
        logging.WithFormat(logging.Standard),
    )
    newLogger.Info("Started app.")
}
```

This would result in the following log line.
``` text
2019-03-09T15:52:13 INFO app_name="MyGlobalApp" Started app.
```

### Advanced Usage
There are actually two types of default `Extras`; `StaticExtras` and
`FunctionalExtras`. The former takes extras in the format you would normally for
a specific log line as discussed above and the latter takes in a function that
is evaluated at log-time. Let's see `FunctionlExtras` in action.
``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    logCount := 0
    newLogger := logging.NewLogger(
        "MyLoggerWithSharedExtras",
        logging.WithFormat(logging.Standard),
        logging.WithDefaultExtras(
            Functional(ExtrasFuncs{
                "total_log_statements": func() (interface{}, error) {
                    logCount++
                    return logCount, nil
                },
            }),
        ),
    )
    newLogger.Info("Started app.")
    newLogger.Info("Quitting app.")
}
```

This would result in the following logs:
``` text
2019-03-09T15:55:24 INFO total_log_statements="1" Started app.
2019-03-09T15:55:25 INFO total_log_statements="2" Quitting app.
```

It's really quite powerful when used properly.

## Logging formats
Three formats are currently supported:
+ JSON
+ Standard
+ Standard Extended

### JSON Example
``` json
{"log_level":"INFO","message":"Hello world!","extra_key":"extra_value","timestamp":1552172390}
```

### Standard Example
``` text
2019-03-09T14:59:50 INFO extra_key="extra_value" Hello world!
```

### Standard Extended Example
``` text
timestamp           | log_level | extra_key   | message
2019-03-09T14:59:50 | INFO      | extra_value | Hello World!
```

## Compatibility
You may find yourself needing to provide a Logger to a library that expects a 
`io.Writer` or a `*log.Logger`. For this slogging provides the `PseudoWriter`.

``` go
package main

import (
    "os"
    "strings"

    "github.com/daihasso/slogging"
)

func main() {
    pseudoWriter := logging.NewPseudoWriter(
        logging.ERROR, logging.GetRootLogger(),
    )

    server := http.Server{
        Addr: "localhost:8080",
        ErrorLog: log.New(pseudoWriter, "", 0),
    }
    
    server.ListenAndServe()
}
```

The `PseudoWriter` is an ultra simple wrapper which simply wraps your logger and
logs to the provided LogLevel when `Write` is called on it.
