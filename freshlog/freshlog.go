//freshlog.go
package freshlog

import (
    "log"
    "os"
)

var (
    Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Warn = log.New(os.Stdout, "[WARNING] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
    Fatal = log.New(os.Stderr, "[FATAL] ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
)
