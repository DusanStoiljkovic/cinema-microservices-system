USE userDB;

SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE users;
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO users (id, name, email, password, role, created_at)
VALUES
(
    1,
    'Admin User',
    'admin@cinema.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'admin',
    NOW()
),
(
    2,
    'Regular User',
    'user@cinema.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'user',
    NOW()
),
(
    3,
    'Dusan Stoiljkovic',
    'dusan@cinema.com',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
    'user',
    NOW()
);