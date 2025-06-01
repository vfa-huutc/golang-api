package logger_test

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"

	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
)

func TestLogger(t *testing.T) {

	t.Run("Test Init", func(t *testing.T) {
		logger.Init()
		assert.NotNil(t, logrus.StandardLogger().Formatter)
	})

	t.Run("Info level logs", func(t *testing.T) {
		t.Run("Info", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.InfoLevel)
			defer hook.Reset()

			logger.Info("hello world")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.InfoLevel, entry.Level)
			assert.Equal(t, "hello world", entry.Message)
		})

		t.Run("Infof", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.InfoLevel)
			defer hook.Reset()

			logger.Infof("hello %s", "world")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.InfoLevel, entry.Level)
			assert.Equal(t, "hello world", entry.Message)
		})
	})

	t.Run("Debug level logs", func(t *testing.T) {
		t.Run("Debug", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.DebugLevel)

			defer hook.Reset()

			logger.Debug("debug msg")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.DebugLevel, entry.Level)
			assert.Equal(t, "debug msg", entry.Message)
		})

		t.Run("Debugf", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.DebugLevel)
			defer hook.Reset()

			logger.Debugf("debug %s", "msg")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.DebugLevel, entry.Level)
			assert.Equal(t, "debug msg", entry.Message)
		})
	})

	t.Run("Error level logs", func(t *testing.T) {
		t.Run("Error", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.ErrorLevel)
			defer hook.Reset()

			logger.Error("error: not found")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.ErrorLevel, entry.Level)
			assert.Equal(t, "error: not found", entry.Message)
		})

		t.Run("Errorf", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.ErrorLevel)
			defer hook.Reset()

			logger.Errorf("error: %s", "not found")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.ErrorLevel, entry.Level)
			assert.Equal(t, "error: not found", entry.Message)
		})
	})

	t.Run("Warning level logs", func(t *testing.T) {
		t.Run("Warn", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.WarnLevel)
			defer hook.Reset()

			logger.Warn("this is a warning")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.WarnLevel, entry.Level)
			assert.Equal(t, "this is a warning", entry.Message)
		})

		t.Run("Warnf", func(t *testing.T) {
			hook := test.NewGlobal()
			logrus.SetLevel(logrus.WarnLevel)

			defer hook.Reset()

			logger.Warnf("this is a %s", "warning")

			assert.Len(t, hook.Entries, 1)
			entry := hook.LastEntry()
			assert.Equal(t, logrus.WarnLevel, entry.Level)
			assert.Equal(t, "this is a warning", entry.Message)
		})
	})
}
