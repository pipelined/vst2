package pool

import (
	//"fmt"
	audio "github.com/dudk/phono/audio"
	vst2 "github.com/dudk/phono/vst2"
	"io/ioutil"
	// "log"
	"strings"
	// "os"
)

var (
	pool []audio.Processor
)

func New() []audio.Processor {
	if pool == nil {
	}

	return pool
}

//Loads all plugins in provided directory
func LoadAll(path string) ([]audio.Processor, error) {
	fileList, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	result := make([]audio.Processor, len(fileList))
	for _, fileInfo := range fileList {
		if strings.HasSuffix(fileInfo.Name(), ".dll") {
			plugin, err := vst2.LoadPlugin(path + "/" + fileInfo.Name())
			if err == nil {
				result = append(result, plugin)
			}
		}
	}
	return result, nil
}
