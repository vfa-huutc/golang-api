package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vfa-khuongdv/golang-cms/internal/constants"
	"github.com/vfa-khuongdv/golang-cms/internal/models"
	"github.com/vfa-khuongdv/golang-cms/internal/services"
	"github.com/vfa-khuongdv/golang-cms/internal/utils"
	"github.com/vfa-khuongdv/golang-cms/pkg/errors"
	"github.com/vfa-khuongdv/golang-cms/pkg/logger"
)

type IUserhandler interface {
	CreateUser(c *gin.Context)
	ForgotPassword(c *gin.Context)
	ResetPassword(c *gin.Context)
	Login(c *gin.Context)
	GetUser(c *gin.Context)
	GetUsers(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetProfile(c *gin.Context)
	UpdateProfile(c *gin.Context)
}

type UserHandler struct {
	userService  services.IUserService
	redisService services.IRedisService
}

func NewUserHandler(userService services.IUserService, redisService services.IRedisService) *UserHandler {
	return &UserHandler{
		userService:  userService,
		redisService: redisService,
	}
}

func (handler *UserHandler) CreateUser(ctx *gin.Context) {

	var input struct {
		Email    string  `json:"email" binding:"required,email"`
		Password string  `json:"password" binding:"required,min=6,max=255"`
		Name     string  `json:"name" binding:"required,min=1,max=45"`
		Birthday *string `json:"birthday" binding:"required,datetime=2006-01-02"` // Assumes YYYY-MM-DD format
		Address  *string `json:"address" binding:"required,min=1,max=255"`
		Gender   int16   `json:"gender" binding:"required,oneof=0 1 2"`
	}

	// Bind and validate the JSON request body to the input struct
	// Return 400 Bad Request if validation fails
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Hash the password using the utils.HashPassword function
	// If hashing fails (returns empty string), return a 400 error
	hashpassword := utils.HashPassword(input.Password)
	if hashpassword == "" {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthPasswordHashFailed, "Hash password failed"),
		)
		return
	}

	// Create a new User model instance with the validated input data
	// Password is stored as the hashed value
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashpassword,
		Birthday: input.Birthday,
		Address:  input.Address,
		Gender:   input.Gender,
	}

	// Attempt to create the user in the database
	// Return 400 Bad Request if creation fails
	if err := handler.userService.CreateUser(&user); err != nil {
		utils.RespondWithError(ctx, http.StatusBadRequest, err)
		return
	}

	utils.RespondWithOK(ctx, http.StatusCreated, gin.H{"message": "Create user successfully"})
}

func (handle *UserHandler) ForgotPassword(ctx *gin.Context) {
	var input struct {
		Email string `json:"email" binding:"required,email"`
	}
	// Bind and validate JSON request body
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get user by email from database
	user, err := handle.userService.GetUserByEmail(input.Email)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Generate random token string for password reset
	newToken := utils.GenerateRandomString(60)

	expiredAt := time.Now().Add(time.Hour).Unix()

	// Set new token on user
	user.Token = &newToken
	user.ExpiredAt = &expiredAt

	// Update user in database with new token
	if err := handle.userService.UpdateUser(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Send password reset email to user
	if err := services.SendMailForgotPassword(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}
	logger.Info("Email sent successfully!")

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Forgot password successfully"})
}

func (handler *UserHandler) ResetPassword(ctx *gin.Context) {
	var input struct {
		Token       string `json:"token" binding:"required"`
		Password    string `json:"password" binding:"required,min=6,max=255"`
		NewPassword string `json:"new_password" binding:"required,min=6,max=255"`
	}
	// Bind and validate JSON request body
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get user by token from database
	user, err := handler.userService.GetUserByToken(input.Token)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Check if token is expired
	if time.Now().Unix() > *user.ExpiredAt {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthTokenExpired, "Token expired"),
		)
		return
	}

	// Check if new password is the same as old password
	if isValid := utils.CheckPasswordHash(input.Password, user.Password); !isValid {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthInvalidPassword, "Password is incorrect"),
		)
		return
	}

	// Hash the password using the utils.HashPassword function
	// If hashing fails (returns empty string), return a 400 error
	hashpassword := utils.HashPassword(input.Password)
	if hashpassword == "" {
		utils.RespondWithError(
			ctx,
			http.StatusInternalServerError,
			errors.New(errors.ErrAuthPasswordHashFailed, "Failed to hash password"),
		)
		return
	}

	// Update user password
	user.Password = hashpassword
	user.Token = nil
	user.ExpiredAt = nil

	// Update user in database
	if err := handler.userService.UpdateUser(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Reset password successfully"})
}

