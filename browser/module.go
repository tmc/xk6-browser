// Package browser provides an entry point to the browser extension.
package browser

import (
	"github.com/dop251/goja"

	"github.com/grafana/xk6-browser/common"

	k6modules "go.k6.io/k6/js/modules"
)

type (
	// RootModule is the global module instance that will create module
	// instances for each VU.
	RootModule struct {
		PidRegistry        *pidRegistry
		BrowserPool        *browserPool
		BrowserProcessPool *browserProcessPool
	}

	// JSModule exposes the properties available to the JS script.
	JSModule struct {
		Browser *goja.Object
		Devices map[string]common.Device
		Version string
	}

	// ModuleInstance represents an instance of the JS module.
	ModuleInstance struct {
		mod *JSModule
	}
)

var (
	_ k6modules.Module   = &RootModule{}
	_ k6modules.Instance = &ModuleInstance{}
)

// New returns a pointer to a new RootModule instance.
func New() *RootModule {
	return &RootModule{
		PidRegistry:        &pidRegistry{},
		BrowserPool:        &browserPool{},
		BrowserProcessPool: newBrowserProcessPool(),
	}
}

// NewModuleInstance implements the k6modules.Module interface to return
// a new instance for each VU.
func (m *RootModule) NewModuleInstance(vu k6modules.VU) k6modules.Instance {
	return &ModuleInstance{
		mod: &JSModule{
			Browser: mapBrowserToGoja(moduleVU{
				VU:                 vu,
				pidRegistry:        m.PidRegistry,
				browserPool:        m.BrowserPool,
				browserProcessPool: m.BrowserProcessPool,
			}),
			Devices: common.GetDevices(),
		},
	}
}

// Exports returns the exports of the JS module so that it can be used in test
// scripts.
func (mi *ModuleInstance) Exports() k6modules.Exports {
	return k6modules.Exports{Default: mi.mod}
}
