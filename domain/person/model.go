package person

type Person struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Age      int    `json:"age"`
}
