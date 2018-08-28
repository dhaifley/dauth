-- ============================================================================
-- delete_tokens
-- Deletes token records from the database.
-- Author: David Haifley
-- Created: 2018-08-08
-- ============================================================================
CREATE OR REPLACE FUNCTION public.delete_tokens(
	p_id BIGINT DEFAULT NULL,
	p_token CHARACTER VARYING DEFAULT NULL,
	p_user_id BIGINT DEFAULT NULL,
	p_created TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_expires TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_start TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_end TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_old TIMESTAMP WITH TIME ZONE DEFAULT NULL)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	num BIGINT;
BEGIN
WITH n AS (
	DELETE FROM token t
	WHERE t.id = COALESCE(p_id, t.id)
		AND t.token = COALESCE(p_token, t.token)
		AND t.user_id = COALESCE(p_user_id, t.user_id)
		AND t.created = COALESCE(p_created, t.created)
		AND t.expires = COALESCE(p_expires, t.expires)
 		AND t.expires BETWEEN COALESCE(p_start, t.expires) AND COALESCE(p_end, t.expires)
 		AND t.expires < COALESCE(p_old, TIMESTAMP WITH TIME ZONE '12/31/2999')
	RETURNING *
) SELECT INTO num COUNT(*) FROM n;
RETURN num;
END;
$$

/* Test code:
SELECT save_token(-1, 'test', 1, now(), now()) AS id
SELECT * FROM token
SELECT delete_tokens(NULL, 'test') AS num
*/
