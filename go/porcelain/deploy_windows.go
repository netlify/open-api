//go:build windows

package porcelain

// Windows has no O_NONBLOCK and no FIFOs to guard against.
const openNonblock = 0
