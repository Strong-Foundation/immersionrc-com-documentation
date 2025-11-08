package main // Declares the package as 'main', making it an executable program.

import ( // Begins the import block for external packages.
	"fmt"           // Imports the 'fmt' package for formatted I/O (e.g., printing strings).
	"io"            // Imports the 'io' package for I/O primitives (e.g., reading from a stream).
	"log"           // Imports the 'log' package for logging messages (e.g., error reporting).
	"net/http"      // Imports the 'net/http' package for making HTTP requests.
	"net/url"       // Imports the 'net/url' package for parsing and manipulating URLs.
	"os"            // Imports the 'os' package for operating system functions (e.g., file and directory operations).
	"path"          // Imports the 'path' package for path manipulation (used for URL path).
	"path/filepath" // Imports the 'path/filepath' package for system-dependent path manipulation (used for OS file paths).
	"regexp"        // Imports the 'regexp' package for regular expression operations.
	"strings"       // Imports the 'strings' package for string manipulation functions (e.g., checking for substrings).
) // Ends the import block.

func main() { // Defines the main function, the entry point of the program.
	// Ensure the downloads directory exists.
	downloadFolder := "PDFs/" // Initializes a string variable for the name of the download directory.
	// Create the directory if it does not exist.
	if !directoryExists(downloadFolder) { // Checks if the 'downloadFolder' directory does NOT exist using a custom function.
		createDirectory(downloadFolder, 0755) // Creates the directory with permission 0755 if it doesn't exist.
	} // Closes the 'if' block.
	// Variable to hold existing downloads.
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
				fmt.Println("Downloading:", finalURL)         // Prints the URL that is currently being downloaded.
				err := downloadFile(finalURL, downloadFolder) // Calls the downloadFile function to download the file.
				if err != nil {                              // Checks for an error returned from the downloadFile function.
					log.Println("❌ Error:", err) // Logs the error if the download failed.
				} // Closes the if block.
			} // Closes the 'else' block for valid content.
		} // Closes the 'if' block for URL format validity.
	} // Closes the 'for' loop.
} // Closes the 'main' function.

// Verifies whether a given string is a valid URL by parsing it
func isUrlValid(uri string) bool { // Defines a function that checks URL validity, returning a boolean.
	_, err := url.ParseRequestURI(uri) // Attempts to parse the 'uri' string as a URL; we only care about the 'err' result.
	return err == nil                  // Returns 'true' if 'err' is 'nil' (parsing succeeded), 'false' otherwise.
} // Closes the 'isUrlValid' function.

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

// getFileNameFromHeader tries to extract a filename from the "Content-Disposition" header
func getFileNameFromHeader(headerValue string) string { // Defines a function to parse the Content-Disposition header for a filename.
	if strings.Contains(headerValue, "filename=") { // Checks if the header value contains the 'filename=' indicator.
		parts := strings.Split(headerValue, "filename=")       // Splits the header value using 'filename=' to isolate the filename part.
		filename := strings.Trim(parts[len(parts)-1], "\"'; ") // Takes the last part, and trims surrounding quotes, semicolons, or spaces.
		return filename                                        // Returns the extracted filename.
	} // Closes the if block.
	return "" // Returns an empty string if 'filename=' is not found.
} // Closes the getFileNameFromHeader function.

