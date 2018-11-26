# Catch 
Catch is a simple way for inlining Golang recover mechanism.
The implementation is thread safe.
<h2 style="color:red"> Do not use it before reading the perfomance cost section</h2>

# Install :
    go get github.com/potatomasterrace/catch
# Usage
This functions illustrate the cases : 
```Go
    func panicWithObj(){
        panic(42)
    }
    func panicWithNil(){
        panic(nil)
    }

    func noPanic(){
        // does stuff but never panic
    }
    func functionThatPanics(string,string)(string,string){
        panic(42)
    }

    func functionThatSwitches(s1 string,s2 string) (string,string){
        // switching values (see Sanitize function)
        return s2,s1
    }
    
```
## Inlining to interface{}
``` Go
    // callback has to be type func()
    err := catch.Interface(callback)
    // err is type interface{}
```
|    callback      |            err    value         | err typed       |
|:----------------:|:-------------------------------:|:---------------:|
|   panicWithObj   |                42               |       int       |
|   panicWithNil   | "panic called with a nil error" |      error      |
|      noPanic     |               nil               |        -        |

## Inlining to error
returns the same values as the previous example
``` Go
    // callback has to be type func()
    err := catch.Error(callback)
    // err is type error
```
## Inlining with details
``` Go
    panicked, err := catch.Panic(callback)
    // err is type interface{}
    // panicked is type bool 
```
|    callback      |            err    value         | panicked value  |
|:----------------:|:-------------------------------:|:---------------:|
|   panicWithObj   |                42               |      true       |
|   panicWithNil   |               nil               |      true       |
|      noPanic     |               nil               |      false      |
## Sanitize function
```Go
    // callback can be any func type
    sanitizedFunction := catch.Sanitize(callback)
    retValues, err :=  sanitizedFunction("hello","world")
```
### NB :
**no compiler checks for arguments types and number of arguments.**
sanitizedFunction is typed func(...interface{}) ([]interface{},error).


|    callback            |            retValues            |            err value            |
|:----------------------:|:-------------------------------:|:-------------------------------:|
|   panicWithNil         |               nil               | "panic called with a nil error" |
|   functionThatPanics   |               nil               |              42                 |
|   functionThatSwitches |        ["world","hello"]        |             nil                 |

# Performance cost 
here is the output from go test -bench=. comparing panic/recover to catch.

    BenchmarkWithPanic/single_goroutine/pure_go-4                   2000000000               0.02 ns/op
    BenchmarkWithPanic/single_goroutine/catch-4                     2000000000               0.12 ns/op
    BenchmarkWithPanic/multiple_goroutines/pure_go-4                2000000000               0.04 ns/op
    BenchmarkWithPanic/multiple_goroutines/catch-4                  1000000000               0.20 ns/op
    BenchmarkWithoutPanic/single_goroutine/pure_go-4                2000000000               0.01 ns/op
    BenchmarkWithoutPanic/single_goroutine/catch-4                  2000000000               0.11 ns/op
    BenchmarkWithoutPanic/multiple_goroutines/pure_go-4             2000000000               0.03 ns/op
    BenchmarkWithoutPanic/multiple_goroutines/catch-4               2000000000               0.08 ns/op
    
## Bottom line
* catch is about 5 to 6 times slower than pure go when no panic happens.
* catch is about 3 to 11 times slower than pure go if a panic happens.