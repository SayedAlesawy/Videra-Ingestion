package errors

// IsError A function to check if an error has occurred
func IsError(err error) bool {
	return err != nil
}
