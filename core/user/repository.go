package user

type Repository interface {
	GetByPhone(phoneNumber PhoneNumber) (*User, error)
}
