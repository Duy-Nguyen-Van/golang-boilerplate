package models

// User represents a user domain entity
type Company struct {
	BaseModel
	Name       string `gorm:"column:name"`
	KeycloakID string `gorm:"column:keycloak_id"`
	Users      []User `gorm:"many2many:user_companies;"`
}

// Manually set table name
func (Company) TableName() string {
	return "companies"
}
