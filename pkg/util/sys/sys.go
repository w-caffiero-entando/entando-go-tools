package sys

import (
    "fmt"
    "io/ioutil"
    "os"
    "os/exec"
    "path"
    "runtime"
    "time"
)

// Detects runtime the current OS
//
func RunOsDetection() error {
    switch osName := runtime.GOOS; osName {
    case "linux":
        OsInfo = OsInfoType{true, "linux", "/", "\n", true, false, false}
    case "darwin":
        OsInfo = OsInfoType{true, "mac", "/", "\r", false, true, false}
    case "windows":
        OsInfo = OsInfoType{true, "win", "\\", "\r\n", false, false, true}
    default:
        return fmt.Errorf("os \"%s\" not supported", osName)
    }
    return nil
}

type OsInfoType struct {
    initialized bool
    OsType      string
    EOL         string
    Separator   string
    IsLinux     bool
    IsMac       bool
    IsWin       bool
}

var OsInfo OsInfoType

// Tells whether a command in available in the path
//
func IsCommandAvailable(cmd string) bool {
    _, err := exec.LookPath(cmd)
    return err == nil
}

// Captures the standard output of the go print utilities
//
type CaptureGoOutputRes struct {
    Stdout string
    Stderr string
}

func CaptureGoOutput(stdout bool, stderr bool, f func()) CaptureGoOutputRes {
    var or, ow, op, er, ew, ep *os.File
    if stderr {
        ep = os.Stderr
        er, ew, _ = os.Pipe()
        os.Stderr = ew
    }
    if stdout {
        op = os.Stdout
        or, ow, _ = os.Pipe()
        os.Stdout = ow
    }
    f()
    var capturedStdout, capturedStderr []byte
    if stdout {
        _ = ow.Close()
        os.Stdout = op
        capturedStdout, _ = ioutil.ReadAll(or)
    }
    if stdout {
        _ = ew.Close()
        os.Stderr = ep
        capturedStderr, _ = ioutil.ReadAll(er)
    }
    return CaptureGoOutputRes{Stdout: string(capturedStdout), Stderr: string(capturedStderr)}
}

// Tells the directory of the source file that runs it
// Useful for locating test resources
//
func DetermineCurrentSourceFileDir() string {
    _, filename, _, _ := runtime.Caller(1)
    return path.Dir(filename)
}

// Puts the user "on hold" when running a log running tasks
//  - desc       the description of the task used by errors and the "on hold" message
//  - mohTime    time after which the "on hold" message is displayed
//  - timeout    max wait for the task completion
//  - task       the task to run
//  - cleanup    the cleanup function in case of timeout
//
func ExecLongRunningTask(desc string, mohTime time.Duration, timeout time.Duration,
    task func() (interface{}, error), cleanup func()) (LongRunningTaskResult, error) {
    //-
    LSPL := "          "
    LSPL = LSPL + LSPL + LSPL + LSPL + LSPL + LSPL + LSPL + LSPL

    r := make(chan LongRunningTaskResult)

    start := time.Now()

    go func() {
        defer close(r)
        res, err := task()
        r <- LongRunningTaskResult{Result: res, Err: err}
    }()

    nextPrintTime := int64(1000)

    for {
        elapsed := time.Since(start)
        select {
        case res := <-r:
            res.Duration = elapsed
            if elapsed > mohTime {
                _, _ = fmt.Fprintf(os.Stderr, "\r%s\r", LSPL)
            }
            return res, nil
        default:
            elapsedMillis := elapsed.Milliseconds()
            if elapsed > mohTime && elapsedMillis > nextPrintTime {
                nextPrintTime += 1000
                _, _ = fmt.Fprintf(os.Stderr, "\r%03d | Still running %s..%s", elapsedMillis/1000, desc, LSPL)
            }
        }
        if elapsed > timeout {
            if cleanup != nil {
                cleanup()
            }
            return LongRunningTaskResult{Duration: elapsed}, fmt.Errorf("timeout running the task \"%s\"", desc)
        }
        time.Sleep(250 * time.Millisecond)
    }
}

type LongRunningTaskResult struct {
    Err      error
    Result   interface{}
    Duration time.Duration
}
