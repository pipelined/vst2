# Plan: Complete VST2 Host Support

## Context

The `pipelined.dev/audio/vst2` library already supports basic VST2 plugin hosting (load, configure, process audio, get/set parameters, presets). However, only 5 of 41+ host callback opcodes are implemented, there is no way to send MIDI events to hosted plugins, no editor/UI window support, and no plugin capability queries. This plan extends the library to be a fully-featured VST2 host.

**No C bridge changes are needed.** The existing `dispatchHostBridge` and `hostCallbackBridge` in the C layer are fully generic -- they pass any opcode through. All work is in Go.

**Scope:** Phases 1-4 only. Plugin-side Host struct updates (phase 5) and stringer regeneration (phase 6) are deferred to a future iteration.

---

## Phase 1: Extend Host Callbacks (plugin-to-host opcodes)

**Files:** `host.go`, `host_callback.go`

### 1a. Add callback type aliases to `host.go`

Add these typed function aliases alongside the existing 5:

| Type                         | Signature                             | Purpose                         |
| ---------------------------- | ------------------------------------- | ------------------------------- |
| `HostAutomateFunc`           | `func(index int32, value float32)`    | Plugin automated a parameter    |
| `HostIdleFunc`               | `func()`                              | Plugin requests idle            |
| `HostProcessEventsFunc`      | `func(events *EventsPtr)`             | Plugin sends MIDI to host       |
| `HostIOChangedFunc`          | `func() bool`                         | Plugin I/O config changed       |
| `HostSizeWindowFunc`         | `func(width, height int32) bool`      | Plugin editor wants resize      |
| `HostGetInputLatencyFunc`    | `func() int32`                        | Report input latency (samples)  |
| `HostGetOutputLatencyFunc`   | `func() int32`                        | Report output latency (samples) |
| `HostGetAutomationStateFunc` | `func() AutomationState`              | Report automation state         |
| `HostGetVendorStringFunc`    | `func() string`                       | Host vendor name                |
| `HostGetProductStringFunc`   | `func() string`                       | Host product name               |
| `HostGetVendorVersionFunc`   | `func() int32`                        | Host version                    |
| `HostCanDoFunc`              | `func(HostCanDoString) CanDoResponse` | Host capability query           |
| `HostGetLanguageFunc`        | `func() HostLanguage`                 | Host UI language                |
| `HostGetDirectoryFunc`       | `func() string`                       | Plugin install directory        |
| `HostBeginEditFunc`          | `func(index int32) bool`              | Parameter edit started          |
| `HostEndEditFunc`            | `func(index int32) bool`              | Parameter edit ended            |

Types referenced (`AutomationState`, `HostLanguage`, `HostCanDoString`, `CanDoResponse`, `EventsPtr`) already exist in `vst.go`/`events.go`.

### 1b. Add fields to `Host` struct in `host.go`

Add corresponding fields for each new type alias. Existing fields remain untouched for backward compatibility.

### 1c. Extend `Host.Callback()` switch in `host_callback.go`

Add `case` branches for each new opcode. Every branch nil-checks the callback before calling. Pattern follows existing code exactly.

**Special cases:**

- `HostAutomate`: passes `index` (parameter) and `opt` (new value)
- `HostSizeWindow`: `width` comes from `index`, `height` from `value`
- `HostGetVendorString`/`HostGetProductString`: write to `(*ascii64)(ptr)` using `copyASCII`
- `HostCanDo`: read C string from `ptr` via `C.GoString((*C.char)(ptr))`
- `HostUpdateDisplay`: already has a field in `Host` but is **missing from the switch** -- fix this gap

### 1d. HostGetDirectory -- proper cleanup

Instead of leaking the `C.CString` allocation, add a `dirCStr unsafe.Pointer` field to the `Plugin` struct. On first `HostGetDirectory` callback, allocate the C string and store the pointer. Free it in `Plugin.Close()`. This bounds the allocation to one per plugin lifetime with guaranteed cleanup.

### 1e. HostVersion handling

Currently `HostVersion` is handled in `host_export.go` before the user callback runs. This is correct. No change needed.

---

## Phase 2: Send MIDI Events to Hosted Plugins

**File:** `host_callback.go`

Add one method on `*Plugin`:

```go
func (p *Plugin) SendEvents(events *EventsPtr) {
    p.Dispatch(PlugProcessEvents, 0, 0, unsafe.Pointer(events), 0)
}
```

`PlugProcessEvents` is already defined in `opcode.go`. The caller creates events via the existing `Events()` constructor and is responsible for calling `Free()` afterward.

---

## Phase 3: Editor/UI Window Support

**File:** `host_callback.go`

Add four methods on `*Plugin`:

- **`EditorGetRect() *EditorRectangle`** -- dispatches `PlugEditGetRect`, ptr is `**EditorRectangle` (plugin writes its own pointer). Returns `nil` if no editor.
- **`EditorOpen(window unsafe.Pointer)`** -- dispatches `PlugEditOpen` with native window handle (HWND / NSView\* / X11 Window)
- **`EditorClose()`** -- dispatches `PlugEditClose`
- **`EditorIdle()`** -- dispatches `PlugEditIdle`

`EditorRectangle` struct already exists in `vst.go:678`. All opcode constants already defined.