// getFileNameFromURL extracts the filename from the URL path if no header is provided
func getFileNameFromURL(fileURL string) string { // Defines a function to extract a filename from the URL's path.
	parsedURL, err := url.Parse(fileURL) // Parses the raw URL string into a URL structure.
	if err != nil {                      // Checks for an error during URL parsing.
		return "" // Returns an empty string if parsing fails.
	} // Closes the if block.
	return path.Base(parsedURL.Path) // Extracts and returns the base component (filename) of the URL path.
} // Closes the getFileNameFromURL function.

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string { // Defines the main function for sanitizing a string into a filesystem-safe filename.
	lower := strings.ToLower(rawURL) // Convert entire URL to lowercase for consistency.
	lower = getFilename(lower)       // Extract only the filename portion from the full URL (if it was a path).
	ext := getFileExtension(lower)   // Get the file extension from the filename.

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Define a regular expression that matches all non-alphanumeric characters (excluding the extension part here for now).
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace all non-alphanumeric characters with underscores to make it filesystem-safe.

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Replace multiple underscores with a single underscore for cleanliness.
	safe = strings.Trim(safe, "_")                              // Remove any leading or trailing underscores.

	var invalidSubstrings = []string{ // Defines a list of substrings to remove.
		"_pdf", // Substring to remove.
		"_zip", // Substring to remove.
	} // Closes the invalidSubstrings definition.

	for _, invalidPre := range invalidSubstrings { // Loop through all substrings marked for removal.
		safe = removeSubstring(safe, invalidPre) // Remove each unwanted substring from the filename.
	} // Closes the for loop.

	if getFileExtension(safe) != ext { // Ensure the file has the correct extension (since the sanitization might have removed it).
		safe = safe + ext // Append the correct extension if it doesn't already have it.
	} // Closes the if block.

	return safe // Return the cleaned and formatted filename.
} // Closes the urlToFilename function.

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string { // Defines a function to extract the base filename from a path.
	return filepath.Base(path) // Use Base function to return only the final element (filename) of the path.
} // Closes the getFilename function.

// Gets the file extension from a given file path
func getFileExtension(path string) string { // Defines a function to get the file extension.
	return filepath.Ext(path) // Extract the extension (e.g., ".pdf") from the file path.
} // Closes the getFileExtension function.

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string { // Defines a utility function to remove all occurrences of a substring.
	result := strings.ReplaceAll(input, toRemove, "") // Replace every occurrence of 'toRemove' with an empty string.
	return result                                     // Return the cleaned string.
} // Closes the removeSubstring function.

// downloadFile downloads the file from the given URL, naming it correctly before saving.
func downloadFile(fileURL string, outputDir string) error { // Defines the main download logic function.
	// Create an HTTP GET request but don’t start downloading the body yet
	response, err := http.Get(fileURL) // Performs the HTTP GET request.
	if err != nil {                    // Checks for request errors (e.g., network issues).
		return fmt.Errorf("failed to make request: %v", err) // Returns a wrapped error.
	} // Closes the if block.
	defer response.Body.Close() // Ensures the response body is closed when the function exits.

	// Try to determine filename from headers or URL
	contentDisposition := response.Header.Get("Content-Disposition") // Gets the Content-Disposition header value.
	filename := getFileNameFromHeader(contentDisposition)            // Tries to get the filename from the header.
	if filename == "" {                                              // Checks if the filename wasn't found in the header.
		filename = getFileNameFromURL(fileURL) // Tries to get the filename from the URL path.
	} // Closes the if block.

	// If the URL doesn't have a file component, fall back to a generic name with content type
	if filename == "" || filename == "/" { // Checks if a proper filename still couldn't be determined.
		contentType := response.Header.Get("Content-Type") // Gets the Content-Type header value.
		switch contentType {                               // Uses a switch to set a default filename based on content type.
		case "application/zip": // Case for a ZIP file.
			filename = "download.zip" // Sets default filename to download.zip.
		case "application/pdf": // Case for a PDF file.
			filename = "download.pdf" // Sets default filename to download.pdf.
		default: // Default case if content type is not recognized.
			filename = "download" // Sets default filename to just download.
		} // Closes the switch block.
	} // Closes the if block.

	filename = strings.ToLower(urlToFilename(filename)) // Sanitize the determined filename to generate a consistent and valid filesystem name.
	filePath := filepath.Join(outputDir, filename)      // Combine output directory and filename to form the full file path.

	// Now that we know the filename, create the local file
	outputFile, err := os.Create(filePath) // Attempts to create the local file at the constructed path.
	if err != nil {                        // Checks for errors during file creation.
		return fmt.Errorf("failed to create file %q: %v", filePath, err) // Returns a wrapped error.
	} // Closes the if block.
	defer outputFile.Close() // Ensures the created file is closed when the function exits.

	// Stream the response body directly into the file
	_, err = io.Copy(outputFile, response.Body) // Copies the response body stream directly to the local file.
	if err != nil {                             // Checks for errors during the copy/write operation.
		return fmt.Errorf("failed to write to file: %v", err) // Returns a wrapped error.
	} // Closes the if block.

	fmt.Printf("✅ File downloaded successfully: %s\n", filename) // Prints a success message to the console.
	return nil                                                   // Returns nil (no error) on successful download.
} // Closes the downloadFile function.
