
DROP TABLE IF EXISTS apikeys CASCADE;

-- @table apikeys
-- @description stores apikeys for users
CREATE TABLE IF NOT EXISTS apikeys (
	user_id			UUID NOT NULL,
    name            VARCHAR NOT NULL,
    apikey          VARCHAR(32) NOT NULL PRIMARY KEY DEFAULT md5(random()::text),
    secret    		VARCHAR(32) NOT NULL DEFAULT md5(random()::text),
	is_active       BOOLEAN DEFAULT TRUE,
    is_deleted      BOOLEAN DEFAULT FALSE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

COMMENT ON TABLE apikeys IS 'User apikeys';

-- @trigger update_apikeys
-- @description update apikey record
DROP TRIGGER IF EXISTS update_apikeys ON apikeys;
CREATE TRIGGER update_apikeys
    BEFORE UPDATE ON apikeys
        FOR EACH ROW
            EXECUTE PROCEDURE update_modified_column();
