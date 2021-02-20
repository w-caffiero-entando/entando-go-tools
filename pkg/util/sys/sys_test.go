package sys

import (
    "fmt"
    "os"
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestCaptureGoOutput(t *testing.T) {
    res := CaptureGoOutput(true, false, func() {
        fmt.Print("Hey\nThere")
        _, _ = fmt.Fprint(os.Stderr, "\n\n<<This text should not be captured>>\n\n")
    })
    assert.Equal(t, "Hey\nThere", res.Stdout)
    assert.Equal(t, "", res.Stderr)

    res = CaptureGoOutput(true, true, func() {
        fmt.Print("Hey\nThere")
        _, _ = fmt.Fprint(os.Stderr, "\nGuys")
    })
    assert.Equal(t, "Hey\nThere", res.Stdout)
    assert.Equal(t, "\nGuys", res.Stderr)
}
