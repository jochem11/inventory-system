package account

type Account struct {
	ID         string `json:"id"`
	FirstName  string `json:"firstName"`
	Insertion  string `json:"insertion"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	CardNumber string `json:"cardNumber"`
}
