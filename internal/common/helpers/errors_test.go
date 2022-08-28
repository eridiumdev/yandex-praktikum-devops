package helpers

import (
	"strconv"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
)

type sourceMock string

func (s sourceMock) LogSource() string {
	return string(s)
}

func TestErrorAfterWrapping(t *testing.T) {
	var (
		err    = errors.New("error")
		source = sourceMock("")
	)

	t.Run("single wrap", func(t *testing.T) {
		wrapped := WrapErr(err, source, "wrap")
		assert.EqualValues(t, "wrap: error", wrapped.Error())
	})
	t.Run("double wrap", func(t *testing.T) {
		wrapped := WrapErr(WrapErr(err, source, "wrap"), source, "wrap")
		assert.EqualValues(t, "wrap: wrap: error", wrapped.Error())
	})
	t.Run("10-times wrap", func(t *testing.T) {
		wrapped := err
		for i := 0; i < 10; i++ {
			wrapped = WrapErr(wrapped, source, "wrap")
		}
		assert.EqualValues(t, strings.Repeat("wrap: ", 10)+"error", wrapped.Error())
	})
}

func TestTraceAfterWrapping(t *testing.T) {
	err := errors.New("error")

	t.Run("single wrap", func(t *testing.T) {
		wrapped := WrapErr(err, sourceMock("1"), "")
		traceable, ok := wrapped.(logger.Traceable)
		require.True(t, ok)
		assert.EqualValues(t, "1", traceable.LogTrace())
	})
	t.Run("double wrap", func(t *testing.T) {
		wrapped := WrapErr(WrapErr(err, sourceMock("2"), ""), sourceMock("1"), "")
		traceable, ok := wrapped.(logger.Traceable)
		require.True(t, ok)
		assert.EqualValues(t, "1->2", traceable.LogTrace())
	})
	t.Run("10-times wrap", func(t *testing.T) {
		wrapped := err
		expected := ""
		for i := 0; i < 10; i++ {
			wrapped = WrapErr(wrapped, sourceMock(strconv.Itoa(9-i)), "")
			if i == 0 {
				expected += strconv.Itoa(i)
			} else {
				expected += "->" + strconv.Itoa(i)
			}
		}
		traceable, ok := wrapped.(logger.Traceable)
		require.True(t, ok)
		assert.EqualValues(t, expected, traceable.LogTrace())
	})
}
