package logging

import (
    "strings"
    "testing"
    "errors"

    gm "github.com/onsi/gomega"
)

func TestInstantLoggerBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestInstantLoggerError(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Error("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"ERROR","message":"Foo","timestamp":\d+}`,
    ))
}

func TestInstantLoggerWarn(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Warn("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"WARN","message":"Foo","timestamp":\d+}`,
    ))
}

func TestInstantLoggerStandard(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo")

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `^timestamp\s* \| log_level\s* \| message\s*\n` +
            `[^|]+ \| INFO\s* \| Foo\s*$`,
    ))
}

func TestInstantLoggerStandardLong(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(StandardLong),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo")

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `[^\s]+ INFO Foo`,
    ))
}

func TestInstantLoggerStandardExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

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

func TestInstantLoggerStandardLongExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(StandardLong),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo", Extras{
        "test": "baz",
    })

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `[^\s]+ INFO test="baz" Foo`,
    ))
}

func TestInstantLoggerJSONExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

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

func TestInstantLoggerJSONExtraFunc(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo", Extra("test_log", "baz"))

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","test_log":"baz"` +
            `,"timestamp":\d+}`,
    ))
}

func TestInstantLoggerJSONExtraFuncAndStruct(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

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

func TestInstantLoggerLogLevel(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithLogLevel(INFO),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Debug("Foo")

    g.Expect(builder.String()).To(gm.BeEmpty())
}

func TestInstantLoggerException(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Exception(errors.New("test err!"), "Oh noes, an error.")

    stringResult := builder.String()

    t.Log(stringResult)

    g.Expect(stringResult).To(gm.MatchRegexp(
        `{"error":"test err![^"]+","log_level":"ERROR",` +
            `"message":"Oh noes, an error.","timestamp":\d+}`,
    ))
}

func TestInstantLoggerGlobalExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    AddGlobalExtras(StaticExtras(Extras{
        "foo": "bar",
    }))
    defer SetGlobalExtras()

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestInstantLoggerDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    val := "bar"

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(StaticExtras(Extras{
            "foo": val,
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    newLogger.Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"foo":"bar","log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestInstantLoggerDefaultExtrasFuncs(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    val := "bar"

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(FunctionalExtras(ExtrasFuncs{
            "foo": func() (interface{}, error) {
                return val, nil
            },
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

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

func TestInstantLoggerAddWriter(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())


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

func TestInstantLoggerSetFormat(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder strings.Builder
    )
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())


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

func TestInstantLoggerAddDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
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

func TestInstantLoggerSetDefaultExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithDefaultExtras(StaticExtras(Extras{
            "foo": "bar",
        })),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

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


func TestInstantLoggerRemoveWriter(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())


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

func TestInstantLoggerSetWriters(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var (
        builder,
        builder2 strings.Builder
    )
    newLogger, err := NewLogger(
        "test",
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())



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
