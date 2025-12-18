package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/gunjan/user-api/internal/models"
	"github.com/gunjan/user-api/internal/repository"
	"go.uber.org/zap"
)

type UserService struct {
	repo   *repository.UserRepository
	logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, logger *zap.Logger) *UserService {
	return &UserService{
		repo:   repo,
		logger: logger,
	}
}

// CalculateAge calculates the age from a date of birth
func CalculateAge(dob time.Time) int {
	now := time.Now()
	age := now.Year() - dob.Year()
	
	// Adjust if birthday hasn't occurred this year
	if now.YearDay() < dob.YearDay() {
		age--
	}
	
	return age
}

// CreateUser creates a new user
func (s *UserService) CreateUser(req models.UserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("failed to parse DOB", zap.Error(err))
		return nil, errors.New("invalid date format")
	}

	userID, err := s.repo.CreateUser(req.Name, dob)
	if err != nil {
		s.logger.Error("failed to create user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user created", zap.Int64("user_id", userID))

	return &models.UserResponse{
		ID:   userID,
		Name: req.Name,
		DOB:  req.DOB,
	}, nil
}

// GetUserByID retrieves a user by ID with calculated age
func (s *UserService) GetUserByID(id int64) (*models.UserResponse, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found", zap.Int64("user_id", id))
			return nil, errors.New("user not found")
		}
		s.logger.Error("failed to get user", zap.Error(err))
		return nil, err
	}

	age := CalculateAge(user.DOB)
	dobStr := user.DOB.Format("2006-01-02")

	return &models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  dobStr,
		Age:  &age,
	}, nil
}

// GetAllUsersSimple retrieves all users without pagination (matches requirements)
func (s *UserService) GetAllUsersSimple() ([]models.UserResponse, error) {
	users, _, err := s.repo.GetAllUsers(0, 10000) // Get all users
	if err != nil {
		s.logger.Error("failed to get users", zap.Error(err))
		return nil, err
	}

	userResponses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		age := CalculateAge(user.DOB)
		dobStr := user.DOB.Format("2006-01-02")
		userResponses = append(userResponses, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  dobStr,
			Age:  &age,
		})
	}

	return userResponses, nil
}

// GetAllUsers retrieves all users with pagination and calculated ages
func (s *UserService) GetAllUsers(page, pageSize int) (*models.UserListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize

	users, total, err := s.repo.GetAllUsers(offset, pageSize)
	if err != nil {
		s.logger.Error("failed to get users", zap.Error(err))
		return nil, err
	}

	userResponses := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		age := CalculateAge(user.DOB)
		dobStr := user.DOB.Format("2006-01-02")
		userResponses = append(userResponses, models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  dobStr,
			Age:  &age,
		})
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	return &models.UserListResponse{
		Users:      userResponses,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(id int64, req models.UserRequest) (*models.UserResponse, error) {
	dob, err := time.Parse("2006-01-02", req.DOB)
	if err != nil {
		s.logger.Error("failed to parse DOB", zap.Error(err))
		return nil, errors.New("invalid date format")
	}

	err = s.repo.UpdateUser(id, req.Name, dob)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found for update", zap.Int64("user_id", id))
			return nil, errors.New("user not found")
		}
		s.logger.Error("failed to update user", zap.Error(err))
		return nil, err
	}

	s.logger.Info("user updated", zap.Int64("user_id", id))

	return &models.UserResponse{
		ID:   id,
		Name: req.Name,
		DOB:  req.DOB,
	}, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(id int64) error {
	err := s.repo.DeleteUser(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn("user not found for deletion", zap.Int64("user_id", id))
			return errors.New("user not found")
		}
		s.logger.Error("failed to delete user", zap.Error(err))
		return err
	}

	s.logger.Info("user deleted", zap.Int64("user_id", id))
	return nil
}

