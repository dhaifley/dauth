-- ============================================================================
-- save_user
-- Saves a user record into the database.
-- Author: David Haifley
-- Created: 2018-08-06
-- ============================================================================
CREATE OR REPLACE FUNCTION public.save_user(
	p_id BIGINT,
	p_user CHARACTER VARYING,
	p_pass CHARACTER VARYING,
	p_name CHARACTER VARYING DEFAULT NULL,
	p_email CHARACTER VARYING DEFAULT NULL)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	new_id BIGINT;
BEGIN
	DELETE FROM "user" u WHERE u.id = p_id;
	DELETE FROM "user" u WHERE u.user = p_user;
	SELECT INTO new_id COALESCE(MAX(u.id), 0) + 1 FROM "user" u;
	INSERT INTO "user" ("id", "user", pass, name, email)
		VALUES (new_id, p_user, p_pass, p_name, p_email);
	RETURN new_id AS "id";
END;
$$;

/* Test code:
SELECT save_user(1, 'test', 'test', 'test', 'test') AS id
SELECT * FROM get_users(NULL, 'test')
SELECT delete_users(NULL, 'test') AS num
*/
