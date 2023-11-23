-- +goose Up
create table bank_account (
    id SERIAL PRIMARY KEY,
    holder_name VARCHAR(255) NOT NULL,
    balance INT NOT NULL,
    opening_date TIMESTAMP NOT NULL,
    bank_name VARCHAR(255) NOT NULL
);

-- +goose Down
DROP TABLE bank_account;