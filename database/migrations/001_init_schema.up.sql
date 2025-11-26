-- ROLES
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);


-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    role_id UUID REFERENCES roles(id),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);


-- PERMISSIONS
CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    resource VARCHAR(50) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description TEXT
);

-- ROLE PERMISSIONS
CREATE TABLE role_permissions (
    role_id UUID REFERENCES roles(id),
    permission_id UUID REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

-- STUDENTS
CREATE TABLE students (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    student_id VARCHAR(20) UNIQUE NOT NULL,
    program_study VARCHAR(100),
    academic_year VARCHAR(10),
    advisor_id UUID,
    created_at TIMESTAMP DEFAULT NOW()
);

-- LECTURERS
CREATE TABLE lecturers (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    lecturer_id VARCHAR(20) UNIQUE NOT NULL,
    department VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

-- ACHIEVEMENT REFERENCES
CREATE TYPE achievement_status AS ENUM (
    'draft', 
    'submitted', 
    'verified', 
    'rejected'
);


CREATE TABLE achievement_references (
    id UUID PRIMARY KEY,
    student_id UUID REFERENCES students(id),
    mongo_achievement_id VARCHAR(24) NOT NULL,
    status achievement_status DEFAULT 'draft',
    submitted_at TIMESTAMP,
    verified_at TIMESTAMP,
    verified_by UUID REFERENCES users(id),
    rejection_note TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);