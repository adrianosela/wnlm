//go:build windows

package wnlm

import (
	"github.com/go-ole/go-ole"
)

var globalOleConn ole.Connection = ole.Connection{}

// Initialize initializes the COM connection for the current program.
func Initialize() error { return globalOleConn.Initialize() }

// Uninitialize uninitializes the COM connection for the current program.
func Uninitialize() { globalOleConn.Uninitialize() }
