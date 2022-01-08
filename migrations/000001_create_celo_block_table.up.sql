CREATE TABLE IF NOT EXISTS celo_blocks (
   id         SERIAL                           PRIMARY KEY,
   height     BIGINT                           UNIQUE NOT NULL,
   hash       VARCHAR (100)                    NOT NULL,
   timestamp  BIGINT                           NOT NULL,
   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
