# ErrorX
    Error package with useful features for micro-service.


# Features
* Modularized features
* Wrap the error to retrieve the cause of an error
* Formatted printing of errors
* Support error stacktrace
* Define Error Code and Http Status Code from error
* Error ID (or) Name to support localization and unique identify of an error

# Usage

### Define new error with message formatting
```
    err := errorx.NewErrorX("File not found , File : %s", "/tmp/dat")
```
___

### Wrap an error and retrieve cause error
```
    dat, err := ioutil.ReadFile("/tmp/dat")
    if err != nil{
        fileNotFoundErr := errorx.NewErrorX("File not found , File : %s", "/tmp/dat").Wrap(err)
    }

    // retrieving cause error
    type causer interface {
		Cause() error
	}

    errorCauser, ok := err.(causer)
	if ok {
		errRootCause = errorCauser.Cause()
        fmt.Println(errRootCause) // Result = "not found"
	}
```
___

### Retriving Error Code and Http Status Code
```
    type MyCustomError struct{
        *errorx.ErrorX
        *errorx.HttpError
    }

    var err error
    err = &MyCustomError{
        errorx.NewErrorX("File not found , File : %s", "/tmp/dat")
        errorx.ErrorWithHttpStatus(400)
    }

    errWithHttpStatus, ok := err.(errorx.HttpError)
    if ok {
        httpStatus := errWithHttpStatus.StatusCode()
        // First 3 Digits of ErrorCode is HttpStatus of that error.
        fmt.Println(httpStatus) // result = 400 (Bad Request)
    }
```

___


### Localization Support
```
    type MyCustomError struct{
        *errorx.ErrorX
        *errorx.ErrorID
    }

    var err error
    err = &MyCustomError{
        errorx.NewErrorX("File not found , File : %s", "/tmp/dat")
        errorx.NewErrorWithID("file_not_found_error")
    }

    errWithCode, ok := err.(errorx.ErrorID)
    if ok {
        errID := errWithCode.ID()
        fmt.Println(errID) // result = file_not_found_error
    }
```

ErrorID can be used as key in localization. For Example:
##### Sample localization json
```
{
    "en-US": {
        "file_not_found_error": "File not found. File path : {{var_file_path}}"
    }
}
```

##### Sample localization code
```
    type MyCustomError struct{
        *errorx.ErrorX
        *errorx.ErrorID
    }

    var err error
    err = &MyCustomError{
        errorx.NewErrorX("File not found , File : %s", "/tmp/dat")
        errorx.NewErrorWithID("file_not_found_error")
    }

    var args []interface{}
    errWithArgs, ok := err.(errorx.ErrorFormatter)
    if ok {
        args = errWithArgs.GetArgs()
    }

    errWithCode, ok := err.(errorx.ErrorID)
    if ok {
        errID := errWithCode.ID()
        localizedMessage := loclize.Transalte(errID, args...)
        fmt.Println(localizedMessage) // result = "File not found. File path : '/tmp/dat'"
    }
```
*Note: Localization feature need to implement as separate package.*

___

### Error's stacktrace
```
    type MyCustomError struct{
        *errorx.ErrorX
        *errorx.ErrorStacktrace
    }

    var err error
    err = &MyCustomError{
        errorx.NewErrorX("File not found , File : %s", "/tmp/dat")
        errorx.NewErrorWithStackTrace(2,2)
    }
    
    errWithStackTrace, ok := err.(errorx.StackTracer)
    if ok {
        stackTraceJSON := errWithStackTrace.GetStackAsJSON()
        fmt.Println(stackTraceJSON)
    }
```
___


### Check error by type casting

file_not_found_error.go
```
    const(
        code := 400001
        id = "file_not_found_error"
        message = "File not found , File : %s"
    )
    type FileNotFoundError struct{
        baseError *errorx.ErrorX
    }

    func FileNotFoundError(filepath string) *FileNotFoundError{
        return FileNotFoundErr{
            baseError: errorx.New(id, code, message, filepath)
            }
        }
    }

    func (err *FileNotFoundError) Wrap(cause error){
        err.Wrap(cause)
    }
```

___

main.go
```
    
    package main

    import "fmt"

    func A() error{
        return errorx.New("file_not_found_error", 400001, "File not found , File : %s", "/tmp/dat").Wrap(err)
    }

    func B() error{
        return A()
    }

    func main(){
        err := B()
        if err != nil{
            fileNotFoundErr, ok := err.(*FileNotFoundError)
            if ok{
                // This is file not found error
            }
            // Unknown error detected
        }
    }
```