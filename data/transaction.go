package data

import "github.com/pkg/errors"

// WithTransaction is a helper function to easily wrap an action in a transaction.
func WithTransaction(db SqlxWrapper, action func(tx TxWrapper) error) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	if err := action(tx); err != nil {
		if txErr := tx.Rollback(); txErr != nil {
			return errors.Wrap(err, txErr.Error())
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
