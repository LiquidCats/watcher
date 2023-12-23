package port

type TransactionExecutor interface {
	Transaction(cb func() error) error
}
