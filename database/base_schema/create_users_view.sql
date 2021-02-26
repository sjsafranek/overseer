
DROP VIEW IF EXISTS users_view CASCADE;

CREATE OR REPLACE VIEW users_view AS (
    SELECT
        *,
        json_build_object(
            'id', users.id,
            'email', users.email,
            'username', users.username,
            'is_active', users.is_active,
            'is_deleted', users.is_deleted,
            'created_at', to_char(users.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
            'updated_at', to_char(users.updated_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
            'apikeys', (
                SELECT json_agg(c) FROM (
                    SELECT
                    json_build_object(
                        'user_id', apikeys.user_id,
                        'name', apikeys.name,
                        'apikey', apikeys.apikey,
                        'secret', apikeys.secret,
                        'is_active', apikeys.is_active,
                        'is_deleted', apikeys.is_deleted,
                        'created_at', to_char(apikeys.created_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"'),
                        'updated_at', to_char(apikeys.updated_at, 'YYYY-MM-DD"T"HH:MI:SS"Z"')
                    )
                    FROM apikeys WHERE apikeys.user_id = users.id AND apikeys.is_deleted = false
                ) c
            )
        ) AS user_json
    FROM users
);
