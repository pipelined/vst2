package vst2

//#cgo darwin LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
import "C"
import "unsafe"

import "fmt"

func getEntryPoint(path string) (uintptr, error) {
	//create C string
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	//convert to CF string
	cfpath := C.CFStringCreateWithCString(nil, cpath, C.kCFStringEncodingUTF8)
	defer C.free(unsafe.Pointer(cfpath))

	//get bundle url
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == nil {
		return 0, fmt.Errorf("Failed to create bundle url at %v", path)
	}
	defer C.free(unsafe.Pointer(bundleURL))
	//open bundle
	bundleRef := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	if bundleRef == nil {
		return 0, fmt.Errorf("Failed to create bundle at %v", path)
	}
	defer C.CFRelease(C.CFTypeRef(bundleRef))

	//create C string
	cvstMain := C.CString(vstMain)
	defer C.free(unsafe.Pointer(cvstMain))
	//create CF string
	cfvstMain := C.CFStringCreateWithCString(nil, cvstMain, C.kCFStringEncodingUTF8)
	defer C.free(unsafe.Pointer(cfvstMain))

	mainEntryPoint := uintptr(C.CFBundleGetFunctionPointerForName(bundleRef, cfvstMain))
	return mainEntryPoint, nil
}
