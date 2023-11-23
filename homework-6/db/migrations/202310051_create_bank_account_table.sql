-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table bank_account
(
    id           UUID PRIMARY KEY,
    holder_name  VARCHAR(255) NOT NULL,
    balance      INT          NOT NULL,
    opening_date TIMESTAMP    NOT NULL,
    bank_name    VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bank_account;
-- +goose StatementEnd