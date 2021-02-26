
DROP TABLE IF EXISTS social_accounts CASCADE;
CREATE TABLE social_accounts (
	id		        VARCHAR,
	name	        VARCHAR,
	type	        VARCHAR DEFAULT 'unknown',
	email			VARCHAR NOT NULL CHECK(email != ''),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	CONSTRAINT account PRIMARY KEY(email, type),
	FOREIGN KEY (email) REFERENCES users(username) ON DELETE CASCADE
);

-- @trigger update_social_accounts
-- @description update social_accounts record
DROP TRIGGER IF EXISTS update_social_accounts ON social_accounts;
CREATE TRIGGER update_social_accounts
    BEFORE UPDATE ON social_accounts
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();
