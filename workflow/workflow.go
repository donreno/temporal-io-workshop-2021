package workflow

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

type Transfer struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Amount      int    `json:"amount"`
}

func TransferWorkflow(ctx workflow.Context, transfer Transfer) error {
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Second * 30,
		StartToCloseTimeout:    time.Second * 30,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)

	getCustomerInfoExec := workflow.ExecuteActivity(ctx, GetCustomerDetails, transfer.Origin)
	isCustomerRiskyExec := workflow.ExecuteActivity(ctx, IsRiskyCustomer, transfer.Origin)

	var customerName string
	err := getCustomerInfoExec.Get(ctx, &customerName)
	if err != nil {
		notifyFailedTransfer(ctx, transfer)
		logger.Error("Error obteniendo informacion de cliente", transfer.Origin)
		return err
	}

	var isRisky bool
	err = isCustomerRiskyExec.Get(ctx, &isRisky)
	if err != nil {
		notifyFailedTransfer(ctx, transfer)
		logger.Error("Error Resolviendo riesgo de cliente", transfer.Origin)
		return err
	}

	if isRisky {
		notifyFailedTransfer(ctx, transfer)
		logger.Error("Cliente", transfer.Origin, "es riesgoso")
	}

	logger.Info("Cliente ", customerName, "Numero de cuenta", transfer.Origin, "no es riesgoso")

	return nil
}

func notifyFailedTransfer(ctx workflow.Context, transfer Transfer) {
	workflow.ExecuteActivity(ctx, NotifyFailedTransfer, transfer.Origin, transfer.Destination, transfer.Amount)
}
