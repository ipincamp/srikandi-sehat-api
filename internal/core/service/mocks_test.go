package service_test

import (
	"context"
	"errors"
	"time"

	"github.com/ipincamp/srikandi-sehat/internal/core/domain"
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/token"
)

// This file contains mock implementations for all interfaces
// required by the services in this package, for unit testing.

// --- Hasher Mock ---

type MockHasher struct {
	HashFunc    func(password string) (string, error)
	CompareFunc func(hash string, password string) bool
}

func (m *MockHasher) Hash(password string) (string, error) {
	return m.HashFunc(password)
}
func (m *MockHasher) Compare(hash string, password string) bool {
	return m.CompareFunc(hash, password)
}

// --- TokenMaker Mock ---

type MockTokenMaker struct {
	CreateTokenFunc   func(userID, roleID, useFor string, duration time.Duration) (string, *token.Payload, error)
	ValidateTokenFunc func(tokenStr string) (*token.Payload, error)
}

func (m *MockTokenMaker) CreateToken(userID, roleID, useFor string, duration time.Duration) (string, *token.Payload, error) {
	return m.CreateTokenFunc(userID, roleID, useFor, duration)
}
func (m *MockTokenMaker) ValidateToken(tokenStr string) (*token.Payload, error) {
	return m.ValidateTokenFunc(tokenStr)
}

// --- MailService Mock ---

type MockMailService struct {
	SendFunc func(ctx context.Context, to, subject, plainBody, htmlBody string) error
}

func (m *MockMailService) Send(ctx context.Context, to, subject, plainBody, htmlBody string) error {
	return m.SendFunc(ctx, to, subject, plainBody, htmlBody)
}

// --- UserRepository Mock ---

type MockUserRepository struct {
	SaveFunc          func(ctx context.Context, user *domain.User) error
	FindByIDFunc      func(ctx context.Context, id string) (*domain.User, error)
	FindByEmailFunc   func(ctx context.Context, email string) (*domain.User, error)
	FindAllFunc       func(ctx context.Context) ([]*domain.User, error)
	UpdateFunc        func(ctx context.Context, user *domain.User) error
	DeleteFunc        func(ctx context.Context, id string) error
	FindByEmailCalled bool
}

func (m *MockUserRepository) Save(ctx context.Context, user *domain.User) error {
	// Add nil check for safety
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, user)
	}
	return errors.New("SaveFunc not stubbed")
}
func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	// Add nil check for safety
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, id)
	}
	return nil, errors.New("FindByIDFunc not stubbed")
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	// Add nil check for safety
	if m.FindByEmailFunc != nil {
		return m.FindByEmailFunc(ctx, email)
	}
	return nil, errors.New("FindByEmailFunc not stubbed")
}
func (m *MockUserRepository) FindAll(ctx context.Context) ([]*domain.User, error) {
	// return m.FindAllFunc(ctx)
	panic("not implemented")
}
func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	// This is the implemented mock function
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return errors.New("UpdateFunc not stubbed")
}
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	// return m.DeleteFunc(ctx)
	panic("not implemented")
}

// --- PersonalTokenRepository Mock ---

type MockPersonalTokenRepository struct {
	SaveFunc           func(ctx context.Context, token *domain.PersonalToken) error
	FindByIDFunc       func(ctx context.Context, jti string) (*domain.PersonalToken, error)
	DeleteFunc         func(ctx context.Context, jti string) error
	DeleteByUserIDFunc func(ctx context.Context, userID string) error
}

func (m *MockPersonalTokenRepository) Save(ctx context.Context, token *domain.PersonalToken) error {
	return m.SaveFunc(ctx, token)
}
func (m *MockPersonalTokenRepository) FindByID(ctx context.Context, jti string) (*domain.PersonalToken, error) {
	return m.FindByIDFunc(ctx, jti)
}
func (m *MockPersonalTokenRepository) Delete(ctx context.Context, jti string) error {
	return m.DeleteFunc(ctx, jti)
}
func (m *MockPersonalTokenRepository) DeleteByUserID(ctx context.Context, userID string) error {
	if m.DeleteByUserIDFunc != nil {
		return m.DeleteByUserIDFunc(ctx, userID)
	}
	return errors.New("DeleteByUserIDFunc not stubbed")
}

// --- UserTokenRepository Mock ---

type MockUserTokenRepository struct {
	SaveFunc                     func(ctx context.Context, token *domain.UserToken) error
	FindByUserIDAndPurposeFunc   func(ctx context.Context, userID string, purpose string) (*domain.UserToken, error)
	DeleteByUserIDAndPurposeFunc func(ctx context.Context, userID string, purpose string) error
}

func (m *MockUserTokenRepository) Save(ctx context.Context, token *domain.UserToken) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, token)
	}
	return errors.New("SaveFunc not stubbed")
}

func (m *MockUserTokenRepository) FindByUserIDAndPurpose(ctx context.Context, userID string, purpose string) (*domain.UserToken, error) {
	if m.FindByUserIDAndPurposeFunc != nil {
		return m.FindByUserIDAndPurposeFunc(ctx, userID, purpose)
	}
	return nil, errors.New("FindByUserIDAndPurposeFunc not stubbed")
}

func (m *MockUserTokenRepository) DeleteByUserIDAndPurpose(ctx context.Context, userID string, purpose string) error {
	if m.DeleteByUserIDAndPurposeFunc != nil {
		return m.DeleteByUserIDAndPurposeFunc(ctx, userID, purpose)
	}
	return errors.New("DeleteByUserIDAndPurposeFunc not stubbed")
}

// --- UnitOfWork & Transaction Mocks ---

// MockTransaction implements the Transaction port
type MockTransaction struct {
	// We embed the mock repos here so the transaction can return them
	MockUserRepo      ports.UserRepository
	MockTokenRepo     ports.PersonalTokenRepository
	MockUserTokenRepo ports.UserTokenRepository
	// We add functions to control Commit and Rollback
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *MockTransaction) GetUserRepository() ports.UserRepository {
	return m.MockUserRepo // Returns the mock repo
}
func (m *MockTransaction) GetPersonalTokenRepository() ports.PersonalTokenRepository {
	return m.MockTokenRepo // Returns the mock repo
}
func (m *MockTransaction) GetUserTokenRepository() ports.UserTokenRepository {
	return m.MockUserTokenRepo // Returns the mock repo
}
func (m *MockTransaction) Commit() error {
	return m.CommitFunc()
}
func (m *MockTransaction) Rollback() error {
	return m.RollbackFunc()
}

// MockUnitOfWork implements the UnitOfWork port
type MockUnitOfWork struct {
	BeginFunc func(ctx context.Context) (ports.Transaction, error)
}

func (m *MockUnitOfWork) Begin(ctx context.Context) (ports.Transaction, error) {
	return m.BeginFunc(ctx)
}
