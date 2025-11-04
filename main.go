package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	// Ensure the downloads directory exists.
	downloadFolder := "PDFs/"
	// Create the directory if it does not exist.
	if !directoryExists(downloadFolder) {
		createDirectory(downloadFolder, 0755)
	}
	// The local downloads file.
	localDownloadsFile := "downloads.txt"
	// Variable to hold existing downloads.
	var existingDownloads string
	// Read the local file if it exists.
	if fileExists(localDownloadsFile) {
		// Read the existing downloads.
		existingDownloads = readAFileAsString(localDownloadsFile)
	}
	// Base URL for downloads.
	url := "https://www.immersionrc.com/?download="
	// Loop though 0 to 10000.
	for index := 0; index <= 10000; index++ {
		// The final URL.
		finalURL := url + fmt.Sprint(index)
		// Check if there is a valid content at the URL.
		if isUrlValid(finalURL) {
			// Get the data from the URL.
			data := getDataFromURL(finalURL)
			// Check if data is not empty.
			if strings.Contains(string(data), "Invalid download.") {
				log.Println("Invalid:", finalURL)
			} else {
				if strings.Contains(existingDownloads, finalURL) {
					log.Println("Already exists in downloads file:", finalURL)
					continue
				}
				log.Println("Valid:", finalURL)
				// Append the data to a file.
				err := appendByteToFile("downloads.txt", []byte(finalURL+"\n"))
				if err != nil {
					log.Println("Error appending to file:", err)
				}
			}
		}
	}

}

// Appends the given data (byte slice) to a file; creates the file if it doesn’t exist
func appendByteToFile(filename string, data []byte) error { // Defines a function to append bytes to a file, returning an error if one occurs.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644) // Opens the file with flags: Append, Create (if not exist), Write-Only, and permissions 0644.
	if err != nil {                                                               // Checks if opening or creating the file failed.
		return err // Returns the error to the caller.
	} // Closes the 'if' block.
	defer file.Close() // Schedules the file to be closed when the function exits (even if an error occurs).

	_, err = file.Write(data) // Writes the 'data' byte slice to the opened file.
	return err                // Returns 'nil' on successful write, or the error if writing failed.
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
func readAFileAsString(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Println(err)
	}
	return string(content)
}

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
