package catch

import (
	"fmt"
	"reflect"
)

// trivial implementation of a functionCall
type funcCall struct {
	err      interface{}
	panicked bool
}

// encapsulate the panic prone function and notify the channel.
func encapsulate(panicProne func(), errChan chan funcCall) {
	defer func() {
		errChan <- funcCall{
			err:      recover(),
			panicked: true,
		}
	}()
	panicProne()
	errChan <- funcCall{
		err:      nil,
		panicked: false,
	}
}

// Panic catches a panic prone func.
func Panic(panicProne func()) (panicked bool, recoverReturn interface{}) {
	errChan := make(chan funcCall)
	go encapsulate(panicProne, errChan)
	ret := <-errChan
	return ret.panicked, ret.err
}

// Interface returns the recovered value as an interface.
func Interface(panicProne func()) interface{} {
	panicked, err := Panic(panicProne)
	if err != nil {
		return err
	}
	if panicked {
		return fmt.Errorf("panic called with a nil error")
	}
	return nil
}

// Error returns the recovered value as error.
func Error(panicProne func()) error {
	recoveredValue := Interface(panicProne)
	if recoveredValue != nil {
		return fmt.Errorf("%s", recoveredValue)
	}
	return nil
}

// valuesToInterfaces convert reflect values to interfaces.
func valuesToInterfaces(values []reflect.Value) []interface{} {
	if values == nil {
		return nil
	}
	interfaces := make([]interface{}, len(values))
	for i, value := range values {
		interfaces[i] = value.Interface()
	}
	return interfaces
}

// Sanitize converts a panic prone function to a function that returns an error.
func SanitizeFunc(panicProneFunc interface{}) func(args ...interface{}) (returnedValues []interface{}, err interface{}) {
	callbackValue := reflect.ValueOf(panicProneFunc)
	return func(args ...interface{}) ([]interface{}, interface{}) {
		in := make([]reflect.Value, 0)
		for _, arg := range args {
			argValue := reflect.ValueOf(arg)
			in = append(in, argValue)
		}
		var retValues []interface{}
		err := Interface(func() {
			returnedValues := callbackValue.Call(in)
			retValues = valuesToInterfaces(returnedValues)
		})
		return retValues, err
	}
}
