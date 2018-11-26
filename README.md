# Catch 
Catch is a simple way for inlining Golang recover mechanism.
The implementation is thread safe.
<h2 style="color:red"> Usage in production strongly advised against</h2>

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

    func Switch(s1 string,s2 string) (string,string){
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
|   panicWithNil   |               nil               |      false      |
|      noPanic     |               nil               |      true       |
## Sanitize function
```Go
    // callback can be any func type
    sanitizedFunction := catch.Sanitize(callback)
    retValues, err :=  sanitizedFunction("hello","world")
```
### NB :
sanitizedFunction is typed func(...interface{}) ([]interface{},error).
Compiler checks for arguments types and number do not work.

|    callback      |            retValues            |
|:----------------:|:-------------------------------:|
|   panicWithObj   |                42               |
|   panicWithNil   |               nil               |
|      noPanic     |               nil               |

# Performance cost 
here is the output from go test -bench=. comparing pure go panic/recover to catch.

## Numbers 