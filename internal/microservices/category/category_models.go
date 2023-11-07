package category

import "github.com/google/uuid"

type (
	TagInput struct {
		UserId      uuid.UUID `json:"user_id"`
		ParentId    uuid.UUID `json:"parent_id"`
		Name        string    `json:"name"`
		ShowIncome  bool      `json:"show_income"`
		ShowOutcome bool      `json:"show_outcome"`
		Regular     bool      `json:"regular"`
	}
)
