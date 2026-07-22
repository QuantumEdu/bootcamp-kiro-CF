package entities

// User represents a system user (admin or cashier).
type User struct {
	ID        int64  `json:"id"`
	Nombre    string `json:"nombre"`
	PinHash   string `json:"-"`
	Rol       string `json:"rol"`
	Activo    bool   `json:"activo"`
	CreatedAt string `json:"created_at"`
}
