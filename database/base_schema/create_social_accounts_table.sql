
DROP TABLE IF EXISTS social_accounts CASCADE;
CREATE TABLE social_accounts (
	id		        VARCHAR,
	name	        VARCHAR,
	type	        VARCHAR DEFAULT 'unknown',
	user_id			UUID NOT NULL,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	-- email			VARCHAR NOT NULL CHECK(email != ''),
	-- CONSTRAINT account PRIMARY KEY(email, type),
	-- FOREIGN KEY (email) REFERENCES users(email) ON DELETE CASCADE
	CONSTRAINT account PRIMARY KEY(user_id, type),
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- @trigger update_social_accounts
-- @description update social_accounts record
DROP TRIGGER IF EXISTS update_social_accounts ON social_accounts;
CREATE TRIGGER update_social_accounts
    BEFORE UPDATE ON social_accounts
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();
