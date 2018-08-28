-- ============================================================================
-- save_token
-- Saves a token record into the database.
-- Author: David Haifley
-- Created: 2018-08-06
-- ============================================================================
CREATE OR REPLACE FUNCTION public.save_token(
	p_id BIGINT,
	p_token CHARACTER VARYING,
	p_user_id BIGINT,
	p_created TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_expires TIMESTAMP WITH TIME ZONE DEFAULT NULL)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	new_id BIGINT;
BEGIN
	DELETE FROM token t WHERE t.id = p_id;
	DELETE FROM token t	WHERE t.token = p_token;
	SELECT INTO new_id COALESCE(MAX(t.id), 0) + 1 FROM token t;
	INSERT INTO token ("id", "token", user_id, created, expires)
		VALUES (new_id, p_token, p_user_id, p_created, p_expires);
	RETURN new_id AS "id";
END;
$$;

/* Test code:
SELECT save_token(-1, 'test', 1, now(), now()) AS id
SELECT * FROM token
SELECT delete_tokens(NULL, 'test') AS num
*/
