package services

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/user/chatbot-whatsapp/models"
)

type CustomerService struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewCustomerService(db *gorm.DB, log *zap.Logger) *CustomerService {
	return &CustomerService{db: db, log: log}
}

func (s *CustomerService) FindByPhone(phone string) (*models.Customer, error) {
	var customer models.Customer
	err := s.db.Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *CustomerService) FindOrCreate(phone string) (*models.Customer, error) {
	customer, err := s.FindByPhone(phone)
	if err == nil {
		return customer, nil
	}

	if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	customer = &models.Customer{
		Phone: phone,
	}

	if err := s.db.Create(customer).Error; err != nil {
		return nil, err
	}

	s.log.Info("new customer created", zap.String("phone", phone))
	return customer, nil
}

func (s *CustomerService) OptIn(phone string) (*models.Customer, error) {
	customer, err := s.FindOrCreate(phone)
	if err != nil {
		return nil, err
	}

	if customer.OptIn {
		return customer, nil
	}

	now := time.Now()
	if err := s.db.Model(customer).Updates(map[string]interface{}{
		"opt_in":   true,
		"opt_in_at": &now,
	}).Error; err != nil {
		return nil, err
	}

	customer.OptIn = true
	customer.OptInAt = &now

	s.log.Info("customer opted in", zap.String("phone", phone))
	return customer, nil
}

func (s *CustomerService) OptOut(phone string) error {
	return s.db.Model(&models.Customer{}).
		Where("phone = ?", phone).
		Update("opt_in", false).
		Error
}

func (s *CustomerService) UpdateName(phone, name string) error {
	return s.db.Model(&models.Customer{}).
		Where("phone = ?", phone).
		Update("name", name).
		Error
}

func (s *CustomerService) ListOptedIn(limit, offset int) ([]models.Customer, error) {
	var customers []models.Customer
	err := s.db.Where("opt_in = ?", true).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&customers).Error
	return customers, err
}

func (s *CustomerService) CountOptedIn() (int64, error) {
	var count int64
	err := s.db.Model(&models.Customer{}).
		Where("opt_in = ?", true).
		Count(&count).Error
	return count, err
}
