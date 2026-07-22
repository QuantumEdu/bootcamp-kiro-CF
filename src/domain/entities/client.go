package entities

// Client represents a customer.
type Client struct {
	ID        int64  `json:"id"`
	Nombre    string `json:"nombre"`
	Telefono  string `json:"telefono"`
	Direccion string `json:"direccion"`
	CreatedAt string `json:"created_at"`
}
