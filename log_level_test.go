package logging

import (
    "testing"

    gm "github.com/onsi/gomega"
)

func TestLogLevelFromString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logLevel := LogLevelFromString("debug")

    g.Expect(logLevel).To(gm.BeIdenticalTo(DEBUG))
}

func TestLogsEnabledForString(t *testing.T) {
    g := gm.NewGomegaWithT(t)

    logLevel, err := GetLogLevelsForString("info")
    g.Expect(err).ToNot(gm.HaveOccurred())

    logsForInfo, err := logsEnabledFromLevel(INFO)
    g.Expect(err).ToNot(gm.HaveOccurred())

    g.Expect(logLevel).To(gm.Equal(logsForInfo))
}
