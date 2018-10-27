package catch

// funcCall represents a function
// call.
type funcCall struct {
	err      interface{}
	finished bool
}

// Panic catches a panic prone func
// and returns
func Panic(panicProne func()) (bool, interface{}) {
	errChan := make(chan funcCall)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				errChan <- funcCall{
					err:      err,
					finished: false,
				}
			}
		}()
		panicProne()
		errChan <- funcCall{
			err:      nil,
			finished: true,
		}
	}()
	ret := <-errChan
	return ret.finished, ret.err
}
