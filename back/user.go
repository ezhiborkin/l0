package back

type User struct {
	ID       int64  `json:"-"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}
