package category

import "github.com/google/uuid"

// input output models
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

<<<<<<< HEAD
	TagDeleteInput struct {
		ID uuid.UUID `json:"id" valid:"-"`
=======
	CategoryCreateResponse struct {
		CategoryID uuid.UUID `json:"category_id"`
>>>>>>> 3eec182f5e5b585c01ea07901fb0af5156153490
	}
)

// category errors