---

## Phase 4: Plugin Query Methods

**File:** `host_callback.go`

Add methods on `*Plugin`:

| Method                                            | Opcode                    | Returns                     |
| ------------------------------------------------- | ------------------------- | --------------------------- |
| `CanDo(PluginCanDoString) CanDoResponse`          | `PlugCanDo`               | YesCanDo/NoCanDo/MaybeCanDo |
| `GetTailSize() int`                               | `PlugGetTailSize`         | Tail in samples             |
| `GetInputProperties(int) (*PinProperties, bool)`  | `PlugGetInputProperties`  | Channel info                |
| `GetOutputProperties(int) (*PinProperties, bool)` | `PlugGetOutputProperties` | Channel info                |
| `GetVSTVersion() int`                             | `PlugGetVstVersion`       | e.g. 2400                   |
| `GetPluginName() string`                          | `PlugGetPluginName`       | Up to 32 chars              |
| `GetVendorString() string`                        | `PlugGetVendorString`     | Up to 64 chars              |
| `GetProductString() string`                       | `PlugGetProductString`    | Up to 64 chars              |
| `GetVendorVersion() int`                          | `PlugGetVendorVersion`    | Version int                 |
| `GetCategory() PluginCategory`                    | `PlugGetPlugCategory`     | Category enum               |
| `NumInputs() int`                                 | (read `p.p.numInputs`)    | Channel count               |
| `NumOutputs() int`                                | (read `p.p.numOutputs`)   | Channel count               |
| `UniqueID() int32`                                | (read `p.p.uniqueID`)     | Plugin ID                   |
| `InitialDelay() int`                              | (read `p.p.initialDelay`) | Latency samples             |

**Note:** `CanDo` needs `C.CString` (null-terminated) + `defer C.free`, not `stringToCString` (which is not null-terminated).

**Note:** `GetPluginName` uses `ascii32` which lacks a `String()` method. Add `func (s ascii32) String() string` to `opcode.go` for consistency with `ascii8`, `ascii24`, `ascii64`.

---

## Testing Strategy

### Mock-based unit tests (new file: `host_callback_test.go`)

Test the `Host.Callback()` dispatch logic without loading a real plugin. Create a `Host` with specific callback functions, call the returned `HostCallbackFunc` directly with each opcode, and verify return values and side effects.

Tests to add:

- **TestHostCallbackAutomate**: Set `Automate` callback, verify it receives correct `index` and `opt` values
- **TestHostCallbackVendorStrings**: Set `GetVendorString`/`GetProductString`, pass `*ascii64` ptr, verify string is written
- **TestHostCallbackCanDo**: Set `CanDo` callback, pass C string ptr, verify response
- **TestHostCallbackLatency**: Set `GetInputLatency`/`GetOutputLatency`, verify return values
- **TestHostCallbackAutomationState**: Set `GetAutomationState`, verify return
- **TestHostCallbackSizeWindow**: Set `SizeWindow`, verify `index`/`value` passed as width/height
- **TestHostCallbackBeginEndEdit**: Set `BeginEdit`/`EndEdit`, verify index and return
- **TestHostCallbackUpdateDisplay**: Verify the missing case is now handled
- **TestHostCallbackNilSafety**: Create `Host{}` with no callbacks set, call every opcode, verify no panic and default 0 return
- **TestHostCallbackDirectory**: Set `GetDirectory`, verify C string returned and cleanup

### Plugin method tests (extend `host_test.go`)

These require a real plugin. Guarded by platform build tags since TAL-NoiseMaker is only available on macOS/Windows:

- **TestPluginCanDo**: Query `PluginCanReceiveMIDIEvent` on TAL-NoiseMaker
- **TestPluginGetVSTVersion**: Should return 2400
- **TestPluginGetPluginName**: Should return non-empty string
- **TestPluginGetCategory**: Should return `PluginCategorySynth`
- **TestPluginNumInputsOutputs**: Verify channel counts
- **TestPluginSendEvents**: Send MIDI note-on, process audio, verify no crash

---

## Files Modified Summary

| File                    | Changes                                                                                     |
| ----------------------- | ------------------------------------------------------------------------------------------- |
| `host.go`               | +16 type aliases, +16 fields on `Host` struct                                               |
| `host_callback.go`      | +17 case branches in `Callback()`, +19 methods on `*Plugin`, `dirCStr` cleanup in `Close()` |
| `opcode.go`             | +`ascii32.String()` method                                                                  |
| `host_callback_test.go` | New file: mock-based unit tests for all host callbacks                                      |
| `host_test.go`          | Extended: plugin method tests (platform-guarded)                                            |

No changes to: `include/host/host.c`, `include/vst.h`, `host_export.go`, `buffer.go`, `events.go`, `vst.go`, `plugin.go`

---

## Verification

1. **Build host mode**: `go build ./...` (no plugin tag)
2. **Build plugin mode**: `go build --buildmode c-archive --tags plugin -o demoplugin.a ./demoplugin`
3. **Run all tests**: `go test --race ./...`
4. **Backward compat**: Existing `Host{}` literals and `NoopHostCallback` must compile unchanged
5. **Mock tests pass on Linux**: `host_callback_test.go` runs on all platforms without a real plugin
