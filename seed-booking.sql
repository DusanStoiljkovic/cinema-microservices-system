USE bookingDB;

SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE tickets;
TRUNCATE TABLE orders;
TRUNCATE TABLE projections;
TRUNCATE TABLE halls;
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO halls (id, name, location, capacity)
VALUES
(1, 'Sala 1', 'Paracin', 120),
(2, 'Sala 2', 'Paracin', 80),
(3, 'Sala 3', 'Nis', 150),
(4, 'VIP Sala', 'Beograd', 50);

INSERT INTO projections (id, movie_id, hall_id, start_time, end_time, price)
VALUES
(
    1,
    1,
    1,
    '2026-05-10 18:00:00',
    '2026-05-10 20:40:00',
    500
),
(
    2,
    2,
    1,
    '2026-05-10 21:00:00',
    '2026-05-10 23:50:00',
    650
),
(
    3,
    3,
    2,
    '2026-05-11 19:00:00',
    '2026-05-11 21:30:00',
    600
),
(
    4,
    4,
    3,
    '2026-05-11 20:00:00',
    '2026-05-11 22:40:00',
    700
),
(
    5,
    5,
    4,
    '2026-05-12 18:30:00',
    '2026-05-12 21:10:00',
    900
),
(
    6,
    6,
    2,
    '2026-05-12 21:30:00',
    '2026-05-13 00:20:00',
    750
),
(
    7,
    7,
    3,
    '2026-05-13 17:00:00',
    '2026-05-13 18:30:00',
    450
),
(
    8,
    8,
    1,
    '2026-05-13 19:00:00',
    '2026-05-13 20:40:00',
    500
);

INSERT INTO orders (id, user_id, status, total_price)
VALUES
(1, 2, 'reserved', 1200),
(2, 3, 'paid', 1400),
(3, 2, 'cancelled', 900);

INSERT INTO tickets (id, projection_id, order_id, seat_number, created_at)
VALUES
(1, 3, 1, 15, NOW()),
(2, 3, 1, 16, NOW()),

(3, 4, 2, 21, NOW()),
(4, 4, 2, 22, NOW()),

(5, 5, 3, 5, NOW());