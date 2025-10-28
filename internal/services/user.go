package services

import (
	"context"

	"golang-boilerplate/internal/cache"
	"golang-boilerplate/internal/dtos"
	"golang-boilerplate/internal/errors"
	"golang-boilerplate/internal/models"
	"golang-boilerplate/internal/repositories"

	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

type UserService interface {
	Create(ctx context.Context, req *dtos.CreateUserRequest) (*models.User, error)
	GetOneByID(ctx context.Context, userID string) (*models.User, error)
	Update(ctx context.Context, userID string, req *dtos.UpdateUserRequest) (*models.User, error)
	Delete(ctx context.Context, userID string) error
	List(ctx context.Context, pageableRequest *dtos.UserPageableRequest) (*dtos.DataResponse[models.User], error)
}

// UserService handles user business logic
type userService struct {
	userRepo    repositories.UserRepository
	companyRepo repositories.CompanyRepository
	cache       cache.Cache
}

// NewUserService creates a new user service
func ProvideUserService(
	userRepo repositories.UserRepository,
	companyRepo repositories.CompanyRepository,
	cache cache.Cache,
) UserService {
	return &userService{
		userRepo:    userRepo,
		companyRepo: companyRepo,
		cache:       cache,
	}
}

func (s *userService) Create(ctx context.Context, req *dtos.CreateUserRequest) (*models.User, error) {
	companies := []models.Company{}
	for _, company := range req.Companies {
		company, err := s.companyRepo.GetOneByID(company.ID)
		if err != nil {
			// Report to Sentry with context
			if hub := sentry.GetHubFromContext(ctx); hub != nil {
				hub.WithScope(func(scope *sentry.Scope) {
					scope.SetTag("service", "user_service")
					scope.SetTag("operation", "create_user")
					scope.SetExtra("company_id", company.ID)
					hub.CaptureException(err)
				})
			}

			log.WithFields(log.Fields{
				"company_id": company.ID,
			}).Errorf("Failed to get company for user creation: %v", err)

			return nil, errors.NotFoundError("Company", err).
				WithOperation("create_user").
				WithResource("company").
				WithContext("company_id", company.ID)
		}
		companies = append(companies, *company)
	}

	user := &models.User{
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		KeycloakID: req.KeycloakID,
		Companies:  companies,
	}

	user, err := s.userRepo.Create(user)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "create_user")
				scope.SetExtra("body_request", req)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"body_request": req,
		}).Errorf("Failed to create user: %v", err)

		return nil, errors.DatabaseError("Failed to create user", err).
			WithOperation("create_user").
			WithResource("user").
			WithContext("request", req)
	}

	return user, nil
}

func (s *userService) GetOneByID(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.GetOneByID(userID)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "get_user")
				scope.SetExtra("user_id", userID)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"user_id": userID,
		}).Errorf("Failed to get user: %v", err)

		return nil, errors.NotFoundError("User", err).
			WithOperation("get_user").
			WithResource("user").
			WithContext("user_id", userID)
	}

	return user, nil
}

