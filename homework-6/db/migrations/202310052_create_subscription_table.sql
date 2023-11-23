-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscription
(
    id UUID PRIMARY KEY,
    subscription_name VARCHAR(255) NOT NULL,
    price             INT          NOT NULL,
    start_date        TIMESTAMP    NOT NULL,
    end_date          TIMESTAMP,
    account_id        UUID          NOT NULL,
    FOREIGN KEY (account_id) REFERENCES bank_account (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN;
ALTER TABLE subscription DROP CONSTRAINT subscription_account_id_fkey;
DROP TABLE subscription;
COMMIT;
-- +goose StatementEnd
