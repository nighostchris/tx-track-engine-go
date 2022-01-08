CREATE TABLE IF NOT EXISTS transactions (
    id              SERIAL          PRIMARY KEY,
    block_height    BIGINT          UNIQUE NOT NULL,
    hash            VARCHAR (100)   NOT NULL,
    origin          VARCHAR (100)   NOT NULL,
    destination     VARCHAR (100)   NOT NULL,
    contract        VARCHAR (100)   NOT NULL,
    value           VARCHAR (100)   NOT NULL,
    type            SMALLINT        NOT NULL,
    memo            TEXT            NOT NULL,
    blockchain      SMALLINT        REFERENCES  blockchains (id),
    timestamp       BIGINT          NOT NULL,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
