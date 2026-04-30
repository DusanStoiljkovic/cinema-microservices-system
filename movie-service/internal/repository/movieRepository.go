package repository

import (
	"context"
	"errors"
	"movie-service/internal/models"

	"gorm.io/gorm"
)

var (
	ErrGenreNotFound  = errors.New("one or more genres do not exits")
	ErrRecordNotFound = errors.New("movie does not exist")
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (repo *MovieRepository) GetMovies(
	ctx context.Context,
	limit, offset int,
	sort string,
	genre string,
	minYear, maxYear int,
	minRating float64,
) ([]*models.Movie, error) {
	var movies []*models.Movie

	query := repo.db.WithContext(ctx).Model(&models.Movie{}).Preload("Genres")

	if genre != "" {
		query = query.
			Joins("JOIN movie_genres mg ON mg.movie_id = movies.id").
			Joins("JOIN genres g ON g.id = mg.genre_id").
			Where("LOWER(g.name) = LOWER(?)", genre)
	}

	if minYear != 0 && maxYear != 0 {
		query = query.Where("year BETWEEN ? AND ?", minYear, maxYear)
	} else if minYear != 0 {
		query = query.Where("year >= ?", minYear)
	} else if maxYear != 0 {
		query = query.Where("year <= ?", maxYear)
	}

	if minRating != 0 {
		query = query.Where("rating >= ?", minRating)
	}

	if sort != "" {
		query = query.Order(sort)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&movies).Error
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (repo *MovieRepository) GetMovieByID(ctx context.Context, id uint) (*models.Movie, error) {
	var movie models.Movie

	err := repo.db.WithContext(ctx).Preload("Genres").First(&movie, id).Error
	if err != nil {
		return nil, err
	}

	return &movie, err
}

func (repo *MovieRepository) Create(ctx context.Context, movie *models.Movie) (*models.Movie, error) {
	genreIDs := extractGenreIDs(movie.Genres)

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		genres, err := repo.findGenresByIDs(tx, genreIDs)
		if err != nil {
			return err
		}

		movie.Genres = nil

		if err := tx.Create(movie).Error; err != nil {
			return err
		}

		if len(genres) > 0 {
			if err := tx.Model(movie).Association("Genres").Replace(genres); err != nil {
				return err
			}
		}

		return tx.Preload("Genres").First(movie, movie.ID).Error
	})

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (repo *MovieRepository) Update(ctx context.Context, id uint, updatedMovie *models.Movie) (*models.Movie, error) {
	var movie models.Movie

	genreIds := extractGenreIDs(updatedMovie.Genres)
	shouldUpdateGenres := updatedMovie.Genres != nil

	err := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Preload("Genres").First(&movie, id).Error; err != nil {
			return err
		}

		movie.Title = updatedMovie.Title
		movie.Description = updatedMovie.Description
		movie.Year = updatedMovie.Year
		movie.ImageURL = updatedMovie.ImageURL
		movie.Duration = updatedMovie.Duration
		movie.Rating = updatedMovie.Rating

		if err := tx.Save(&movie).Error; err != nil {
			return nil
		}

		if shouldUpdateGenres {
			genres, err := repo.findGenresByIDs(tx, genreIds)
			if err != nil {
				return err
			}

			if err := tx.Model(&movie).Association("Genres").Replace(genres); err != nil {
				return err
			}
		}

		return tx.Preload("Genres").First(&movie, id).Error
	})
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

func (repo *MovieRepository) Delete(ctx context.Context, id uint) error {
	return repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var movie models.Movie

		if err := tx.First(&movie, id).Error; err != nil {
			return err
		}

		if err := tx.Model(&movie).Association("Genres").Clear(); err != nil {
			return err
		}

		return tx.Delete(&movie).Error
	})
}

func (repo *MovieRepository) findGenresByIDs(tx *gorm.DB, genreIDs []uint) ([]models.Genre, error) {
	if len(genreIDs) == 0 {
		return []models.Genre{}, nil
	}

	var genres []models.Genre

	if err := tx.Where("id IN ?", genreIDs).Error; err != nil {
		return nil, err
	}

	if len(genres) != len(genreIDs) {
		return nil, ErrGenreNotFound
	}

	return genres, nil
}

func extractGenreIDs(genres []models.Genre) []uint {
	if genres == nil {
		return nil
	}

	ids := make([]uint, 0, len(genres))
	seen := make(map[uint]bool)

	for _, genre := range genres {
		if genre.ID == 0 {
			continue
		}

		if !seen[genre.ID] {
			seen[genre.ID] = true
			ids = append(ids, genre.ID)
		}
	}

	return ids
}
