#GO-ERRORS

##Features
* Support Error Stacktrace
* Support Error Code 
* Support Error Message with arguments
* Support Causer Interface 
* Support JSON formatted stacktrace 

## Example
```
    err1 := clerrors.New("custom error")
    
    err2 := fmt.Errorf("root cause error")
    
    err3 := clerrors.WrapWithCode(err2, 11001, "my error")
	
    err4 := clerrors.WrapWithFormat(err3, 11001, "%s not found", "record")

    err5 := clerrors.WrapWithFormat(err4, 11001, "user record not found")

    fmt.Println(clerrors.Cause(err5))
    fmt.Println(clerrors.GetErrorMessagesWithStack(err5))
    fmt.Println(clerrors.GetErrorMessages(err5))
```