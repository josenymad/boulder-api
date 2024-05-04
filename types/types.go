package types

type EnvVars struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
