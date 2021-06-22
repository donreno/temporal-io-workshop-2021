package workflow

import (
	"errors"
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/workflow"
)

type Transfer struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	Amount      int    `json:"amount"`
}

func TransferWorkflow(ctx workflow.Context, transfer Transfer) error {
	activityOptions := workflow.ActivityOptions{
		ScheduleToCloseTimeout: time.Minute,
		StartToCloseTimeout:    time.Second * 15,
	}

	ctx = workflow.WithActivityOptions(ctx, activityOptions)
	logger := workflow.GetLogger(ctx)

	if err := verifyCustomer(ctx, logger, transfer); err != nil {
		notifyFailedTransfer(ctx, transfer)
		return err
	}

	if err := executeTransfer(ctx, logger, transfer); err != nil {
		notifyFailedTransfer(ctx, transfer)
		return err
	}

	notifySuccessfulTransfer(ctx, transfer)

	return nil
}

func verifyCustomer(ctx workflow.Context, logger log.Logger, transfer Transfer) error {
	getCustomerInfoExec := workflow.ExecuteActivity(ctx, GetCustomerDetails, transfer.Origin)
	isCustomerRiskyExec := workflow.ExecuteActivity(ctx, IsRiskyCustomer, transfer.Origin)

	var customerName string
	err := getCustomerInfoExec.Get(ctx, &customerName)
	if err != nil {
		logger.Error("Error obteniendo informacion de cliente", transfer.Origin)
		return err
	}

	var isRisky bool
	err = isCustomerRiskyExec.Get(ctx, &isRisky)
	if err != nil {
		logger.Error("Error Resolviendo riesgo de cliente", transfer.Origin)
		return err
	}

	if isRisky {
		logger.Error("Cliente", transfer.Origin, "es riesgoso")
		return err
	}

	logger.Info("Cliente ", customerName, "Numero de cuenta", transfer.Origin, "no es riesgoso")

	return nil
}

func executeTransfer(ctx workflow.Context, logger log.Logger, transfer Transfer) error {
	chargeAccountExec := workflow.ExecuteActivity(ctx, ChargeAccount, transfer.Origin, transfer.Amount)
	payToAccountExec := workflow.ExecuteActivity(ctx, PayToAccount, transfer.Destination, transfer.Amount)

	chargeErr := chargeAccountExec.Get(ctx, nil)
	paymentError := payToAccountExec.Get(ctx, nil)

	if chargeErr != nil && paymentError != nil {
		logger.Error("Cargo fallido", chargeErr)
		logger.Error("Abono fallido", paymentError)

		return errors.New(chargeErr.Error() + " | " + paymentError.Error())
	} else if chargeErr != nil {
		logger.Error("Cargo fallido", chargeErr)
		workflow.ExecuteActivity(ctx, RevertPayment, transfer.Destination, transfer.Amount)
		return chargeErr
	} else if paymentError != nil {
		logger.Error("Abono fallido", paymentError)
		workflow.ExecuteActivity(ctx, RevertCharge, transfer.Origin, transfer.Amount)
		return paymentError
	}

	return nil
}

func notifyFailedTransfer(ctx workflow.Context, transfer Transfer) {
	defer workflow.ExecuteActivity(ctx, NotifyFailedTransfer, transfer.Origin, transfer.Destination, transfer.Amount).Get(ctx, nil)
}

func notifySuccessfulTransfer(ctx workflow.Context, transfer Transfer) {
	defer workflow.ExecuteActivity(ctx, NotifySuccessfulTransfer, transfer.Origin, transfer.Destination, transfer.Amount).Get(ctx, nil)
}
