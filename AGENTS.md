# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go library implementing the VST2 SDK API, providing type-safe bindings for building both VST2 audio plugin hosts and plugins without direct CGO/unsafe interaction. Module path: `github.com/cwbudde/vst2`.

## Build & Test Commands

```bash
# Run tests (host mode)
go test --race --coverprofile=coverage.txt --covermode=atomic ./...

# Build plugin (verifies plugin-tagged code compiles)
go build --buildmode c-archive --tags plugin -o demoplugin.a ./demoplugin

# Run a single test
go test --race -run TestName ./...

# Regenerate opcode string methods
go generate ./...

# Lint (uses trunk.io)
trunk check
```

## Architecture

### Dual Build Modes

The package compiles in two mutually exclusive modes controlled by build tags:

- **Host mode** (default, `!plugin` tag): Code for loading and controlling VST2 plugins. Files: `host.go`, `host_callback.go`, `host_export.go`, `vst_linux.go`, `vst_darwin.go`, `vst_windows.go`
- **Plugin mode** (`--tags plugin`): Code for creating VST2 plugins. Files: `plugin.go`, `plugin_export.go`, `plugin_export_windows.go`, `plugin_export_notwin.go`

### Key Types

- **`Plugin`** (plugin mode): Defines a VST2 effect with audio processing callbacks (`ProcessDoubleFunc`/`ProcessFloatFunc`), parameters, and metadata. Allocated via the global `PluginAllocator` function.
- **`Dispatcher`** (plugin mode): Handles host-to-plugin opcodes (parameter queries, buffer size changes, MIDI events, chunk get/set).
- **`Host`** (both modes): Callback struct for querying the host (sample rate, buffer size, time info, process level, display update).
- **`VST`** (host mode): Reference to a loaded plugin, wraps platform-specific library loading.
- **`DoubleBuffer`/`FloatBuffer`**: C-compatible audio buffers with channel-based access. Manual C memory allocation — must be freed.
- **`EventsPtr`/`MIDIEvent`/`SysExMIDIEvent`**: MIDI event types bridging Go and the C event structures.

### C Bridge Layer

The `include/` directory contains C bridge code:

- `include/vst.h` — Main VST SDK header with type definitions
- `include/plugin/plugin.c` — Plugin-side C bridge (entry point, dispatch, process callbacks)
- `include/host/host.c` — Host-side C bridge (callback routing)
- `include/event/event.c` — Event system helpers

Go files use `#include` directives in CGO preambles to pull in the appropriate C code.

### Platform-Specific Plugin Loading (Host Mode)

Each platform has its own `vst_<platform>.go` for loading `.dll`/`.vst`/`.so` files:

- **Windows**: `syscall.LoadDLL` + `NewProc`
- **macOS**: CoreFoundation `CFBundle` API via CGO
- **Linux**: `dlopen`/`dlsym` via CGO

### Pipelined Integration

`processor.go` integrates with the `pipelined.dev/pipe` framework, providing a `Processor` allocator that wraps VST plugins as pipe components with automatic buffer format selection.

### Code Generation

`opcode_string.go` is generated via `go:generate stringer -type=PluginOpcode,HostOpcode`. Do not edit manually.

## CGO Conventions

- C types use `C.int64_t` (not `C.longlong`) for cross-platform correctness.
- `Events.numEvents` is `int32_t` (not `int64_t`).
- Plugin parameters use `C.int64_t` for the value field.
- Memory allocated via `C.CBytes` or `C.calloc` must be explicitly freed.

## CI

GitHub Actions runs on macOS and Windows (not Linux). Tests run in host mode; plugin mode is verified via `c-archive` build of `./demoplugin`.
