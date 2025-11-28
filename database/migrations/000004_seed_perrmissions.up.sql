-- Create permissions
INSERT INTO permissions (name, resource, action) VALUES
('achievement:create', 'achievement', 'create'),
('achievement:read', 'achievement', 'read'),
('achievement:update', 'achievement', 'update'),
('achievement:delete', 'achievement', 'delete'),
('achievement:verify', 'achievement', 'verify'),
('user:manage', 'user', 'manage');


-- ADMIN full access
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id FROM roles r, permissions p
WHERE r.name = 'Admin';

-- MAHASISWA
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Mahasiswa'
AND p.name IN (
  'achievement:create',
  'achievement:read',
  'achievement:update'
);

-- DOSEN WALI
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r, permissions p
WHERE r.name = 'Dosen Wali'
AND p.name IN (
  'achievement:read',
  'achievement:verify'
);
