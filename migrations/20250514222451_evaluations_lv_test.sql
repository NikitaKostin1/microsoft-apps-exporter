-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS evaluations_lv_test (
    id             VARCHAR(8)   NOT NULL PRIMARY KEY,
    list_id        VARCHAR(40)  NOT NULL,
    site_id        VARCHAR(100) NOT NULL,
    etag           VARCHAR(45)  NOT NULL,
    gp_avg_score   REAL,
    gp_nickname    VARCHAR(20),
    gp_hrid        VARCHAR(20),
    is_attachments BOOLEAN,
    FOREIGN KEY (list_id) REFERENCES sharepoint_lists(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE evaluations_lv_test;
-- +goose StatementEnd
