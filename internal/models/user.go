package models

// User represents a user domain entity
type User struct {
	BaseModel
	FirstName        string    `gorm:"column:first_name"`
	LastName         string    `gorm:"column:last_name"`
	Email            string    `gorm:"column:email"`
	KeycloakID       string    `gorm:"column:keycloak_id"`
	StripeCustomerID string    `gorm:"column:stripe_customer_id"`
	Companies        []Company `gorm:"many2many:user_companies;"`
}

// Manually set table name
func (User) TableName() string {
	return "users"
}
