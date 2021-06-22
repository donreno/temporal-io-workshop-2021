package main

import (
	"log"

	"github.com/donreno/temporal-io-workshop-2021/workflow"
	"github.com/gofiber/fiber/v2"
	"go.temporal.io/sdk/client"
)

func main() {
	app := fiber.New()

	workflowOpts := client.StartWorkflowOptions{
		ID:        "transfer-workflow",
		TaskQueue: "transfer-workflow-queue",
	}

	app.Post("/transfer", func(ctx *fiber.Ctx) error {
		var transfer workflow.Transfer
		ctx.BodyParser(&transfer)

		c, err := client.NewClient(client.Options{})
		if err != nil {
			log.Println("Error al crear cliente", err)
			return ctx.Status(500).SendString("Error creando client de temporal")
		}

		defer c.Close()

		exec, err := c.ExecuteWorkflow(ctx.Context(), workflowOpts, workflow.TransferWorkflow, transfer)
		if err != nil {
			log.Println("Error iniciando workflow", err)
			return ctx.Status(500).SendString("Error iniciando workflow")
		}

		log.Println("Workflow ID", exec.GetID(), "| Run ID", exec.GetRunID())

		if err = exec.Get(ctx.Context(), nil); err != nil {
			log.Println("Error obteniendo resultado de workflow", err)
			return ctx.Status(500).SendString("Error obteniendo resultado de workflow")
		}

		return ctx.Status(200).SendString("Transferencia realizada de forma exitosa!")
	})

	log.Fatal(app.Listen(":3000"))
}
