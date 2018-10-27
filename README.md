# Catch
 Catch is a simple way for inlining Golang recover
### Benchmarks show no difference between the two methods. The performance impact may be understated due to compiler optimization.
# Golang using Catch
``` golang
import (
	"github.com/potatomasterrace/catch"
)

func main() {
	finished,caughtError := catch.Panic(func(){
        // panic prone logic 
    })
    // finished is type bool true if callback finished without any error.
    // caughtError is whatever recover method returned
}
```
 # Golang panic/recover

```golang
func main() {
	defer func() {
		if caughtError := recover(); caughtError != nil {
            // Do something with caught error.
            // No way to know if error is nil
            // because no panic was called or because
            // panic was called with an nil parameter
		}
	}()}
	thisFuncCanPanic()
	// Awkward logic
}
```