func (handler *UserHandler) ChangePassword(ctx *gin.Context) {
	// Get user ID from the context
	// If user ID is 0 or not found, return bad request error
	userId := ctx.GetUint("UserID")
	if userId == 0 {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, "Invalid UserID"),
		)
		return
	}

	var input struct {
		OldPassword     string `json:"old_password" binding:"required,min=6,max=255"`
		NewPassword     string `json:"new_password" binding:"required,min=6,max=255"`
		ConfirmPassword string `json:"confirm_password" binding:"required,min=6,max=255"`
	}
	// Bind and validate JSON request body
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get user by ID from database
	user, err := handler.userService.GetUser(uint(userId))
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Check if old password is correct
	if isValid := utils.CheckPasswordHash(input.OldPassword, user.Password); !isValid {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthPasswordMismatch, "Old password is incorrect"),
		)
		return
	}

	// Check if new password is the same as old password
	if input.OldPassword == input.NewPassword {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthPasswordMismatch, "New password must be different from old password"),
		)
		return
	}

	// Check if new password and confirm password match
	if input.NewPassword != input.ConfirmPassword {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrAuthInvalidPassword, "New password and confirm password do not match"),
		)
		return
	}

	// Hash the password using the utils.HashPassword function
	// If hashing fails (returns empty string), return a 400 error
	hashpassword := utils.HashPassword(input.NewPassword)
	if hashpassword == "" {
		utils.RespondWithError(
			ctx,
			http.StatusInternalServerError,
			errors.New(errors.ErrAuthPasswordHashFailed, "Hash password failed"),
		)
		return
	}

	// Update user password
	user.Password = hashpassword

	// Update user in database
	if err := handler.userService.UpdateUser(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Change password successfully"})
}

func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	// Get user ID from the context
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)

	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Get user from database
	item, err := handler.userService.GetUser(uint(userId))
	if item == nil {
		utils.RespondWithError(
			ctx,
			http.StatusNotFound,
			err,
		)
		return
	}

	// Delete user from database
	if err := handler.userService.DeleteUser(uint(userId)); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Delete user successfully"})
}

func (handler *UserHandler) UpdateUser(ctx *gin.Context) {
	// Get user ID from the context
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)

		return
	}

	// Define input struct with validation tags
	var input struct {
		Name     string `json:"name" validate:"min=1,max=45"`           // Name must be between 1-45 chars
		Birthday string `json:"birthday" validate:"valid_birthday"`     // Birthday must be valid date
		Address  string `json:"address" validate:"min=1,max=255"`       // Address must be between 1-255 chars
		Gender   int16  `json:"gender" validate:"required,oneof=0 1 2"` // Gender must be 0, 1 or 2
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, err.Error()),
		)
		return
	}

	// Get existing user from database
	user, err := handler.userService.GetUser(uint(userId))
	// Return error if user not found
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Update user fields with input values
	user.Name = input.Name
	user.Birthday = &input.Birthday
	user.Address = &input.Address
	user.Gender = input.Gender

	// Save updated user to database
	if err := handler.userService.UpdateUser(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Update user successfully"})
}

