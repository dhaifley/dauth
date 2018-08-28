-- ============================================================================
-- delete_user_perms
-- Deletes user permission assignment records from the database.
-- Author: David Haifley
-- Created: 2018-08-14
-- ============================================================================
CREATE OR REPLACE FUNCTION public.delete_user_perms(
	p_id BIGINT DEFAULT NULL,
	p_user_id BIGINT DEFAULT NULL,
	p_perm_id BIGINT DEFAULT NULL)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	num BIGINT;
BEGIN
WITH n AS (
	DELETE FROM user_perm up
	WHERE up.id = COALESCE(p_id, up.id)
		AND up.user_id = COALESCE(p_user_id, up.user_id)
		AND up.perm_id = COALESCE(p_perm_id, up.perm_id)
	RETURNING *
) SELECT INTO num COUNT(*) FROM n;
RETURN num;
END;
$$

/* Test code:
SELECT save_user_perm(1, 1, 1) AS id
SELECT * FROM get_user_perms(NULL, 1)
SELECT delete_user_perms(NULL, 1) AS num
*/
