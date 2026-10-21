package solana

import (
	"database/sql/driver"
	"fmt"
)

// liquid-collective patch: base58 text (de)serialization for database/sql, so PublicKey, Hash
// and Signature can be stored in text columns. Not part of upstream solana-go.

// scanText returns the text of a database value; drivers hand text columns over as string or []byte.
func scanText(value interface{}, typ string) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %s", value, typ)
	}
}

// Scan implements the sql.Scanner interface.
func (p *PublicKey) Scan(value interface{}) error {
	str, err := scanText(value, "PublicKey")
	if err != nil {
		return err
	}
	pubKey, err := PublicKeyFromBase58(str)
	if err != nil {
		return err
	}
	*p = pubKey
	return nil
}

// Value implements the driver.Valuer interface.
func (p PublicKey) Value() (driver.Value, error) {
	return p.String(), nil
}

// Scan implements the sql.Scanner interface.
func (h *Hash) Scan(value interface{}) error {
	str, err := scanText(value, "Hash")
	if err != nil {
		return err
	}
	hash, err := HashFromBase58(str)
	if err != nil {
		return err
	}
	*h = hash
	return nil
}

// Value implements the driver.Valuer interface.
func (h Hash) Value() (driver.Value, error) {
	return h.String(), nil
}

// Scan implements the sql.Scanner interface.
func (s *Signature) Scan(value interface{}) error {
	str, err := scanText(value, "Signature")
	if err != nil {
		return err
	}
	sig, err := SignatureFromBase58(str)
	if err != nil {
		return err
	}
	*s = sig
	return nil
}

// Value implements the driver.Valuer interface.
func (s Signature) Value() (driver.Value, error) {
	return s.String(), nil
}
