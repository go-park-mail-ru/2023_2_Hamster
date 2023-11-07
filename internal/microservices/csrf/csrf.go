package csrf

import "github.com/google/uuid"

type Usecase interface {
	GenerateCSRFToken(userID uuid.UUID) (string, error)
	CheckCSRFToken(acessToken string) (uuid.UUID, error)
}
