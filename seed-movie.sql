USE movieDB;

SET FOREIGN_KEY_CHECKS = 0;
TRUNCATE TABLE movie_genres;
TRUNCATE TABLE movies;
TRUNCATE TABLE genres;
SET FOREIGN_KEY_CHECKS = 1;

INSERT INTO genres (id, name, created_at, updated_at)
VALUES
(1, 'action', NOW(), NOW()),
(2, 'crime', NOW(), NOW()),
(3, 'thriller', NOW(), NOW()),
(4, 'comedy', NOW(), NOW()),
(5, 'drama', NOW(), NOW()),
(6, 'horror', NOW(), NOW()),
(7, 'sci-fi', NOW(), NOW()),
(8, 'adventure', NOW(), NOW()),
(9, 'romance', NOW(), NOW()),
(10, 'animation', NOW(), NOW());

INSERT INTO movies (id, title, description, year, image_url, duration, rating, created_at, updated_at)
VALUES
(
    1,
    'American Gangster',
    'A crime drama about power, loyalty and the rise of an organized crime empire.',
    2007,
    'https://example.com/american-gangster.jpg',
    157,
    7.8,
    NOW(),
    NOW()
),
(
    2,
    'Interstellar',
    'A team of explorers travel through a wormhole in space in an attempt to save humanity.',
    2014,
    'https://example.com/interstellar.jpg',
    169,
    8.7,
    NOW(),
    NOW()
),
(
    3,
    'Inception',
    'A thief enters dreams to steal secrets and is offered a chance to erase his past.',
    2010,
    'https://example.com/inception.jpg',
    148,
    8.8,
    NOW(),
    NOW()
),
(
    4,
    'The Dark Knight',
    'Batman faces the Joker, a criminal mastermind who wants to create chaos in Gotham.',
    2008,
    'https://example.com/dark-knight.jpg',
    152,
    9.0,
    NOW(),
    NOW()
),
(
    5,
    'Dune',
    'A noble family becomes involved in a war for control over the galaxy’s most valuable resource.',
    2021,
    'https://example.com/dune.jpg',
    155,
    8.0,
    NOW(),
    NOW()
),
(
    6,
    'John Wick 4',
    'John Wick uncovers a path to defeating the High Table.',
    2023,
    'https://example.com/john-wick-4.jpg',
    169,
    7.7,
    NOW(),
    NOW()
),
(
    7,
    'Scary Movie',
    'A comedy parody of popular horror movies.',
    2000,
    'https://example.com/scary-movie.jpg',
    88,
    6.3,
    NOW(),
    NOW()
),
(
    8,
    'Inside Out',
    'A young girl’s emotions guide her through a major life change.',
    2015,
    'https://example.com/inside-out.jpg',
    95,
    8.1,
    NOW(),
    NOW()
);

INSERT INTO movie_genres (movie_id, genre_id)
VALUES
(1, 1),
(1, 2),
(1, 5),

(2, 5),
(2, 7),
(2, 8),

(3, 1),
(3, 3),
(3, 7),

(4, 1),
(4, 2),
(4, 3),

(5, 5),
(5, 7),
(5, 8),

(6, 1),
(6, 2),
(6, 3),

(7, 4),
(7, 6),

(8, 4),
(8, 8),
(8, 10);