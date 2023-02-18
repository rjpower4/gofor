package gofor

// Resource represents a remote file that can be downloaded to disk
type Resource struct {
	// FileName is the name of the file to save the given resource to when downloaded
	FileName string `toml:"filename"`

	// Description is a short description of the resource
	Description string `toml:"description" omitempty:"true"`

	// URL is the URL to retrieve the resource from
	URL string `toml:"url"`

	// Tags is the collection of helpful tags for the resource to allow searching
	Tags []string `toml:"tags" omitempty:"true"`
}
