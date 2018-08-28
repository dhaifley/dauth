-- ============================================================================
-- get_tokens
-- Retrieves token records from the database.
-- Author: David Haifley
-- Created: 2018-08-10
-- ============================================================================
CREATE OR REPLACE FUNCTION public.get_tokens(
	p_id BIGINT DEFAULT NULL,
	p_token CHARACTER VARYING DEFAULT NULL,
	p_user_id BIGINT DEFAULT NULL,
	p_created TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_expires TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_start TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_end TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	p_old TIMESTAMP WITH TIME ZONE DEFAULT NULL)
RETURNS TABLE(
	"id" BIGINT,
	"token" CHARACTER VARYING,
	"user_id" CHARACTER VARYING,
	"created" TIMESTAMP WITH TIME ZONE,
	"expires" TIMESTAMP WITH TIME ZONE)
LANGUAGE 'plpgsql'
AS $$
BEGIN
RETURN QUERY
SELECT
	t.id,
	t.token,
	t.user_id,
	t.created,
	t.expires
FROM token t
WHERE t.id = COALESCE(p_id, t.id)
	AND t.token = COALESCE(p_token, t.token)
	AND t.user_id = COALESCE(p_user_id, t.user_id)
	AND t.created = COALESCE(p_created, t.created)
	AND t.expires = COALESCE(p_expires, t.expires)
	AND t.expires BETWEEN COALESCE(p_start, t.expires) AND COALESCE(p_end, t.expires)
	AND t.expires < COALESCE(p_old, TIMESTAMP WITH TIME ZONE '12/31/2999');
END;
$$

/* Test code:
SELECT save_token(-1, 'test', 1, now(), now()) AS id
SELECT * FROM get_tokens(NULL, 'test')
SELECT delete_tokens(NULL, 'test') AS num
*/
