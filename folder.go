package folder

import "io"

// Folder allows reading messages one by one.
type Folder interface {
	// Message returns an io.Reader for the next message in the folder,
	// and a close function to call when reading is complete. After the
	// last message, Message returns a nil reader. After calling
	// Message, it is an error to call it again before first calling the
	// close function.
	Message() (io.Reader, func() error, error)
}
