package types

type User struct {
	Name     string
	Email    string
	Password string
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:     name,
		Email:    email,
		Password: password,
	}
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		Name:  u.Name,
		Email: u.Email,
	}
}

type GasStation struct {
	Name     string
	Address  string
	Location Location
	Prices   PriceRecord
}

type History []HistoryRecord

type GasStationWithHistory struct {
	*GasStation
	History History
}

type PriceRecord map[string]float64

type HistoryRecord struct {
	Prices PriceRecord
	Time   int64
}

type Location struct {
	Lat float64
	Lon float64
}

type CreateUserRequest struct {
	Name     string
	Email    string
	Password string
}

type UserResponse struct {
	Name  string
	Email string
}
