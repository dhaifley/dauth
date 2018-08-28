-- ============================================================================
-- get_users
-- Retrieves user records from the database.
-- Author: David Haifley
-- Created: 2018-08-11
-- ============================================================================
CREATE OR REPLACE FUNCTION public.get_users(
	p_id BIGINT DEFAULT NULL,
	p_user CHARACTER VARYING DEFAULT NULL,
	p_pass CHARACTER VARYING DEFAULT NULL,
	p_name CHARACTER VARYING DEFAULT NULL,
	p_email CHARACTER VARYING DEFAULT NULL)
RETURNS TABLE(
	"id" BIGINT,
	"user" CHARACTER VARYING,
	"pass" CHARACTER VARYING,
	"name" CHARACTER VARYING,
	"email" CHARACTER VARYING)
LANGUAGE 'plpgsql'
AS $$
BEGIN
RETURN QUERY
SELECT
	u.id,
	u.user,
	u.pass,
	u.name,
	u.email
FROM "user" u
WHERE u.id = COALESCE(p_id, u.id)
	AND u.user = COALESCE(p_user, u.user)
	AND u.pass = COALESCE(p_pass, u.pass)
	AND u.name = COALESCE(p_name, u.name)
	AND u.email = COALESCE(p_email, u.email);
END;
$$

/* Test code:
SELECT save_user(1, 'test', 'test', 'test', 'test') AS user
SELECT * FROM get_users(NULL, 'test')
SELECT delete_users(NULL, 'test') AS num
*/
