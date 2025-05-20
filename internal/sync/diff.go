package sync

// diffFull compares two datasets (existing and incoming) and determines which records should be inserted, updated, or deleted.
// - existing: A slice of existing records, each with an ID and ETag.
// - incoming: A slice of incoming records, each with an ID and ETag.
// - getID: A function that takes an object and returns its ID as a string.
// - getETag: A function that takes an object and returns its ETag as a string.
func diffFull[T any](existing, incoming []T, getID, getETag func(T) string) (toInsert, toUpdate []T, toDelete []string) {
	existingMap := make(map[string]string, len(existing))
	for _, e := range existing {
		existingMap[getID(e)] = getETag(e)
	}

	for _, i := range incoming {
		iID, iEtag := getID(i), getETag(i)

		if eEtag, found := existingMap[iID]; found {
			if eEtag != iEtag {
				toUpdate = append(toUpdate, i) // Update if ETag changed
			}
			delete(existingMap, iID) // Mark as processed
		} else {
			toInsert = append(toInsert, i) // Insert new record
		}
	}

	// Remaining records are deletions
	for eID := range existingMap {
		toDelete = append(toDelete, eID)
	}

	return
}

// diffDelta compares two datasets (existing and changes) and determines which records should be inserted, updated, or deleted in a delta sync.
// - existing: A slice of existing records, each with an ID and ETag.
// - changes: A slice of changed records, each with an ID and ETag.
// - getID: A function that takes an object and returns its ID as a string.
// - getETag: A function that takes an object and returns its ETag as a string.
func diffDelta[T any](existing, changes []T, getID, getETag func(T) string) (toInsert, toUpdate []T, toDelete []string) {
	existingMap := make(map[string]string, len(existing))
	for _, e := range existing {
		existingMap[getID(e)] = getETag(e)
	}

	for _, change := range changes {
		id, eTag := getID(change), getETag(change)

		if eTag == "" {
			toDelete = append(toDelete, id) // Deletion
		} else if _, found := existingMap[id]; found {
			toUpdate = append(toUpdate, change) // Update if changed
		} else {
			toInsert = append(toInsert, change) // Insert new record
		}
	}

	return
}
