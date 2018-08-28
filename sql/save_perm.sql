-- ============================================================================
-- save_perm
-- Saves a permission record into the database.
-- Author: David Haifley
-- Created: 2018-08-14
-- ============================================================================
CREATE OR REPLACE FUNCTION public.save_perm(
	p_id BIGINT,
	p_service CHARACTER VARYING,
	p_name CHARACTER VARYING)
RETURNS BIGINT
LANGUAGE 'plpgsql'
AS $$
DECLARE
	new_id BIGINT;
BEGIN
	DELETE FROM perm p WHERE p.id = p_id;
	DELETE FROM perm p WHERE p.service = p_service
		AND p.name = p_name;
	SELECT INTO new_id COALESCE(MAX(p.id), 0) + 1 FROM perm p;
	INSERT INTO perm ("id", service, name)
		VALUES (new_id, p_service, p_name);
	RETURN new_id AS "id";
END;
$$;

/* Test code:
SELECT save_perm(1, 'test', 'test') AS id
SELECT * FROM perm
SELECT delete_perms(NULL, 'test') AS num
*/
