/* #nosec G404 */
package logging

import (
    "math/rand"
    "os"
    "strconv"
    "strings"
    "testing"
    "errors"

    gm "github.com/onsi/gomega"
)

func TestGetLogger(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    identifier := "test" + strconv.Itoa(rand.Int())
    loggerFromGet := GetLogger(identifier)
    g.Expect(loggerFromGet).To(gm.BeNil())

    newLogger, err := NewLogger(identifier)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    loggerFromGet = GetLogger(identifier)
    g.Expect(loggerFromGet).To(gm.BeIdenticalTo(newLogger))
}

func TestSetRootLogger(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var (
        builder,
        builder2 strings.Builder
    )
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        err := SetRootLoggerExisting(initialRootLoggerName)
        g.Expect(err).ToNot(gm.HaveOccurred())
    }()

    identifier := "test" + strconv.Itoa(rand.Int())
    newLogger, err := NewLogger(
        identifier,
        WithLogWriters(&builder2),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
    g.Expect(builder2.String()).To(gm.BeEmpty())

    builder.Reset()
    builder2.Reset()

    err = SetRootLogger(identifier, newLogger)

    Info("Bar")

    g.Expect(err).ToNot(gm.HaveOccurred())
    g.Expect(builder.String()).To(gm.BeEmpty())
    g.Expect(builder2.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Bar","timestamp":\d+}`,
    ))
}

func TestGlobalException(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var (
        builder strings.Builder
    )
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    defer func() {
        rootLogger.SetWriters(os.Stdout)
    }()

    Exception(errors.New("Test err"), "Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"error":"Test err[^"]+","log_level":"ERROR",` +
            `"message":"Foo","timestamp":\d+}`,
    ))
}

func TestGlobalDebug(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    Debug("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foo","timestamp":\d+}`,
    ))
}

func TestGlobalWarn(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(WARN)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    Warn("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"WARN","message":"Foo","timestamp":\d+}`,
    ))
}

func TestGlobalError(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(ERROR)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    Error("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"ERROR","message":"Foo","timestamp":\d+}`,
    ))
}

func TestGlobalInfo(t *testing.T) {
    g := gm.NewGomegaWithT(t)
    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(INFO)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    Info("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foo","timestamp":\d+}`,
    ))
}

func TestSetDefaultLoggerLogLevel(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    Debug("Foo")
    g.Expect(builder.String()).To(gm.BeEmpty())

    err := rootLogger.SetLogLevel(DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())

    Debug("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foo","timestamp":\d+}`,
    ))
}

func TestAddGlobalExtrasBasic(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

    AddGlobalExtras(StaticExtras(Extras{
        "test": "bar",
    }))
    defer func() {
        SetGlobalExtras()
    }()

    Debug("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foo","test":"bar","timestamp":\d+}`,
    ))
}

func TestAddGlobalExtrasFunc(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

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

    Debug("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foo","test":"baz","timestamp":\d+}`,
    ))

    builder.Reset()

    testVal = "baz2"

    Debug("Foo")

    g.Expect(builder.String()).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foo","test":"baz2","timestamp":\d+}`,
    ))
}

func TestSetGetGlobalExtras(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    rootLogger := GetRootLogger()
    rootLogger.SetWriters(&builder)
    err := rootLogger.SetLogLevel(DEBUG)
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer func() {
        rootLogger.SetWriters(os.Stdout)
        _ = rootLogger.SetLogLevel(INFO)
    }()

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
