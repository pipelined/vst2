# Go bindings for VST2 SDK

[![GoDoc](https://godoc.org/github.com/pipelined/vst2?status.svg)](https://godoc.org/github.com/pipelined/vst2)
[![Go Report Card](https://goreportcard.com/badge/github.com/pipelined/vst2)](https://goreportcard.com/report/github.com/pipelined/vst2)
[![Build Status](https://travis-ci.org/pipelined/vst2.svg?branch=master)](https://travis-ci.org/pipelined/vst2)

## Dependencies 

VST2 SDK is also required and since it is under commercial lincense, it cannot be a part of public repository. 

To use this package:

1. Go to [Steinberg web site](https://www.steinberg.net/en/company/developers.html) to get the SDK
2. Extract downloaded archive
3. Set environment variable CGO_CFLAGS="-I<PATH TO /VST2_SDK/public.sdk/source/vst2.x>"
* aeffect.h
* aeffectx.h
* vstfxstore.h
4. To build this package, you also need to change next lines in **aeffect.h**:
    ```
    #ifdef  __cplusplus
    #define VST_INLINE inline
    #else
    #define VST_INLINE 
    #endif
    ```
    to 
    ```
    #ifdef  __cplusplus
    #define VST_INLINE inline
    #else
    #define VST_INLINE inline
    #endif
    ```
