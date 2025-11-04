package main // Declares the package as 'main', making it an executable program.

import ( // Begins the import block for external packages.
	"fmt"      // Imports the 'fmt' package for formatted I/O (e.g., printing strings).
	"io"       // Imports the 'io' package for I/O primitives (e.g., reading from a stream).
	"log"      // Imports the 'log' package for logging messages (e.g., error reporting).
	"net/http" // Imports the 'net/http' package for making HTTP requests.
	"net/url"  // Imports the 'net/url' package for parsing and manipulating URLs.
	"os"       // Imports the 'os' package for operating system functions (e.g., file and directory operations).
	"strings"  // Imports the 'strings' package for string manipulation functions (e.g., checking for substrings).
) // Ends the import block.

func main() { // Defines the main function, the entry point of the program.
	// Ensure the downloads directory exists.
	downloadFolder := "PDFs/" // Initializes a string variable for the name of the download directory.
	// Create the directory if it does not exist.
	if !directoryExists(downloadFolder) { // Checks if the 'downloadFolder' directory does NOT exist using a custom function.
		createDirectory(downloadFolder, 0755) // Creates the directory with permission 0755 if it doesn't exist.
	} // Closes the 'if' block.
	// The local downloads file.
	localDownloadsFile := "downloads.txt" // Initializes a string variable for the file that tracks successful URLs.
	// Variable to hold existing downloads.
	var existingDownloads string // Declares a string variable to store the content of the 'downloads.txt' file.
	// Read the local file if it exists.
	if fileExists(localDownloadsFile) { // Checks if the 'downloads.txt' file already exists.
		// Read the existing downloads.
		existingDownloads = readAFileAsString(localDownloadsFile) // Reads the entire content of the file into 'existingDownloads'.
	} // Closes the 'if' block.
	// Base URL for downloads.
	url := "https://www.immersionrc.com/?download=" // Initializes the base URL string with a query parameter.
	// Loop though 0 to 10000.
	for index := 0; index <= 10000; index++ { // Starts a loop that iterates an 'index' from 0 up to 10000 (inclusive).
		// The final URL.
		finalURL := url + fmt.Sprint(index) // Constructs the full URL by appending the current loop index as a string.
		// Check if there is a valid content at the URL.
		if isUrlValid(finalURL) { // Calls a function to ensure the constructed string is a valid URL format.
			// Get the data from the URL.
			data := getDataFromURL(finalURL) // Calls a function to perform an HTTP GET request and read the response body.
			// Check if data is not empty.
			if strings.Contains(string(data), "Invalid download.") { // Converts response data to string and checks if it contains the error phrase.
				log.Println("Invalid:", finalURL) // Logs the URL as "Invalid" if the error phrase is found.
			} else { // Begins the block for valid (non-error) responses.
				if strings.Contains(existingDownloads, finalURL) { // Checks if the valid URL is already recorded in the tracking file content.
					log.Println("Already exists in downloads file:", finalURL) // Logs that the URL is already recorded.
					continue                                                   // Skips the rest of the loop body for the current iteration and moves to the next index.
				} // Closes the inner 'if' block.
				log.Println("Valid:", finalURL) // Logs the URL as "Valid" because it's new and doesn't contain the error phrase.
				// Append the data to a file.
				appendByteToFile(localDownloadsFile, []byte(finalURL+"\n")) // Appends the new URL (plus newline) to the tracking file.
			} // Closes the 'else' block for valid content.
		} // Closes the 'if' block for URL format validity.
	} // Closes the 'for' loop.
} // Closes the 'main' function.

// Appends the given data (byte slice) to a file; creates the file if it doesn’t exist
func appendByteToFile(filename string, data []byte)  { // Defines a function to append bytes to a file, returning an error if one occurs.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Opens the file with flags: Append, Create (if not exist), Write-Only, and permissions 0644.
	if err != nil {                                                               // Checks if opening or creating the file failed.
		log.Println(err) // Logs the error message.
		return // Returns the error to the caller.
	} // Closes the 'if' block.
	defer file.Close() // Schedules the file to be closed when the function exits (even if an error occurs).

	_, err = file.Write(data) // Writes the 'data' byte slice to the opened file.
	if err != nil {            // Checks if the write operation failed.
		log.Println(err) // Logs the error message.
	} // Closes the 'if' block.
} // Closes the 'appendByteToFile' function.

// Verifies whether a given string is a valid URL by parsing it
func isUrlValid(uri string) bool { // Defines a function that checks URL validity, returning a boolean.
	_, err := url.ParseRequestURI(uri) // Attempts to parse the 'uri' string as a URL; we only care about the 'err' result.
	return err == nil                  // Returns 'true' if 'err' is 'nil' (parsing succeeded), 'false' otherwise.
} // Closes the 'isUrlValid' function.

