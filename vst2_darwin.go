package vst2

//#cgo darwin LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
import "C"
import (
	"fmt"
	"unsafe"
)

const (
	// Extension of Vst2 files
	Extension = ".vst"
)

var (
	// ScanPaths of Vst2 files
	scanPaths = []string{
		"~/Library/Audio/Plug-Ins/VST",
		"/Library/Audio/Plug-Ins/VST",
	}
)

// handle keeps a bundle reference to clean up on close.
type handle struct {
	bundle uintptr
}

// Open loads the plugin entry point into memory. It's CFBundle in OS X.
func open(path string) (effectMain, handle, error) {
	// create C string.
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	// convert to CF string.
	cfpath := C.CFStringCreateWithCString(0, cpath, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	// get bundle url.
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == 0 {
		return nil, handle{}, fmt.Errorf("failed to create bundle url at %v", path)
	}
	defer C.CFRelease(C.CFTypeRef(bundleURL))

	// open bundle and release it only if it failed.
	// bundle ref should be released in the end of program with entryPoint.Close call.
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	if bundle == 0 {
		return nil, handle{}, fmt.Errorf("failed to create bundle ref at %v", path)
	}

	// create C string.
	cvstMain := C.CString(main)
	defer C.free(unsafe.Pointer(cvstMain))
	// create CF string.
	cfvstMain := C.CFStringCreateWithCString(0, cvstMain, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	ep := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if ep == nil {
		C.CFRelease(C.CFTypeRef(C.CFBundleRef(bundle)))
		return nil, handle{}, fmt.Errorf("failed to find entry point in bundle %v", path)
	}

	return effectMain(ep), handle{bundle: uintptr(bundle)}, nil
}

// close cleans up bundle reference.
func (h handle) close() error {
	C.CFRelease(C.CFTypeRef(C.CFBundleRef(h.bundle)))
	return nil
}
