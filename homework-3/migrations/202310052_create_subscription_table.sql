-- +goose Up
CREATE TABLE subscription
(
    id                SERIAL PRIMARY KEY,
    subscription_name VARCHAR(255) NOT NULL,
    price             INT          NOT NULL,
    start_date        TIMESTAMP    NOT NULL,
    end_date          TIMESTAMP,
    account_id        INT          NOT NULL,
    FOREIGN KEY (account_id) REFERENCES bank_account (id)
);

-- +goose Down
DROP TABLE subscription;
