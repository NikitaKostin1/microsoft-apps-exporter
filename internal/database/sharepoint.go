package database

import (
	"context"
	"database/sql"
	"fmt"
	"microsoft-apps-exporter/internal/models"
	"strings"
)

/*
Lists
*/

func (db *Database) GetList(ID string) ([]models.ListMetadata, error) {
	query := `
	SELECT 
		id, site_id, etag, name, display_name, delta_link
	FROM sharepoint_lists
	WHERE id = $1;`

	var metadata models.ListMetadata

	err := db.Connection.QueryRowContext(context.Background(), query, ID).Scan(
		&metadata.ID,
		&metadata.SiteID,
		&metadata.ETag,
		&metadata.Name,
		&metadata.DisplayName,
		&metadata.DeltaLink,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return []models.ListMetadata{}, nil
		}
		return []models.ListMetadata{}, err
	}

	return []models.ListMetadata{metadata}, nil
}

func (db *Database) InsertLists(m *[]models.ListMetadata) error {
	query := `
		INSERT INTO sharepoint_lists (
			id, site_id, etag, name, display_name, delta_link
		) VALUES (
			$1, $2, $3, $4, $5, $6
		);`

	return db.withTransaction(func(tx *sql.Tx) error {
		for _, metadata := range *m {
			_, err := tx.ExecContext(context.Background(), query,
				metadata.ID,
				metadata.SiteID,
				metadata.ETag,
				metadata.Name,
				metadata.DisplayName,
				metadata.DeltaLink,
			)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db *Database) UpdateListIgnoreDelta(metadata models.ListMetadata) error {
	query := `
		UPDATE sharepoint_lists
		SET 
			site_id = $2,
			etag = $3,
			name = $4,
			display_name = $5
		WHERE sharepoint_lists.id = $1;
	`
	return db.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.ExecContext(context.Background(), query,
			metadata.ID,
			metadata.SiteID,
			metadata.ETag,
			metadata.Name,
			metadata.DisplayName,
		)
		return err
	})
}

func (db *Database) DeleteList(ID string) error {
	query := `
		DELETE FROM sharepoint_lists
		WHERE id = $1;
	`
	return db.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.ExecContext(context.Background(), query, ID)
		return err
	})
}

func (db *Database) GetDeltaLink(listID string) (*string, error) {
	query := `
		SELECT delta_link
		FROM sharepoint_lists
		WHERE id = $1;`

	var deltaLink sql.NullString

	err := db.Connection.QueryRowContext(context.Background(), query, listID).Scan(&deltaLink)
	if err != nil {
		return nil, err
	}

	if deltaLink.Valid {
		return &deltaLink.String, nil // Return pointer to string if not NULL
	}

	return nil, nil // Return nil if delta_link is NULL
}

func (db *Database) SaveDeltaLink(listID, deltaLink string) error {
	query := `
		UPDATE sharepoint_lists 
		SET delta_link = $2 
		WHERE id = $1;`

	return db.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.ExecContext(context.Background(), query, listID, deltaLink)
		return err
	})
}

func (db *Database) DeleteDeltaLink(listID string) error {
	query := `
		UPDATE sharepoint_lists 
		SET delta_link = NULL 
		WHERE id = $1;`

	return db.withTransaction(func(tx *sql.Tx) error {
		_, err := tx.ExecContext(context.Background(), query, listID)
		return err
	})
}

/*
List Items
*/

func (db *Database) GetListItems(table, siteID, listID string) (*[]models.ListItem, error) {
	query := fmt.Sprintf(`SELECT * FROM %s WHERE site_id = $1 AND list_id = $2;`, table)

	rows, err := db.Connection.QueryContext(context.Background(), query, siteID, listID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var listItems []models.ListItem
	for rows.Next() {
		listItem, err := scanListItem(rows)
		if err != nil {
			return nil, fmt.Errorf("item_id \"%s\": %w", listItem.Metadata.ID, err)
		}
		listItems = append(listItems, listItem)
	}

	return &listItems, nil
}

func (db *Database) InsertListItems(table string, columnsMap map[string]string, listItems *[]models.ListItem) error {
	metadataColumns := models.ListItemMetadata{}.DbColumns()
	fieldsColumns := extractKeys(columnsMap)
	allColumns := append(metadataColumns, fieldsColumns...)
	placeholders := generatePlaceholders(len(allColumns))

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", table, strings.Join(allColumns, ", "), strings.Join(placeholders, ", "))

	return db.withTransaction(func(tx *sql.Tx) error {
		for _, listItem := range *listItems {
			values := append(listItem.Metadata.AsArray(), mapFieldValues(listItem.MappedFields, columnsMap, fieldsColumns)...)
			if _, err := tx.ExecContext(context.Background(), query, values...); err != nil {
				return fmt.Errorf("item_id \"%s\": %w", listItem.Metadata.ID, err)
			}
		}
		return nil
	})
}

func (db *Database) UpdateListItem(table string, columnsMap map[string]string, listItem models.ListItem) error {
	metadataColumns := listItem.Metadata.DbColumns()
	setClauses, values := buildUpdateClauses(metadataColumns, columnsMap, listItem)

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE id = $1;`, table, strings.Join(setClauses, ", "))
	values = append([]interface{}{listItem.Metadata.ID}, values...)

	return db.withTransaction(func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(context.Background(), query, values...); err != nil {
			return fmt.Errorf("item_id \"%s\": %w", listItem.Metadata.ID, err)
		}
		return nil
	})
}

func (db *Database) DeleteListItem(table, ID string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id = $1;`, table)

	return db.withTransaction(func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(context.Background(), query, ID); err != nil {
			return fmt.Errorf("item_id \"%s\": %w", ID, err)
		}
		return nil
	})
}
