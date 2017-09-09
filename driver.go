package console

import (
	"fmt"
	"sync"
)

type Driver func(host string) (Sessions, error)

var (
	driversMu sync.RWMutex
	drivers   = map[string]Driver{}
)

func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("console: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("console: Register called twice for driver " + name)
	}
	drivers[name] = driver
}

func GetDrivers(name, host string) (Sessions, error) {
	driversMu.RLock()
	defer driversMu.RUnlock()
	dri, ok := drivers[name]
	if !ok {
		return nil, fmt.Errorf("console: Unknown driver %q (forgotten import?)", name)
	}
	return dri(host)
}
