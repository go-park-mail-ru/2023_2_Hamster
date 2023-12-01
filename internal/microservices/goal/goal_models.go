package goal

import "github.com/google/uuid"

type (
	GoalCreateRequest struct {
		UserId      uuid.UUID `json:"user_id" valid:"required"`
		Name        string    `json:"name" valid:"required"`
		Description string    `json:"description" valid:"-"`
		Total       float64   `json:"total" valid:"required,greaterzero"`
		Date        string    `json:"date" valid:"isdate"`
	}

	GoalUpdateRequest struct {
		ID          uuid.UUID `json:"id" valid:"required"`
		UserId      uuid.UUID `json:"user_id" valid:"required"`
		Name        string    `json:"name" valid:"required"`
		Description string    `json:"description" valid:"-"`
		Total       float64   `json:"total" valid:"required,greaterzero"`
		Date        string    `json:"date" valid:"isdate"`
	}

	GoalDeleteRequest struct {
		ID uuid.UUID `json:"id" valid:"required"`
	}
)
