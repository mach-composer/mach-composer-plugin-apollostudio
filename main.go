package main

import (
	"github.com/mach-composer/mach-composer-plugin-apollostudio/internal"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
)

func main() {
	p := internal.NewApollostudioPlugin()
	plugin.ServePlugin(p)
}
