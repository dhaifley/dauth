-- ============================================================================
-- save_user_perm
-- Saves a user permission assignment record into the database.
-- Author: David Haifley
-- Created: 2018-08-14
-- ============================================================================
CREATE OR REPLACE FUNCTION public.save_user_perm(
	p_id BIGINT,
	p_user_id BIGINT,
	p_perm_id BIGINT)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	new_id BIGINT;
BEGIN
	DELETE FROM user_perm up WHERE up.id = up_id;
	DELETE FROM user_perm up WHERE up.user_id = p_user_id
		AND up.perm_id = p_perm_id;
	SELECT INTO new_id COALESCE(MAX(up.id), 0) + 1 FROM user_perm up;
	INSERT INTO user_perm ("id", user_id, perm_id)
		VALUES (new_id, p_user_id, p_perm_id);
	RETURN new_id AS "id";
END;
$$;

/* Test code:
SELECT save_user_perm(1, 1, 1) AS id
SELECT * FROM user_perm
SELECT delete_user_perms(NULL, 'test') AS num
*/
