package catch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCatch(t *testing.T) {
	t.Run("catches panic", func(t *testing.T) {
		expectedError := "Hello World !"
		finished, err := Panic(func() {
			panic(expectedError)
		})
		assert.False(t, finished)
		assert.Equal(t, expectedError, err)
	})
	t.Run("no false catch", func(t *testing.T) {
		finished, err := Panic(func() {
		})
		assert.True(t, finished)
		assert.Nil(t, err)
	})
}

func BenchmarkCatch(b *testing.B) {
	panicError := "Hello World !"
	emptyFunc, panicProneFunc := func() {}, func() { panic(panicError) }
	b.Run("with panic", func(b *testing.B) {
		b.Run("golang basic panic/recover", func(b *testing.B) {
			defer func() {
				recover()
			}()
			panicProneFunc()
		})
		b.Run("catch.panic way", func(b *testing.B) {
			Panic(panicProneFunc)
		})
	})
	b.Run("without panic", func(b *testing.B) {
		b.Run("golang basic panic/recover", func(b *testing.B) {
			defer func() {
				recover()
			}()
			emptyFunc()
		})
		b.Run("catch.panic way", func(b *testing.B) {
			Panic(emptyFunc)
		})
	})
}
