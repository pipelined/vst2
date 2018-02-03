package vst2

//#cgo darwin LDFLAGS: -framework CoreFoundation
//#include <CoreFoundation/CoreFoundation.h>
//#include "vst2.h"
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

//TODO: refactor to plugin.load method
func (plugin *Plugin) load() error {
	//create C string
	cpath := C.CString(plugin.Path)
	defer C.free(unsafe.Pointer(cpath))
	//convert to CF string
	cfpath := C.CFStringCreateWithCString(nil, cpath, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfpath))

	//get bundle url
	bundleURL := C.CFURLCreateWithFileSystemPath(C.kCFAllocatorDefault, cfpath, C.kCFURLPOSIXPathStyle, C.true)
	if bundleURL == nil {
		return fmt.Errorf("Failed to create bundle url at %v", plugin.Path)
	}
	defer C.free(unsafe.Pointer(bundleURL))
	//open bundle and release it only if it failed
	bundle := C.CFBundleCreate(C.kCFAllocatorDefault, bundleURL)
	if bundle == nil {
		return fmt.Errorf("Failed to create bundle at %v", plugin.Path)
	}
	plugin.library = unsafe.Pointer(bundle)
	//bundle ref should be released in the end of program with plugin.unload call

	//create C string
	cvstMain := C.CString(vstMain)
	defer C.free(unsafe.Pointer(cvstMain))
	//create CF string
	cfvstMain := C.CFStringCreateWithCString(nil, cvstMain, C.kCFStringEncodingUTF8)
	defer C.CFRelease(C.CFTypeRef(cfvstMain))

	mainEntryPoint := unsafe.Pointer(C.CFBundleGetFunctionPointerForName(bundle, cfvstMain))
	if mainEntryPoint == nil {
		plugin.Unload()
		return fmt.Errorf("Failed to find entry point in bundle %v", plugin.Path)
	}
	plugin.Name = getBundleString(bundle, "CFBundleName")
	plugin.effect = C.loadEffect(C.vstPluginFuncPtr(mainEntryPoint))

	return nil
}

//Unload cleans up plugin refs
func (plugin *Plugin) Unload() {
	bundle := C.CFBundleRef(plugin.library)
	C.CFBundleUnloadExecutable(bundle)
	C.CFRelease(C.CFTypeRef(bundle))
	C.free(unsafe.Pointer(plugin.effect))
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
