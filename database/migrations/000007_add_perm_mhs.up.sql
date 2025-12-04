INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Mahasiswa'
AND p.name IN (
  'achievement:delete',
);