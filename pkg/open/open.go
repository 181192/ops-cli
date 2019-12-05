/*
Package open opens a file, directory, or URI using the OS's default
application for that object type.  Optionally, you can
specify an application to use.
This is a proxy for the following commands:
					OSX: "open"
			Windows: "start"
	Linux/Other: "xdg-open"
*/
package open

// Run opens a file, directory, or URI using the OS's default
// application for that object type. Wait for the open
// command to complete.
func Run(input string) error {
	return open(input).Run()
}

// Start opens a file, directory, or URI using the OS's default application for that object type.
// Don't wait for the open command to complete.
func Start(input string) error {
	return open(input).Start()
}

// RunWith opens a file, directory, or URI using the specified application.
// Wait for the open command to complete.
func RunWith(input string, appName string) error {
	return openWith(input, appName).Run()
}

// StartWith opens a file, directory, or URI using the specified application.
// Don't wait for the open command to complete.
func StartWith(input string, appName string) error {
	return openWith(input, appName).Start()
}
