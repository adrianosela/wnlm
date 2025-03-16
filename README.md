# wnlm

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/wnlm)](https://goreportcard.com/report/github.com/adrianosela/wnlm)
[![Documentation](https://godoc.org/github.com/adrianosela/wnlm?status.svg)](https://godoc.org/github.com/adrianosela/wnlm)
[![GitHub issues](https://img.shields.io/github/issues/adrianosela/wnlm.svg)](https://github.com/adrianosela/wnlm/issues)
[![license](https://img.shields.io/github/license/adrianosela/wnlm.svg)](https://github.com/adrianosela/wnlm/blob/master/LICENSE)

Go bindings for the Windows [Network List Manager API](https://learn.microsoft.com/en-us/windows/win32/nla/portal).

> The Network List Manager API enables applications to retrieve a list of available network connections. Applications can filter networks, based on attributes and signatures, and choose the networks best suited to their task. The Network List Manager infrastructure notifies applications of changes in the network environment, thus enabling applications to dynamically update network connections.

This is built on top of [github.com/go-ole/go-ole](https://github.com/go-ole/go-ole) (Go bindings for Windows COM using shared libraries instead of cgo), using the Network List Manager API interface definitions from [gitlab.winehq.org/wine/wine](https://gitlab.winehq.org/wine/wine/-/blob/1c4350ac/include/netlistmgr.idl).

### Usage

```
import "github.com/adrianosela/wnlm"

func main () {
    wnlm.Initialize()
    defer wnlm.Uninitialize()

    nlm, err := wnlm.NewNetworkListManager()
    if err != nil {
        // handle err
    }
    defer nlm.Release()

    // do stuff...
}
```

### Examples

- [Enumerate Networks](./_examples_/enumerate_networks/)
