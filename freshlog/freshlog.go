//freshlog.go
package freshlog

import (
    "log"
    "os"
    "io/ioutil"
)

var (
    Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Warn = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Fatal = log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
)

func SetLogLevel(log_level string) {
    level := nameToLevel(log_level)
    if level < 5 {
        Debug = log.New(ioutil.Discard, "", log.LstdFlags)
    }
    if level < 4 {
        Info = log.New(ioutil.Discard, "", log.LstdFlags)
    }
    if level < 3 {
        Warn = log.New(ioutil.Discard, "", log.LstdFlags)
    }
    if level < 2 {
        Error = log.New(ioutil.Discard, "", log.LstdFlags)
    }
    return
}

func nameToLevel(level string) int {
    switch level {
    case "fatal":
        return 1
    case "error":
        return 2
    case "warning":
        return 3
    case "info":
        return 4
    case "debug":
        return 5
    default:
        return 5
    }
}
