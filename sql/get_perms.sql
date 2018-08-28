-- ============================================================================
-- get_perms
-- Retrieves permission records from the database.
-- Author: David Haifley
-- Created: 2018-08-14
-- ============================================================================
CREATE OR REPLACE FUNCTION public.get_perms(
	p_id BIGINT DEFAULT NULL,
	p_service CHARACTER VARYING DEFAULT NULL,
	p_name CHARACTER VARYING DEFAULT NULL)
RETURNS TABLE(
	"id" BIGINT,
	"service" CHARACTER VARYING,
	"name" CHARACTER VARYING)
LANGUAGE 'plpgsql'
AS $$
BEGIN
RETURN QUERY
SELECT
	p.id,
	p.service,
	p.name
FROM perm p
WHERE p.id = COALESCE(p_id, p.id)
		AND p.service = COALESCE(p_service, p.service)
		AND p.name = COALESCE(p_name, p.name);
END;
$$

/* Test code:
SELECT save_perm(1, 'test', 'test') AS id
SELECT * FROM get_perms(NULL, 'test')
SELECT delete_perms(NULL, 'test')
*/
