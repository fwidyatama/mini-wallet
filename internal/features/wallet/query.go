package wallet

const (
	InsertWalletQuery = `
		INSERT INTO wallets (owned_by)
		VALUES ($1) RETURNING id,
        owned_by,
        enabled_at,
        balance,
        CASE when(status=1) THEN 'enabled'
        END AS is_enable `

	GetWalletByOwnerIdQuery = `
		SELECT id,
       	owned_by,
       	CASE when(status=1) THEN 'enabled' when(status=0) THEN 'disabled'
       	END AS status,
       	balance,
       	enabled_at,
       	disabled_at
		FROM wallets
		WHERE owned_by = $1  `

	DisableWalletQuery = `
		UPDATE wallets
		SET status = 0,
		    enabled_at=NULL,
		    disabled_at=now()
		WHERE owned_by=$1 RETURNING 
		id,
		owned_by,
		disabled_at,
		balance,
		CASE when(status=0) THEN 'disabled'
		END AS status`

	EnableWalletQuery = `
		UPDATE wallets
		SET status = 1,
		    enabled_at=now(),
		    disabled_at=NULL
		WHERE owned_by=$1 RETURNING id,
		owned_by,
		enabled_at,
		balance,
		CASE when(status=1) THEN 'enabled'
		END AS status`

	TransactionWalletQuery = `
		INSERT INTO histories (wallet_id, status, transaction_by, TYPE, amount, reference_id)
		VALUES ($1, $2, $3, $4, $5,$6) RETURNING
		id,
	    status,
	    transaction_by,
	    transaction_at,
	    amount,
	    reference_id
`
	UpdateBalanceWalletQuery = `
	UPDATE wallets
		SET balance = $1
		WHERE id=$2 
`
	GetWalletTransactionQuery = `
		SELECT *
		FROM histories
		WHERE transaction_by = $1
		ORDER BY transaction_at DESC
		LIMIT 10
`
)
