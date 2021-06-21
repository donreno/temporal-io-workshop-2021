package main

import (
	"log"

	"github.com/donreno/temporal-io-workshop-2021/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("Error al crear cliente", err)
	}

	defer c.Close()

	w := worker.New(c, "transfer-workflow", worker.Options{})

	w.RegisterWorkflow(workflow.TransferWorkflow)
	w.RegisterActivity(workflow.GetCustomerDetails)
	w.RegisterActivity(workflow.IsRiskyCustomer)
	w.RegisterActivity(workflow.ChargeAccount)
	w.RegisterActivity(workflow.PayToAccount)
	w.RegisterActivity(workflow.RevertCharge)
	w.RegisterActivity(workflow.RevertPayment)
	w.RegisterActivity(workflow.NotifyFailedTransfer)
	w.RegisterActivity(workflow.NotifySuccessfulTransfer)

	if err = w.Run(worker.InterruptCh()); err != nil {
		log.Fatalln("Error ejecutando worker", err)
	}
}