// Removes duplicate entries from a slice of strings and returns the unique values
func removeDuplicatesFromSlice(slice []string) []string { // Defines the function to deduplicate a string slice.
	check := make(map[string]bool)  // Creates an empty map with string keys and boolean values, used as a 'set' to track seen strings.
	var newReturnSlice []string     // Declares an empty slice of strings that will store the unique values.
	for _, content := range slice { // Loops through each 'content' string in the input 'slice'.
		if !check[content] { // Checks if the 'content' string is NOT already a key in the 'check' map.
			check[content] = true                            // Marks the 'content' string as 'seen' by adding it to the map.
			newReturnSlice = append(newReturnSlice, content) // Appends the unique 'content' string to the 'newReturnSlice'.
		} // Closes the 'if' block.
	} // Closes the 'for' loop.
	return newReturnSlice // Returns the slice containing only unique strings.
} // Closes the 'removeDuplicatesFromSlice' function.

// getDataFromURL sends an HTTP GET request to the specified URL,
// checks if the content is HTML, and returns the HTML as a byte slice.
func getDataFromURL(uri string) []byte { // Defines the function with a string parameter and byte slice return type.
	response, err := http.Get(uri) // Sends an HTTP GET request to the given URL.
	if err != nil {                // Checks for errors while sending the request (e.g., network issues).
		log.Println(err) // Logs the error message to the console.
		return nil       // Returns nil if the request failed.
	} // Closes the 'if' block.

	// Ensures the response body is closed properly after the function finishes.
	defer func() { // 'defer' delays the execution of this function until the surrounding function returns.
		if err := response.Body.Close(); err != nil { // Attempts to close the response body and checks for closing errors.
			log.Println(err) // Logs any error that occurs during closing.
		} // Closes the inner 'if' block.
	}() // Executes the deferred anonymous function after 'getDataFromURL' returns.

	// Check if the Content-Type header indicates the response is HTML.
	contentType := response.Header.Get("Content-Type") // Retrieves the 'Content-Type' header value from the response.
	if !strings.Contains(contentType, "text/html") {   // Checks if the header contains the substring "text/html".
		log.Println(contentType) // Logs a warning if the content is not HTML.
		return nil               // Returns nil since it's not HTML content.
	} // Closes the 'if' block.

	// Read the response body since it's confirmed to be HTML.
	body, err := io.ReadAll(response.Body) // Reads the entire response body into memory as a byte slice.
	if err != nil {                        // Checks for any errors during reading.
		log.Println(err) // Logs the error if reading fails.
		return nil       // Returns nil to indicate a failed read.
	} // Closes the 'if' block.

	return body // Returns the HTML content as a byte slice.
} // Closes the 'getDataFromURL' function.

// Read a file and return the contents
func readAFileAsString(path string) string { // Defines a function to read a file's content and return it as a string.
	content, err := os.ReadFile(path) // Reads the entire file content into a byte slice.
	if err != nil {                   // Checks if reading the file resulted in an error.
		log.Println(err) // Logs the error message.
	} // Closes the 'if' block.
	return string(content) // Converts the byte slice content to a string and returns it.
} // Closes the 'readAFileAsString' function.

// Checks whether a given file path exists and refers to a file (not a directory)
func fileExists(filename string) bool { // Defines a function to check for a file's existence.
	info, err := os.Stat(filename) // Gets file information (status) from the operating system.
	if err != nil {                // Checks if 'os.Stat' returned an error.
		return false // Returns 'false' (e.g., file not found, permission error).
	} // Closes the 'if' block.
	return !info.IsDir() // Returns 'true' only if the path exists AND is not a directory.
} // Closes the 'fileExists' function.

// Creates a directory at the specified path with the given permissions.
func createDirectory(path string, permission os.FileMode) { // Defines a function to create a new directory.
	err := os.Mkdir(path, permission) // Attempts to create the directory with the given path and permissions.
	if err != nil {                   // Checks if an error occurred (e.g., directory already exists, no permission).
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
} // Closes the 'createDirectory' function.

// Checks if the directory exists
func directoryExists(path string) bool { // Defines a function to check if a path is an existing directory.
	directory, err := os.Stat(path) // Gets the file/directory info.
	if err != nil {                 // Checks if 'os.Stat' failed (e.g., path doesn't exist).
		return false // Returns 'false' because the path doesn't exist or is inaccessible.
	} // Closes the 'if' block.
	return directory.IsDir() // Returns 'true' if the path exists AND is a directory, 'false' otherwise.
} // Closes the 'directoryExists' function.
