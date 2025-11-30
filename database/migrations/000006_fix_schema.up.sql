-- ubah fk advisor
ALTER TABLE students
ADD CONSTRAINT fk_students_advisor
FOREIGN KEY (advisor_id) REFERENCES lecturers(id)
ON DELETE SET NULL;

-- add deleted status
ALTER TYPE achievement_status
ADD VALUE IF NOT EXISTS 'deleted';
