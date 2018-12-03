package catch

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// adding some work so benchmark is accurate
const benchmarkSampleSize = 400 * 1000

func repeatForBenchmark(f func(), multithreadedContext bool) {
	for i := 0; i < benchmarkSampleSize; i++ {
		if multithreadedContext {
			go f()
		} else {
			f()
		}
	}
}

func getPureGoBenchFunc(panic bool) func() {
	wrap := func(f func()) func() {
		return func() {
			defer func() {
				recover()
			}()
			f()
		}
	}
	if panic {
		return wrap(panicProneFunc)
	}
	return wrap(panicLessFunc)

}
func getCatchBenchFunc(panic bool) func() {
	if panic {
		return func() {
			Panic(panicProneFunc)
		}
	}
	return func() {
		Panic(panicLessFunc)
	}
}

func runBenchmarks(panic bool, multithreadedContext bool, b *testing.B) {
	b.Run("pure go", func(b *testing.B) {
		repeatForBenchmark(getPureGoBenchFunc(panic), multithreadedContext)
	})
	b.Run("catch", func(b *testing.B) {
		repeatForBenchmark(getCatchBenchFunc(panic), multithreadedContext)
	})
}

func BenchmarkWithPanic(b *testing.B) {
	panics := true
	b.Run("single goroutine", func(b *testing.B) {
		runBenchmarks(panics, false, b)
	})
	b.Run("multiple goroutines", func(b *testing.B) {
		runBenchmarks(panics, true, b)
	})
}

func BenchmarkWithoutPanic(b *testing.B) {
	panics := false
	b.Run("single goroutine", func(b *testing.B) {
		runBenchmarks(panics, false, b)
	})
	b.Run("multiple goroutines", func(b *testing.B) {
		runBenchmarks(panics, true, b)
	})
}

var expectedError = "42"

func panicFunc(value interface{}) func() {
	return func() {
		panic(value)
	}
}
func panicProneFunc() {
	panic(expectedError)
}
func panicLessFunc() {}

func compareAsString(t *testing.T, expected interface{}, actual interface{}) {
	expectedString, actualString := fmt.Sprintf("%s", expected), fmt.Sprintf("%s", actual)
	assert.Equal(t, expectedString, actualString)
}
func TestPanic(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		panicked, err := Panic(
			panicFunc(expectedError),
		)
		assert.True(t, panicked)
		assert.Equal(t, expectedError, err)
	})
	t.Run("no panic", func(t *testing.T) {
		panicked, err := Panic(panicLessFunc)
		assert.False(t, panicked)
		assert.Nil(t, err)
	})
}

func TestInterface(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		err := Interface(panicFunc(expectedError))
		compareAsString(t, expectedError, err)
	})
	t.Run("type conservation", func(t *testing.T) {
		err := Interface(panicFunc(42))
		assert.Equal(t, 42, err)
	})
	t.Run("panic with nil ", func(t *testing.T) {
		err := Interface(panicFunc(nil))
		assert.NotNil(t, err)
	})
	t.Run("no panic", func(t *testing.T) {
		err := Interface(panicLessFunc)
		assert.Nil(t, err)
	})
}

func TestValuesToInterface(t *testing.T) {
	t.Run("on nil", func(t *testing.T) {
		interfaces := valuesToInterfaces(nil)
		assert.Nil(t, interfaces)
	})
	t.Run("convert values", func(t *testing.T) {
		values := []reflect.Value{
			reflect.ValueOf("hello"),
			reflect.ValueOf("world"),
		}
		interfaces := valuesToInterfaces(values)
		assert.Equal(t, len(interfaces), 2)
		for i, actualValue := range interfaces {
			compareAsString(t, values[i].Interface(), actualValue)
		}
	})
}

func TestError(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		err := Error(panicProneFunc)
		compareAsString(t, expectedError, err)
	})
	t.Run("nil error", func(t *testing.T) {
		err := Error(panicLessFunc)
		assert.Nil(t, err)
	})
}

func TestCancall(t *testing.T) {
	panicless := func() (int, string) {
		return 0, ""
	}
	panicky := func() (string, int) {
		panic(42)
	}
	t.Run("positive", func(t *testing.T) {
		assert.True(t, CanCall(panicless))
	})
	t.Run("negative", func(t *testing.T) {
		assert.False(t, CanCall(panicky))
	})
}

func TestSanitizeWithProcedure(t *testing.T) {
	t.Run("with panic", func(t *testing.T) {
		sanitizedFunction := SanitizeFunc(panicFunc(expectedError))
		values, err := sanitizedFunction()
		compareAsString(t, expectedError, err)
		assert.Nil(t, values)
	})
	t.Run("without panic", func(t *testing.T) {
		sanitizedFunction := SanitizeFunc(func() (string, string) {
			return "hello", "world"
		})
		values, err := sanitizedFunction()
		assert.Nil(t, err)
		compareAsString(t, "[hello world]", values)
	})
}

func TestSanitizeWithFunc(t *testing.T) {
	t.Run("", func(t *testing.T) {
		sanitizedFunction := SanitizeFunc(func(bool) {
			panic(expectedError)
		})
		values, err := sanitizedFunction(false)
		compareAsString(t, expectedError, err)
		assert.Nil(t, values)
	})
	t.Run("", func(t *testing.T) {
		sanitizedFunction := SanitizeFunc(func(s1 string, s2 string) (string, string) {
			return s2, s1
		})
		values, err := sanitizedFunction("world", "hello")
		assert.Nil(t, err)
		compareAsString(t, "[hello world]", values)
	})
}