func (handler *UserHandler) GetUser(ctx *gin.Context) {
	// Get user ID from the context
	id := ctx.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, err.Error()),
		)
		return
	}

	// Get user from database
	user, err := handler.userService.GetUser(uint(userId))
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	utils.RespondWithOK(ctx, http.StatusOK, user)
}

func (handler *UserHandler) GetProfile(ctx *gin.Context) {
	// Get user ID from the context
	userId := ctx.GetUint("UserID")
	if userId == 0 {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, "Invalid UserID"),
		)
		return
	}

	var user models.User

	// Try to get user from Redis cache
	userString, err := handler.redisService.Get(constants.PROFILE)
	if err != nil {
		logger.Warnf("Failed to get user from Redis: %+v", err)
	}
	// If not in cache, get from DB
	if userString == "" {
		dbUser, err := handler.userService.GetUser(userId)
		if err != nil {
			utils.RespondWithError(
				ctx,
				http.StatusBadRequest,
				err,
			)
			return
		}
		user = *dbUser

		// Cache the user data
		if err := handler.cacheUserProfile(&user); err != nil {
			logger.Warnf("Failed to cache user profile: %v", err)
		}
	} else {
		logger.Info("User retrieved from Redis")
		if err := json.Unmarshal([]byte(userString), &user); err != nil {
			logger.Warnf("Failed to unmarshal user from Redis: %v", err)
		}

	}

	utils.RespondWithOK(ctx, http.StatusOK, user)
}

// cacheUserProfile serializes a user object to JSON and stores it in Redis cache
// Parameters:
//   - user: pointer to the user model to be cached
//
// Returns:
//   - error if JSON marshaling or Redis caching fails
func (handler *UserHandler) cacheUserProfile(user *models.User) error {
	// Convert user object to JSON bytes
	userJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	// Create Redis key by concatenating profile prefix with user ID
	profileKey := constants.PROFILE + string(rune(user.ID))

	// Store serialized user data in Redis
	if err := handler.redisService.Set(profileKey, userJSON); err != nil {
		return fmt.Errorf("failed to cache in Redis: %v", err)
	}
	return nil
}

func (handler *UserHandler) UpdateProfile(ctx *gin.Context) {
	// Get user ID from context and validate
	userId := ctx.GetUint("UserID")
	if userId == 0 {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidParse, "Invalid UserID"),
		)
		return
	}

	// Define input struct for profile update with validation rules
	var input struct {
		Name     *string `json:"name" binding:"omitempty,min=1,max=45"`            // Name must be between 1 and 45 characters if provided
		Birthday *string `json:"birthday" binding:"omitempty,datetime=2006-01-02"` // Birthday must be a valid date (YYYY-MM-DD) if provided
		Address  *string `json:"address" binding:"omitempty,min=1,max=255"`        // Address must be between 1 and 255 characters if provided
		Gender   *int16  `json:"gender" binding:"omitempty,oneof=0 1 2"`           // Gender must be 0, 1, or 2 if provided
	}

	// Bind and validate JSON request body
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			errors.New(errors.ErrInvalidData, "Invalid UserID"),
		)
		return
	}

	// Get existing user from database
	user, err := handler.userService.GetUser(userId)

	// Return error if user not found
	if err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}

	// Update user fields if provided in input
	if input.Name != nil {
		user.Name = *input.Name
	}
	if input.Birthday != nil {
		user.Birthday = input.Birthday
	}
	if input.Address != nil {
		user.Address = input.Address
	}
	if input.Gender != nil {
		user.Gender = *input.Gender
	}

	// Save updated user to database
	if err := handler.userService.UpdateUser(user); err != nil {
		utils.RespondWithError(
			ctx,
			http.StatusBadRequest,
			err,
		)
		return
	}
	// Clear cache
	profileKey := constants.PROFILE + string(rune(user.ID))
	if err := handler.redisService.Delete(profileKey); err != nil {
		logrus.Errorf("Failed to clear cache: %v", err)
	}

	utils.RespondWithOK(ctx, http.StatusOK, gin.H{"message": "Update profile successfully"})
}
