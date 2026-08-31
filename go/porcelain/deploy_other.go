//go:build !windows

package porcelain

import "syscall"

// openNonblock keeps a FIFO from hanging a read-open; no effect on regular files.
const openNonblock = syscall.O_NONBLOCK
