package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "" {
		fmt.Println("usage: gensql <table_name>")
		os.Exit(1)
	}
	tableName := os.Args[1]

	template := `
-- +goose Up
-- +goose StatementBegin
CREATE TABLE TABLE_NAME (
  id BIGSERIAL PRIMARY KEY,

  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION update_updated_at_TABLE_NAME()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_TABLE_NAME_updated_at
    BEFORE UPDATE
    ON
        TABLE_NAME
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_TABLE_NAME();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_TABLE_NAME_updated_at ON TABLE_NAME;
DROP FUNCTION update_updated_at_TABLE_NAME();
DROP TABLE TABLE_NAME;
-- +goose StatementEnd`

	result := strings.ReplaceAll(template, "TABLE_NAME", tableName)
	fmt.Println(result)
}
