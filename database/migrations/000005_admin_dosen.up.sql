-- USER ADMIN
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
VALUES (
    'admin',
    'admin@kampus.ac.id',
    '$2b$12$rmQvGQ40Ni2tRNp7/tpjUe52CnKUITA0lmgTT2.elGl1uYPx8At6y',
    'Administrator Sistem',
    (SELECT id FROM roles WHERE name = 'Admin'),
    TRUE
);

-- USER DOSEN
INSERT INTO users (username, email, password_hash, full_name, role_id, is_active)
VALUES (
    'dosen',
    'dosen@kampus.ac.id',
    '$2b$12$7JGrGZqeckYXChzKff/WFu621K/DPVh/i4jC.flDlvu84iZfmQqri',
    'Dosen Wali',
    (SELECT id FROM roles WHERE name = 'Dosen Wali'),
    TRUE
);
