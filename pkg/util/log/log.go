package log

import (
    "fmt"
    "os"
    "runtime"
    "strings"
    "time"
)

// The depth of a log line in the code stack
type CodeDepth int

const (
    Default     = CodeDepth(0)
    Business    = CodeDepth(-100)
    Domain      = CodeDepth(-200)
    Generic     = CodeDepth(-300)
    System      = CodeDepth(-400)
    AbyssalCode = CodeDepth(-1000) // thou shalt never see dese domains
)

const DefaultLogLevel = Generic

type Logger struct {
    LogDevice           *os.File
    LogVisibleFromDepth CodeDepth
    DefaultLoggingDepth CodeDepth
    Qualifier           string
    PrintMethod         bool
}

func (log Logger) At(depth CodeDepth) Logger {
    return Logger{
        LogDevice:           log.LogDevice,
        LogVisibleFromDepth: log.LogVisibleFromDepth,
        DefaultLoggingDepth: depth,
        PrintMethod:         depth <= Generic,
    }
}

func (log *Logger) VisibilityThreshold(depth CodeDepth) {
    log.LogVisibleFromDepth = depth
}

func (log *Logger) Error(f string, msgs ...interface{}) {
    log.printLogLine(os.Stderr, "E", getCaller(2), f, msgs...)
}

func (log *Logger) Warn(f string, msgs ...interface{}) {
    log.printLogLine(os.Stderr, "W", getCaller(2), f, msgs...)
}

func (log *Logger) Trace(f string, msgs ...interface{}) {
    log.printLogLine(os.Stdout, "T", getCaller(2), f, msgs...)
}

func (log *Logger) Info(f string, msgs ...interface{}) {
    log.printLogLine(os.Stdout, "I", getCaller(2), f, msgs...)
}

func (log *Logger) Debug(f string, msgs ...interface{}) {
    log.printLogLine(os.Stdout, "D", getCaller(2), f, msgs...)
}

func (log *Logger) printLogLine(out *os.File, tp string, caller string, format string, args ...interface{}) {
    depth := log.DefaultLoggingDepth
    if out == nil {
        out = log.LogDevice
    }
    visibleFrom := log.LogVisibleFromDepth
    if log.LogVisibleFromDepth == Default {
        visibleFrom = DefaultLogLevel
    }
    if visibleFrom < depth {
        return
    }
    tmp := fmt.Sprintf(format, args...)
    methodPart := ""
    if log.PrintMethod {
        methodPart = "\t\t⇐(" + caller + ")"
    }
    tmp = fmt.Sprintf("➤ %s | [%s] | %s%s", tp, time.Now().Format("2006-01-02 15:04:05"), tmp, methodPart)
    _, _ = fmt.Fprintln(out, tmp)
}

func (log *Logger) NotImpl(format string, args ...interface{}) {
    caller := getCaller(2)
    log.printLogLine(os.Stderr, "W", caller, `/!\ NOT IMPLEMENTED /!\ | `+caller+` | `+format, args)
}

func getCaller(backSteps int) string {
    pc, _, _, _ := runtime.Caller(backSteps)
    fun := runtime.FuncForPC(pc)
    arr := strings.Split(fun.Name(), "/")
    return arr[len(arr)-1]
}
