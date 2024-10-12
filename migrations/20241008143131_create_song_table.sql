-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS songs (
    "id" SERIAL PRIMARY KEY,
    "created_at" timestamp NOT NULL DEFAULT NOW(),
    "updated_at" timestamp NOT NULL DEFAULT NOW(),
    "song_name" TEXT NOT NULL,
    "song_group" TEXT NOT NULL,
    "song_text" TEXT NOT NULL,
    "release_date" DATE NOT NULL,
    "link" VARCHAR(1024) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE songs;
-- +goose StatementEnd
