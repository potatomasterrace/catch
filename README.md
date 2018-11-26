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
|    callback      |            err    value         | err castable to |
|:----------------:|:-------------------------------:|:---------------:|
|   panicWithObj   |                42               |       int       |
|   panicWithNil   | "panic called with a nil error" |      error      |
|      noPanic     |               nil               |        -        |

## Inlining to error
returns the same values as the previous example
``` Go
    // callback has to be type func()
    err := catch.Error(callback)
    // r is type error
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

    BenchmarkWithPanicking/pure_go-4                       1        1090155629 ns/op
    BenchmarkWithPanicking/catch-4                         1        5946460116 ns/op
    BenchmarkWithoutPanicking/pure_go-4             2000000000               0.02 ns/op
    BenchmarkWithoutPanicking/catch-4                      1        3063595345 ns/op
## Bottom line
* catch is about 6 times slower than pure go when no panic happens.
* catch is a **LOT**  slower (litteraly 100 billion times) than pure go if a panic happens.