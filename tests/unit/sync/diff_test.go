//go:build testing && unit

package sync_test

import (
	"microsoft-apps-exporter/internal/sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRecord struct {
	ID   string
	ETag string
}

func getID(r testRecord) string {
	return r.ID
}

func getETag(r testRecord) string {
	return r.ETag
}

func TestDiffFull(t *testing.T) {
	existing := []testRecord{
		{"1", "A"}, // Deleted
		{"2", "B"},
		{"3", "C"},
	}
	incoming := []testRecord{
		{"2", "B"}, // Unchanged
		{"3", "D"}, // Updated
		{"4", "E"}, // New
	}

	toInsert, toUpdate, toDelete := sync.DiffFull(existing, incoming, getID, getETag)

	assert.ElementsMatch(t, []testRecord{{"4", "E"}}, toInsert, "Expected insert records mismatch")
	assert.ElementsMatch(t, []testRecord{{"3", "D"}}, toUpdate, "Expected update records mismatch")
	assert.ElementsMatch(t, []string{"1"}, toDelete, "Expected delete records mismatch")
}

func TestDiffDelta(t *testing.T) {
	existing := []testRecord{
		{"1", "A"},
		{"2", "B"},
		{"3", "C"},
		{"4", "D"}, // Unchanged
	}
	changes := []testRecord{
		{"1", ""},  // Deleted
		{"2", "B"}, // Updated
		{"3", "E"}, // Updated
		{"5", "F"}, // New
	}

	toInsert, toUpdate, toDelete := sync.DiffDelta(existing, changes, getID, getETag)

	assert.ElementsMatch(t, []testRecord{{"5", "F"}}, toInsert, "Expected insert records mismatch")
	assert.ElementsMatch(t, []testRecord{{"2", "B"}, {"3", "E"}}, toUpdate, "Expected update records mismatch")
	assert.ElementsMatch(t, []string{"1"}, toDelete, "Expected delete records mismatch")
}
