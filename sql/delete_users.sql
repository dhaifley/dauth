-- ============================================================================
-- delete_users
-- Deletes user records from the database.
-- Author: David Haifley
-- Created: 2018-08-11
-- ============================================================================
CREATE OR REPLACE FUNCTION public.delete_users(
	p_id BIGINT DEFAULT NULL,
	p_user CHARACTER VARYING DEFAULT NULL,
	p_pass CHARACTER VARYING DEFAULT NULL,
	p_name CHARACTER VARYING DEFAULT NULL,
	p_email CHARACTER VARYING DEFAULT NULL)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	num BIGINT;
BEGIN
WITH n AS (
	DELETE FROM "user" u
	WHERE u.id = COALESCE(p_id, u.id)
	AND u.user = COALESCE(p_user, u.user)
	AND u.pass = COALESCE(p_pass, u.pass)
	AND u.name = COALESCE(p_name, u.name)
	AND u.email = COALESCE(p_email, u.email)
	RETURNING *
) SELECT INTO num COUNT(*) FROM n;
RETURN num;
END;
$$

/* Test code:
SELECT save_user(1, 'test', 'test', 'test', 'test') AS user
SELECT * FROM get_users(NULL, 'test')
SELECT delete_users(null, 'test') AS num
*/
