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

	TagUpdateInput struct {
		UserID      uuid.UUID `json:"user_id" valid:"required"`
		ParentID    uuid.UUID `json:"parent_id" valid:"-"`
		Name        string    `json:"name" valid:"required"`
		ShowIncome  bool      `json:"show_income" valid:"-"`
		ShowOutcome bool      `json:"show_outcome" valid:"-"`
		Regular     bool      `json:"regular" valid:"-"`
	}
)

type CategoryCreateResponse struct {
	CategoryID uuid.UUID `json:"category_id"`
}
