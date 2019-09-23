package api

//#cgo darwin LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
/*
char * MYCFStringCopyUTF8String(CFStringRef aString) {
	if (aString == NULL) {
    	return NULL;
	}

  	CFIndex length = CFStringGetLength(aString);
  	CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
	char *buffer = (char *)malloc(maxSize);
	if (CFStringGetCString(aString, buffer, maxSize, kCFStringEncodingUTF8)) {
    	return buffer;
	}
	free(buffer); // If we failed
	return NULL;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// handle keeps a bundle reference to clean up on close.
type handle struct {
	bundle uintptr
}

// Open loads the plugin entry point into memory. It's CFBundle in OS X.
func Open(path string) (EntryPoint, error) {
	//create C string
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	//convert to CF string
	cfpath := C.CFStringCreateWithCString(0, cpath, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	//get bundle url
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == 0 {
		return EntryPoint{}, fmt.Errorf("Failed to create bundle url at %v", path)
	}
	defer C.CFRelease(C.CFTypeRef(bundleURL))
	//open bundle and release it only if it failed
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	// library = uintptr(bundle)
	//bundle ref should be released in the end of program with plugin.unload call

	//create C string
	cvstMain := C.CString(main)
	defer C.free(unsafe.Pointer(cvstMain))
	//create CF string
	cfvstMain := C.CFStringCreateWithCString(0, cvstMain, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	entryPoint := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if entryPoint == nil {
		C.CFRelease(C.CFTypeRef(C.CFBundleRef(bundle)))
		return EntryPoint{}, fmt.Errorf("Failed to find entry point in bundle %v", path)
	}

	return EntryPoint{
		main: effectMain(entryPoint),
		handle: handle{
			bundle: uintptr(bundle),
		},
	}, nil
}

// close cleans up bundle reference.
func (h handle) close() error {
	C.CFRelease(C.CFTypeRef(C.CFBundleRef(h.bundle)))
	return nil
}

//get string from CFBundle
func getBundleString(bundle C.CFBundleRef, str string) string {
	//create C string
	cstring := C.CString(str)
	defer C.free(unsafe.Pointer(cstring))
	//convert to CF string
	cfstring := C.CFStringCreateWithCString(0, cstring, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfstring))

	bundleString := C.CFStringRef(C.CFBundleGetValueForInfoDictionaryKey(bundle, cfstring))
	defer C.CFRelease(C.CFTypeRef(bundleString))

	convertedString := C.MYCFStringCopyUTF8String(bundleString)
	defer C.free(unsafe.Pointer(convertedString))
	return C.GoString(convertedString)
}
