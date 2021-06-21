package workflow

import (
	"log"
	"time"
)

func GetCustomerDetails(accountNumber string) (string, error) {
	time.Sleep(time.Millisecond * 20)
	log.Println("Cuenta identificada")
	return "Cliente 1", nil
}

func IsRiskyCustomer(accountNumber string) (bool, error) {
	time.Sleep(time.Millisecond * 100)
	log.Println("Cliente no es riesgoso")
	return false, nil
}

func ChargeAccount(accountNumber string, amount int) error {
	time.Sleep(time.Millisecond * 30)
	log.Println("Cargando", amount, "a cuenta", accountNumber)
	return nil
}

func PayToAccount(accountNumber string, amount int) error {
	time.Sleep(time.Millisecond * 30)
	log.Println("Abonando", amount, "a cuenta", accountNumber)
	return nil
}

func RevertCharge(accountNumber string, amount int) error {
	time.Sleep(time.Millisecond * 30)
	log.Println("Reversando cargo de", amount, "a cuenta", accountNumber)
	return nil
}

func RevertPayment(accountNumber string, amount int) error {
	time.Sleep(time.Millisecond * 30)
	log.Println("Reversando abono de", amount, "a cuenta", accountNumber)
	return nil
}

func NotifyFailedTransfer(origin, destination string, amount int) error {
	log.Println("Transaccion fallida de", amount, "desde", origin, "hacia", destination)
	return nil
}

func NotifySuccessfulTransfer(origin, destination string, amount int) error {
	log.Println("Transaccion exitosa de", amount, "desde", origin, "hacia", destination)
	return nil
}
