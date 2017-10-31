package lib

import (
    "fmt"
)

type ErrInvalidSelector string

func (e ErrInvalidSelector) Error() string {
	return fmt.Sprintf("The selector '%v' is invalid.", string(e))
}
