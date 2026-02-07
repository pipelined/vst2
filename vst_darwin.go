//go:build !plugin
// +build !plugin

package vst2

//#cgo LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	// FileExtension of VST2 files in macOS.
	FileExtension = ".vst"

	displayNameKey = "CFBundleName"
)

// ScanPaths of Vst2 files
var scanPaths = []string{
	"~/Library/Audio/Plug-Ins/VST",
	"/Library/Audio/Plug-Ins/VST",
}

// Open loads the plugin entry point into memory. It's CFBundle in OS X.
func Open(path string) (*VST, error) {
	// convert to CF string.
	cfpath := C.CFStringCreateWithCString(0, stringToCString(path), C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	// get bundle url.
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == 0 {
		return nil, fmt.Errorf("failed to create bundle url at %v", path)
	}
	defer C.CFRelease(C.CFTypeRef(bundleURL))

	// open bundle and release it only if it failed.
	// bundle ref should be released in the end of program with entryPoint.Close call.
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	if bundle == 0 {
		return nil, fmt.Errorf("failed to create bundle ref at %v", path)
	}

	// create CF string.
	cfvstMain := C.CFStringCreateWithCString(0, stringToCString(main), C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	ep := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if ep == nil {
		C.CFRelease(C.CFTypeRef(C.CFBundleRef(bundle)))
		return nil, fmt.Errorf("failed to find entry point in bundle %v", path)
	}
	return &VST{
		main:   pluginMain(ep),
		handle: uintptr(bundle),
		Name:   getName(bundle),
	}, nil
}

func getName(bundle C.CFBundleRef) string {
	// create CF string.
	cfNameKey := C.CFStringCreateWithCString(0, stringToCString(displayNameKey), C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfNameKey))
	cfName := C.CFBundleGetValueForInfoDictionaryKey(bundle, cfNameKey)
	defer C.CFRelease(cfName)

	return cfStringRefToString(C.CFStringRef(cfName))
}

// Convert CoreFoundation String to golang string.
func cfStringRefToString(cfStr C.CFStringRef) (goString string) {
	l := C.CFStringGetLength(cfStr) // utf-16 length
	buf := make([]byte, l*2)
	C.CFStringGetCString(cfStr, (*C.char)(unsafe.Pointer(&buf[0])), l*2, C.CFStringGetSystemEncoding())
	return string(buf)
}

// Close frees plugin handle.
func (v *VST) Close() error {
	C.CFRelease(C.CFTypeRef(C.CFBundleRef(v.handle)))
	return nil
}
