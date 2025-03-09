//go:build windows

package wnlm

import (
	"syscall"
)

// globalSyscaller represents how this library will perform
// system calls. By default it actually uses the system, but
// it can be overridden for mocking or testing.
var globalSyscaller Syscaller = defaultSyscaller()

// Syscaller represents an entity capable of performing system calls.
type Syscaller interface {
	SyscallN(trap uintptr, args ...uintptr) (r1 uintptr, r2 uintptr, err syscall.Errno)
}

// SetGlobalSyscaller overrides the global Syscaller.
func SetGlobalSyscaller(s Syscaller) { globalSyscaller = s }

// UnsetGlobalSyscaller sets the global Syscaller to default.
func UnsetGlobalSyscaller() { globalSyscaller = defaultSyscaller() }

// syscaller is the default implementation of Syscaller.
// It actually uses the system to perform system calls.
type syscaller struct{}

// defaultSyscaller returns a new default syscaller.
func defaultSyscaller() Syscaller { return &syscaller{} }

// SyscallN performs a system call.
func (s *syscaller) SyscallN(trap uintptr, args ...uintptr) (r1 uintptr, r2 uintptr, err syscall.Errno) {
	return syscall.SyscallN(trap, args...)
}
