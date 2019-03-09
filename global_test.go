package logging

import (
    "log"
    "strings"
    "testing"

    gm "github.com/onsi/gomega"
)

func TestNewLoggerJSON(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger, err := NewChainLogger(
        "test",
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(logger).To(gm.BeAssignableToTypeOf(&JSONLogger{}))
}

func TestNewLoggerELF(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logger, err := NewChainLogger(
        "test",
        WithFormat(ELF),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(logger).To(gm.BeAssignableToTypeOf(&ELFLogger{}))
}

func TestGlobalDebug(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    err := GetDefaultLogger().SetLogLevel("DEBUG")
    g.Expect(err).ToNot(gm.HaveOccurred())
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    Debug("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"DEBUG","timestamp":\d+}`,
    ))
}

func TestGlobalWarn(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    err := GetDefaultLogger().SetLogLevel("WARN")
    g.Expect(err).ToNot(gm.HaveOccurred())
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    Warn("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"WARN","timestamp":\d+}`,
    ))
}

func TestGlobalError(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    err := GetDefaultLogger().SetLogLevel("ERROR")
    g.Expect(err).ToNot(gm.HaveOccurred())
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    Error("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"ERROR","timestamp":\d+}`,
    ))
}

func TestGlobalInfo(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    err := GetDefaultLogger().SetLogLevel("INFO")
    g.Expect(err).ToNot(gm.HaveOccurred())
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    Info("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"INFO","timestamp":\d+}`,
    ))
}

func TestSetDefaultLoggerLogLevel(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    Debug("Foo").Send()
    g.Expect(builder.String()).To(gm.BeEmpty())

    err := SetDefaultLoggerLogLevel("DEBUG")
    g.Expect(err).ToNot(gm.HaveOccurred())

    Debug("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"DEBUG","timestamp":\d+}`,
    ))
}

func TestAddGlobalExtrasBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    err := SetDefaultLoggerLogLevel("DEBUG")
    g.Expect(err).ToNot(gm.HaveOccurred())

    AddGlobalExtras(StaticExtras(Extras{
        "test": "bar",
    }))
    defer func() {
        SetGlobalExtras()
    }()

    Debug("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"DEBUG","timestamp":\d+,"test":"bar"}`,
    ))
}

func TestAddGlobalExtrasFunc(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    err := SetDefaultLoggerLogLevel("DEBUG")
    g.Expect(err).ToNot(gm.HaveOccurred())

    testVal := "baz"

    AddGlobalExtras(FunctionalExtras(
        ExtrasFuncs{
            "test": func() (interface{}, error) {
                return testVal, nil
            },
        },
    ))
    defer func() {
        SetGlobalExtras()
    }()

    Debug("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"DEBUG","timestamp":\d+,"test":"baz"}`,
    ))

    testVal = "baz2"

    Debug("Foo").Send()

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"message":"Foo","log_level":"DEBUG","timestamp":\d+,"test":"baz2"}`,
    ))
}

func TestSetGetGlobalExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    defaultInternalLogger := GetDefaultLogger().GetInternalLogger()
    GetDefaultLogger().SetInternalLogger(log.New(&builder, "", 0))
    defer func() {
        GetDefaultLogger().SetInternalLogger(defaultInternalLogger)
    }()

    err := SetDefaultLoggerLogLevel("DEBUG")
    g.Expect(err).ToNot(gm.HaveOccurred())

    extra := StaticExtras(Extras{
        "test": "bar",
    })
    SetGlobalExtras(extra)
    defer func() {
        SetGlobalExtras()
    }()
    globalExtras := GetGlobalExtras()
    t.Log(extra)
    g.Expect(globalExtras).To(gm.HaveLen(1))
}
