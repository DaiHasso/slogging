/* #nosec G404 */
package logging

import (
    "math/rand"
    "strconv"
    "strings"
    "testing"
    "errors"

    gm "github.com/onsi/gomega"
)

func TestLoggerBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerIdentifier(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    identifier := "test" + strconv.Itoa(rand.Int())
    newLogger, err := NewLogger(identifier)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    g.Expect(newLogger.Identifier()).To(gm.BeIdenticalTo(identifier))
}

func TestLoggerDuplicateName(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    newLogger, err := NewLogger("test_dupe")
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger2, err := NewLogger("test_dupe")
    g.Expect(err).To(gm.MatchError(
        "Can't create new logger with identifier 'test_dupe'; identifier " +
            "already exists",
    ))

    newLogger.Close()

    newLogger2, err = NewLogger("test_dupe")
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger2.Close()
}

func TestLoggerError(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Error("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"ERROR","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerWarn(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Warn("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"WARN","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerStandardExtended(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(StandardExtended),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `^timestamp\s* \| log_level\s* \| message\s*\n` +
            `[^|]+ \| INFO\s* \| Foo\s*$`,
    ))
}

func TestLoggerStandard(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(Standard),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `[^\s]+ INFO Foo`,
    ))
}

func TestLoggerStandardExtendedExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(StandardExtended),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo", Extras{
        "test": "baz",
    })

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `^timestamp\s* \| log_level\s* \| message\s* \| test\s*\n` +
            `[^|]+ \| INFO\s* \| Foo\s* \| baz\s*$`,
    ))
}

func TestLoggerStandardExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(Standard),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo", Extras{
        "test": "baz",
    })

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `[^\s]+ INFO test="baz" Foo`,
    ))
}

func TestLoggerStandardExtrasStruct(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(Standard),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    testStruct := struct{
        Foo string
        bar int
    }{
        Foo: "baz",
        bar: 6,
    }
    newLogger.Info("Foo", Extras{
        "test": testStruct,
    })

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `[^\s]+ INFO test="{baz 6}" Foo`,
    ))
}

func TestLoggerJSONExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo", Extras{
        "test_log": "baz",
    })

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","test_log":"baz"` +
            `,"timestamp":\d+}`,
    ))
}

func TestLoggerJSONExtraFunc(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo", Extra("test_log", "baz"))

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","test_log":"baz"` +
            `,"timestamp":\d+}`,
    ))
}

func TestLoggerJSONExtraFuncAndStruct(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info(
        "Foo", Extra("test_log", "baz"), Extras{
            "test_2": "baz2",
        },
    )

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","test_2":"baz2",` +
            `"test_log":"baz","timestamp":\d+}`,
    ))
}

func TestLoggerLogLevel(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithLogLevel(INFO),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Debug("Foo")

    g.Expect(builder.String()).To(gm.BeEmpty())
}

func TestLoggerException(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Exception(errors.New("test err!"), "Oh noes, an error.")

    stringResult := builder.String()

    t.Log(stringResult)

    g.Expect(stringResult).To(gm.MatchRegexp(
        `{"error":"test err![^"]+","log_level":"ERROR",` +
            `"message":"Oh noes, an error.","timestamp":\d+}`,
    ))
}

func TestLoggerGlobalExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    AddGlobalExtras(StaticExtras(Extras{
        "foo": "bar",
    }))
    defer SetGlobalExtras()

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    val := "bar"

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(StaticExtras(Extras{
            "foo": val,
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerDefaultExtrasFuncs(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    val := "bar"

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(FunctionalExtras(ExtrasFuncs{
            "foo": func() (interface{}, error) {
                return val, nil
            },
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))

    val = "baz"

    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"baz","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerAddWriter(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()


    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))

    g.Expect(builder2.String()).To(gm.BeEmpty())

    newLogger.AddWriters(&builder2)

    builder.Reset()

    newLogger.Info("Bar")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
    g.Expect(builder2.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
}

func TestLoggerSetFormat(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder strings.Builder
    )
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(StandardExtended),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()


    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `^timestamp\s* \| log_level\s* \| message\s*\n` +
            `[^|]+ \| INFO\s* \| Foo\s*$`,
    ))

    builder.Reset()

    newLogger.SetFormat(JSON)

    newLogger.Info("Bar")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
}

func TestLoggerAddDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))

    builder.Reset()

    newLogger.AddDefaultExtras(StaticExtras(Extras{
        "bar": "baz",
    }))

    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"bar":"baz","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerSetDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(StaticExtras(Extras{
            "foo": "bar",
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))

    builder.Reset()

    newLogger.SetDefaultExtras(StaticExtras(Extras{
        "test": "setextras",
    }))

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","test":"setextras",` +
            `"timestamp":\d+}`,
    ))
}


func TestLoggerRemoveWriter(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()


    newLogger.AddWriters(&builder2)

    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
    g.Expect(builder2.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))

    builder.Reset()
    builder2.Reset()

    newLogger.RemoveWriter(&builder2)

    newLogger.Info("Bar")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
    g.Expect(builder2.String()).To(gm.BeEmpty())
}

func TestLoggerSetWriters(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()



    newLogger.Info("Foo")
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
    g.Expect(builder2.String()).To(gm.BeEmpty())

    builder.Reset()
    builder2.Reset()

    newLogger.SetWriters(&builder2)

    newLogger.Info("Bar")
    g.Expect(builder.String()).To(gm.BeEmpty())
    g.Expect(builder2.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
}

func TestLoggerDefaultExtrasFailure(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(FunctionalExtras(ExtrasFuncs{
            "foo": func() (interface{}, error) {
                return nil, errors.New("Test extras error.")
            },
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"error":"Test extras error.[^"]+","log_level":"ERROR",` +
            `"message":"Error while running logger instance extras.",` +
            `"timestamp":\d+}\n` +
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestLoggerGlobalExtrasFailure(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    AddGlobalExtras(FunctionalExtras(ExtrasFuncs{
        "foo": func() (interface{}, error) {
            return nil, errors.New("Test extras error.")
        },
    }))
    defer SetGlobalExtras()

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"error":"Test extras error.[^"]+","log_level":"ERROR",` +
            `"message":"Error while running global logger extras.",` +
            `"timestamp":\d+}\n` +
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestCloneLogger(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(StaticExtras(Extras{
            "foo": 5,
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    newLogger2, err := CloneLogger(
        "test" + strconv.Itoa(rand.Int()),
        newLogger,
        WithDefaultExtras(StaticExtras(Extras{
            "bar": 10,
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger2.Close()

    newLogger2.Info("Foo")

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"bar":10,"foo":5,"log_level":"INFO","message":"Foo",` +
            `"timestamp":\d+}`,
    ))

    builder.Reset()

    newLogger.Info("Foo")

    logResult = builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":5,"log_level":"INFO","message":"Foo",` +
            `"timestamp":\d+}`,
    ))
}
