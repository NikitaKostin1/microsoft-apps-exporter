-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS sharepoint_lists (
    id             VARCHAR(40)  NOT NULL PRIMARY KEY,
    site_id        VARCHAR(100) NOT NULL,
    etag           VARCHAR(45)  NOT NULL,
    name           VARCHAR(40)  NOT NULL,
    display_name   VARCHAR(40)  NOT NULL,
    delta_link     TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE sharepoint_lists;
-- +goose StatementEnd
