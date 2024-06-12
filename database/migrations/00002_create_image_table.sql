-- +goose Up
-- +goose StatementBegin
CREATE TYPE ImageStatus AS ENUM ('PENDING', 'ERROR', 'FINISHED');

CREATE TABLE image (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID REFERENCES auth.users NOT NULL, 
    batch_id UUID NOT NULL,
    prompt TEXT NOT NULL,
    status ImageStatus NOT NULL DEFAULT 'PENDING', 
    URL TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE FUNCTION update_updated_at_image() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = now();

RETURN NEW;

END;

$$ language 'plpgsql';

CREATE TRIGGER update_image_updated_at BEFORE
UPDATE
    ON image FOR EACH ROW EXECUTE PROCEDURE update_updated_at_image();

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TRIGGER update_image_updated_at ON image;

DROP FUNCTION update_updated_at_image();

DROP TABLE image;
DROP TYPE ImageStatus
-- +goose StatementEnd