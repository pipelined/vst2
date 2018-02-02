package vst2

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
import "unsafe"

import "fmt"

//TODO: refactor to plugin.load method
func getEntryPoint(path string) (unsafe.Pointer, error) {
	//create C string
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	//convert to CF string
	cfpath := C.CFStringCreateWithCString(nil, cpath, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	//get bundle url
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == nil {
		return nil, fmt.Errorf("Failed to create bundle url at %v", path)
	}
	defer C.free(unsafe.Pointer(bundleURL))
	//open bundle and release it only if it failed
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	if bundle == nil {
		return nil, fmt.Errorf("Failed to create bundle at %v", path)
	}
	defer C.CFRelease(C.CFTypeRef(bundle))

	//create C string
	cvstMain := C.CString(vstMain)
	defer C.free(unsafe.Pointer(cvstMain))
	//create CF string
	cfvstMain := C.CFStringCreateWithCString(nil, cvstMain, C.kCFStringEncodingUTF8)
	defer C.free(unsafe.Pointer(cfvstMain))

	mainEntryPoint := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if mainEntryPoint == nil {
		defer C.CFBundleUnloadExecutable(bundle)
	}
	return mainEntryPoint, nil
}

//get string from CFBundle
func getBundleString(bundle C.CFBundleRef, str string) string {
	//create C string
	cstring := C.CString(str)
	defer C.free(unsafe.Pointer(cstring))
	//convert to CF string
	cfstring := C.CFStringCreateWithCString(nil, cstring, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfstring))

	bundleString := C.CFStringRef(C.CFBundleGetValueForInfoDictionaryKey(bundle, cfstring))
	defer C.free(unsafe.Pointer(bundleString))

	convertedString := C.MYCFStringCopyUTF8String(bundleString)
	defer C.free(unsafe.Pointer(convertedString))
	return C.GoString(convertedString)
}
