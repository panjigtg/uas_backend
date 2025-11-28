-- Hapus tabel berurutan dari yang memiliki dependency
DROP TABLE IF EXISTS achievement_references;
DROP TABLE IF EXISTS lecturers;
DROP TABLE IF EXISTS students;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS roles;

-- Hapus ENUM type
DROP TYPE IF EXISTS achievement_status;
