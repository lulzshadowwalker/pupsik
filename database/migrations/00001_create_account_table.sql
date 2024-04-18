-- +goose Up
-- +goose StatementBegin
CREATE TABLE account (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users NOT NULL UNIQUE, 
    user_name TEXT NOT NULL,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION update_updated_at_account() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();

RETURN NEW;

END;

$$ language 'plpgsql';

CREATE TRIGGER update_account_updated_at BEFORE
UPDATE
    ON account FOR EACH ROW EXECUTE PROCEDURE update_updated_at_account();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_account_updated_at ON account;

DROP FUNCTION update_updated_at_account();

DROP TABLE account;
-- +goose StatementEnd