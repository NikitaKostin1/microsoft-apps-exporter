//go:build testing

// Exports internal functions for testing purposes.
// This file is only included in builds with the "testing" tag.
package sync

func DiffFull[T any](existing, incoming []T, getID, getETag func(T) string) (toInsert, toUpdate []T, toDelete []string) {
	return diffFull(existing, incoming, getID, getETag)
}

func DiffDelta[T any](existing, changes []T, getID, getETag func(T) string) (toInsert, toUpdate []T, toDelete []string) {
	return diffDelta(existing, changes, getID, getETag)
}
