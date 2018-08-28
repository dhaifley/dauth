-- ============================================================================
-- get_user_perms
-- Retrieves user permission assignment records from the database.
-- Author: David Haifley
-- Created: 2018-08-14
-- ============================================================================
CREATE OR REPLACE FUNCTION public.get_user_perms(
	p_id BIGINT DEFAULT NULL,
	p_user_id BIGINT DEFAULT NULL,
	p_perm_id BIGINT DEFAULT NULL)
RETURNS TABLE(
	"id" BIGINT,
	"user_id" BIGINT,
	"perm_id" BIGINT)
LANGUAGE 'plpgsql'
AS $$
BEGIN
RETURN QUERY
SELECT
	up.id,
	up.user_id,
	up.perm_id
FROM user_perm up
WHERE up.id = COALESCE(p_id, up.id)
		AND up.user_id = COALESCE(p_user_id, up.user_id)
		AND up.perm_id = COALESCE(p_perm_id, up.perm_id);
END;
$$

/* Test code:
SELECT save_perm(1, 1, 1) AS id
SELECT * FROM perm
SELECT * FROM get_perms(NULL, 1)
SELECT delete_user_perms(NULL, 1) AS num
*/
