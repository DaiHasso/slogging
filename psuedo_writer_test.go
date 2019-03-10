/* #nosec G404 */
package logging

import (
    "math/rand"
    "strconv"
    "strings"
    "testing"

    gm "github.com/onsi/gomega"
)

func TestPsuedoWriter(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    psuedoWriter := NewPseudoWriter(INFO, newLogger)

    _, err = psuedoWriter.Write([]byte("Foobar"))
    g.Expect(err).ToNot(gm.HaveOccurred())

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `{"log_level":"INFO","message":"Foobar","timestamp":\d+}`,
    ))
}

func TestPsuedoWriterDebug(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
        WithLogLevel(DEBUG),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    psuedoWriter := NewPseudoWriter(DEBUG, newLogger)

    _, err = psuedoWriter.Write([]byte("Foobar"))
    g.Expect(err).ToNot(gm.HaveOccurred())

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `{"log_level":"DEBUG","message":"Foobar","timestamp":\d+}`,
    ))
}

func TestPsuedoWriterWarn(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    psuedoWriter := NewPseudoWriter(WARN, newLogger)

    _, err = psuedoWriter.Write([]byte("Foobar"))
    g.Expect(err).ToNot(gm.HaveOccurred())

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `{"log_level":"WARN","message":"Foobar","timestamp":\d+}`,
    ))
}

func TestPsuedoWriterError(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    var builder strings.Builder
    newLogger, err := NewLogger(
        "test" + strconv.Itoa(rand.Int()),
        WithLogWriters(&builder),
        WithFormat(JSON),
    )
    g.Expect(err).ToNot(gm.HaveOccurred())
    defer newLogger.Close()

    psuedoWriter := NewPseudoWriter(ERROR, newLogger)

    _, err = psuedoWriter.Write([]byte("Foobar"))
    g.Expect(err).ToNot(gm.HaveOccurred())

    logResult := builder.String()
    t.Log("\n" + logResult)
    g.Expect(logResult).To(gm.MatchRegexp(
        `{"log_level":"ERROR","message":"Foobar","timestamp":\d+}`,
    ))
}
