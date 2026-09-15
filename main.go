package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"strings"
	"sync"

	"github.com/xtls/xray-core/core"
	_ "github.com/xtls/xray-core/main/distro/all"
	"github.com/xtls/xray-core/infra/conf/serial"
)

var (
	serverInstance *core.Instance
	serverMutex    sync.Mutex
)

//export StartCore
func StartCore(configJson *C.char) *C.char {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverInstance != nil {
		_ = serverInstance.Close()
		serverInstance = nil
	}

	rawConfig := C.GoString(configJson)
	config, err := serial.LoadJSONConfig(strings.NewReader(rawConfig))
	if err != nil {
		return C.CString("failed to parse json config: " + err.Error())
	}

	server, err := core.New(config)
	if err != nil {
		return C.CString("failed to create core instance: " + err.Error())
	}

	if err := server.Start(); err != nil {
		return C.CString("failed to start core: " + err.Error())
	}

	serverInstance = server
	return nil
}

//export StopCore
func StopCore() {
	serverMutex.Lock()
	defer serverMutex.Unlock()

	if serverInstance != nil {
		_ = serverInstance.Close()
		serverInstance = nil
	}
}

func main() {}
