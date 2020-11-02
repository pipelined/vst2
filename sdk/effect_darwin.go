package sdk

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

	displayNameKey = "CFBundleDisplayName"
)

var (
	// ScanPaths of Vst2 files
	scanPaths = []string{
		"~/Library/Audio/Plug-Ins/VST",
		"/Library/Audio/Plug-Ins/VST",
	}
)

// Open loads the plugin entry point into memory. It's CFBundle in OS X.
func Open(path string) (*EntryPoint, error) {
	// create C string.
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	// convert to CF string.
	cfpath := C.CFStringCreateWithCString(0, cpath, C.kCFStringEncodingUTF8)
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

	// create C string.
	cvstMain := C.CString(main)
	defer C.free(unsafe.Pointer(cvstMain))
	// create CF string.
	cfvstMain := C.CFStringCreateWithCString(0, cvstMain, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	ep := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if ep == nil {
		C.CFRelease(C.CFTypeRef(C.CFBundleRef(bundle)))
		return nil, fmt.Errorf("failed to find entry point in bundle %v", path)
	}
	return &EntryPoint{
		main:   effectMain(ep),
		handle: uintptr(bundle),
		Name:   getName(bundle),
	}, nil
}

func getName(bundle C.CFBundleRef) string {
	// create C string.
	nameKey := C.CString(displayNameKey)
	defer C.free(unsafe.Pointer(nameKey))
	// create CF string.
	cfNameKey := C.CFStringCreateWithCString(0, nameKey, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfNameKey))
	cfName := C.CFBundleGetValueForInfoDictionaryKey(bundle, cfNameKey)
	defer C.CFRelease(cfName)

	return cfStrToStr(C.CFStringRef(cfName))
}

func cfStrToStr(cstr C.CFStringRef) (goString string) {
	l := C.CFStringGetLength(cstr) //utf-16 length
	buf := make([]byte, l*2)
	C.CFStringGetCString(cstr, (*C.char)(unsafe.Pointer(&buf[0])), l*2, C.CFStringGetSystemEncoding())
	return string(buf)
}

// Close frees plugin handle.
func (m *EntryPoint) Close() error {
	C.CFRelease(C.CFTypeRef(C.CFBundleRef(m.handle)))
	return nil
}
