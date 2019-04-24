package folder

import "io"

// Folder allows reading messages one by one.
type Folder interface {
	// Message returns an io.ReadCloser for the next message in the folder.
	// After the last message,
	// Message returns a nil reader.
	// After calling Message,
	// it is an error to call it again before first closing the ReadCloser.
	Message() (io.ReadCloser, error)
}
