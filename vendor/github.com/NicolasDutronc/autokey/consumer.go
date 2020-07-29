package autokey

// Consumer defines an object that is able to retrieve the key
type Consumer interface {
	Get() (string, error)
}
