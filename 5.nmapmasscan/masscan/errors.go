package masscan

import (
	"errors"
)

var (
	// ErrNmapNotInstalled means that upon trying to manually locate nmap in the user's path,
	// it was not found. Either use the WithBinaryPath method to set it manually, or make sure that
	// the nmap binary is present in the user's $PATH.
	ErrNmapNotInstalled = errors.New("masscan binary was not found")

	// ErrScanTimeout means that the provided context was done before the scanner finished its scan.
	ErrScanTimeout = errors.New("masscan scan timed out")

	// ErrMallocFailed means that nmap crashed due to insufficient memory, which may happen on large target networks.
	ErrMallocFailed = errors.New("malloc failed, probably out of space")

	// ErrParseOutput means that nmap's output was not parsed successfully.
	ErrParseOutput = errors.New("unable to parse masscan output, see warnings for details")

	// ErrResolveName means that Nmap could not resolve a name.
	ErrResolveName = errors.New("masscan could not resolve a name")
)