func (s *userService) Update(ctx context.Context, userID string, req *dtos.UpdateUserRequest) (*models.User, error) {
	preloads := []string{"Companies"}
	user, err := s.userRepo.GetOneByID(userID, preloads...)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "update_user")
				scope.SetExtra("user_id", userID)
				scope.SetExtra("body_request", req)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"user_id":      userID,
			"body_request": req,
		}).Errorf("Failed to get user for update: %v", err)

		return nil, errors.NotFoundError("User", err).
			WithOperation("update_user").
			WithResource("user").
			WithContext("user_id", userID)
	}

	// Handle many-to-many relationship with companies
	// If req.Companies is provided (even if empty), we need to update the associations
	if len(req.Companies) > 0 {
		// Get requested company IDs
		requestedCompanyIDs := make(map[string]bool)
		for _, company := range req.Companies {
			if company.ID != "" {
				requestedCompanyIDs[company.ID] = true
			}
		}

		// Get current company IDs
		currentCompanyIDs := make(map[string]bool)
		for _, company := range user.Companies {
			currentCompanyIDs[company.ID.String()] = true
		}

		// Build the final companies list with only the requested companies
		var finalCompanies []models.Company
		for companyID := range requestedCompanyIDs {
			// Check if this company is already associated
			if currentCompanyIDs[companyID] {
				// Find the existing company in user.Companies
				for _, company := range user.Companies {
					if company.ID.String() == companyID {
						finalCompanies = append(finalCompanies, company)
						break
					}
				}
			} else {
				// This is a new company, fetch it from the database
				company, err := s.companyRepo.GetOneByID(companyID)
				if err != nil {
					// Report to Sentry with context
					if hub := sentry.GetHubFromContext(ctx); hub != nil {
						hub.WithScope(func(scope *sentry.Scope) {
							scope.SetTag("service", "user_service")
							scope.SetTag("operation", "update_user")
							scope.SetExtra("company_id", companyID)
							hub.CaptureException(err)
						})
					}

					log.WithFields(log.Fields{
						"company_id": companyID,
					}).Errorf("Failed to get company for user update: %v", err)

					return nil, errors.NotFoundError("Company", err).
						WithOperation("update_user").
						WithResource("company").
						WithContext("company_id", companyID)
				}
				finalCompanies = append(finalCompanies, *company)
			}
		}

		// Replace the user's companies with the final list
		// This will be an empty slice if no companies were requested
		user.Companies = finalCompanies
	}

	// Update user fields if provided in request
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.KeycloakID != "" {
		user.KeycloakID = req.KeycloakID
	}

	// Save the updated user
	err = s.userRepo.Update(user)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "update_user")
				scope.SetExtra("user_id", userID)
				scope.SetExtra("body_request", req)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"user_id":      userID,
			"body_request": req,
		}).Errorf("Failed to update user: %v", err)

		return nil, errors.DatabaseError("Failed to update user", err).
			WithOperation("update_user").
			WithResource("user").
			WithContext("user_id", userID)
	}

	return user, nil
}

func (s *userService) Delete(ctx context.Context, userID string) error {
	user, err := s.userRepo.GetOneByID(userID)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "delete_user")
				scope.SetExtra("user_id", userID)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"user_id": userID,
		}).Errorf("Failed to get user for delete: %v", err)

		return errors.NotFoundError("User", err).
			WithOperation("delete_user").
			WithResource("user").
			WithContext("user_id", userID)
	}

	err = s.userRepo.Delete(user)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "delete_user")
				scope.SetExtra("user_id", userID)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"user_id": userID,
		}).Errorf("Failed to delete user: %v", err)

		return errors.DatabaseError("Failed to delete user", err).
			WithOperation("delete_user").
			WithResource("user").
			WithContext("user_id", userID)
	}

	return nil
}

func (s *userService) List(ctx context.Context, pageableRequest *dtos.UserPageableRequest) (*dtos.DataResponse[models.User], error) {
	preloads := []string{"Companies"}
	users, err := s.userRepo.Get(pageableRequest, preloads...)
	if err != nil {
		// Report to Sentry with context
		if hub := sentry.GetHubFromContext(ctx); hub != nil {
			hub.WithScope(func(scope *sentry.Scope) {
				scope.SetTag("service", "user_service")
				scope.SetTag("operation", "get_users")
				scope.SetExtra("pageable_request", pageableRequest)
				hub.CaptureException(err)
			})
		}

		log.WithFields(log.Fields{
			"pageable_request": pageableRequest,
		}).Errorf("Failed to get users: %v", err)

		return nil, errors.DatabaseError("Failed to get users", err).
			WithOperation("get_users").
			WithResource("users").
			WithContext("pageable_request", pageableRequest)
	}

	return users, nil
}
