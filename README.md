# Golang VST package

## Dependencies 
For Go dependencies **dep** is used.

VST SDK is also required
1. Go to [Steinberg web site](https://www.steinberg.net/en/company/developers.html) to get the SDK
2. Extract downloaded archive
3. Go to /VST2_SDK/public.sdk/source/vst2.x
4. Copy VST2 next header files to phono/vendor/vst2
* aeffect.h
* aeffectx.h
* vstfxstore.h

## Building 
1. Go to submodule diretory _vst2
2. Run 
~~~
./build.sh <os> <arch>
~~~

## Releasing 
1. Go to submodule diretory _vst2
2. Run 
~~~
./release.sh <version>
~~~