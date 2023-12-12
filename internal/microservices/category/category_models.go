package category

import "github.com/google/uuid"

// input output models
type (
	TagInput struct {
		UserId   uuid.UUID `json:"user_id"`
		ParentId uuid.UUID `json:"parent_id"`
		// Image       int       `json:"image_id"`
		Name        string `json:"name"`
		ShowIncome  bool   `json:"show_income"`
		ShowOutcome bool   `json:"show_outcome"`
		Regular     bool   `json:"regular"`
	}

	TagUpdateInput struct {
		UserID   uuid.UUID `json:"user_id" valid:"required"`
		ParentID uuid.UUID `json:"parent_id" valid:"-"`
		// Image       int       `json:"image_id"`
		Name        string `json:"name" valid:"required"`
		ShowIncome  bool   `json:"show_income" valid:"-"`
		ShowOutcome bool   `json:"show_outcome" valid:"-"`
		Regular     bool   `json:"regular" valid:"-"`
	}

	TagDeleteInput struct {
		ID uuid.UUID `json:"id" valid:"-"`
	}

	CategoryCreateResponse struct {
		CategoryID uuid.UUID `json:"category_id"`
	}
)

// category errors
