//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package configuration

import "sync"

func ResetConfig() {
	once = sync.Once{}
}
