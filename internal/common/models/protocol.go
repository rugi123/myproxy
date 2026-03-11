package models

type AuthRequest struct {
	Token string // Токен клиента
}

type AuthResponse struct {
	Success bool
	Message string // Описание ошибки/успеха
}
