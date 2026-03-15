package ports

import (
	"github.com/labbs/nexo/application/action/dto"
)

type ActionPort interface {
	CreateAction(input dto.CreateActionInput) (*dto.CreateActionOutput, error)
	ListActions(input dto.ListActionsInput) (*dto.ListActionsOutput, error)
	GetAction(input dto.GetActionInput) (*dto.GetActionOutput, error)
	UpdateAction(input dto.UpdateActionInput) error
	DeleteAction(input dto.DeleteActionInput) error
	GetRuns(input dto.GetRunsInput) (*dto.GetRunsOutput, error)
	ExecuteActions(input dto.ExecuteActionInput)
}
