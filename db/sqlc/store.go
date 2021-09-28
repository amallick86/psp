package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	TransferBankTx(ctx context.Context, arg TransferBankTxParams) (TransferBankTxResult, error)
	ReceiveBankTx(ctx context.Context, arg ReceiveBankTxParams) (ReceiveBankTxResult, error)
}

//SQLStore provides all functions to execute db queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

//NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

//execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err:%v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}

//TransferTxParams contains the input parameters of the transfer transaction
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

//TransferTXResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

//TransferTx performs a money transfer from one account to the other.
//It creates a transfer record, add account entries and update account balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return nil
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return nil
		}
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil

	})
	return result, err
}

//TransferTxParams contains the input parameters of the transfer transaction
type TransferBankTxParams struct {
	AccountID     int64 `json:"accountID"`
	BankAccountID int64 `json:"bankAccountID"`
	Amount        int64 `json:"amount"`
}

//TransferTXResult is the result of the transfer transaction
type TransferBankTxResult struct {
	Transfer    TransfersBank `json:"transfer"`
	FromAccount Account       `json:"from_account"`
}

//TransferTx performs a money transfer from one account to the other.
//It creates a transfer record, add account entries and update account balance within a single database transaction
func (store *SQLStore) TransferBankTx(ctx context.Context, arg TransferBankTxParams) (TransferBankTxResult, error) {
	var result TransferBankTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransferBank(ctx, CreateTransferBankParams{
			AccountID:     arg.AccountID,
			BankAccountID: arg.BankAccountID,
			Amount:        -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: -arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil

	})
	return result, err
}
func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	return
}

//TransferTxParams contains the input parameters of the transfer transaction
type ReceiveBankTxParams struct {
	AccountID     int64 `json:"account_id"`
	BankAccountID int64 `json:"bankAccount_id"`
	Amount        int64 `json:"amount"`
}

//TransferTXResult is the result of the transfer transaction
type ReceiveBankTxResult struct {
	Transfer  TransfersBank `json:"transfer"`
	ToAccount Account       `json:"to_account"`
}

//TransferTx performs a money transfer from one account to the other.
//It creates a transfer record, add account entries and update account balance within a single database transaction
func (store *SQLStore) ReceiveBankTx(ctx context.Context, arg ReceiveBankTxParams) (ReceiveBankTxResult, error) {
	var result ReceiveBankTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransferBank(ctx, CreateTransferBankParams{
			AccountID:     arg.AccountID,
			BankAccountID: arg.BankAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			ID:     arg.AccountID,
			Amount: arg.Amount,
		})
		if err != nil {
			return err
		}

		return nil

	})
	return result, err
}
