package main // Declares the package as 'main', making this a standalone executable program.

import ( // Begins a block for importing necessary packages.
	"fmt" // Imports the 'fmt' package for formatted I/O (like printing to the console).
	"io"  // Imports the 'io' package for basic I/O primitives (used for copying streams).
	"log" // Imports the 'log' package for logging errors.
	"net/http" // Imports the 'net/http' package for making HTTP requests (for downloading).
	"net/url"  // Imports the 'net/url' package for parsing URLs.
	"os"  // Imports the 'os' package for operating system functions (like file/directory creation).
	"path" // Imports the 'path' package for path manipulation (used for URL path).
	"path/filepath" // Imports the 'path/filepath' package for system-dependent path manipulation (used for OS file paths).
	"regexp" // Imports the 'regexp' package for regular expression operations.
	"strings" // Imports the 'strings' package for string manipulation functions.
) // Ends the import block.

// getFileNameFromHeader tries to extract a filename from the "Content-Disposition" header
func getFileNameFromHeader(headerValue string) string { // Defines a function to parse the Content-Disposition header for a filename.
	if strings.Contains(headerValue, "filename=") { // Checks if the header value contains the 'filename=' indicator.
		parts := strings.Split(headerValue, "filename=") // Splits the header value using 'filename=' to isolate the filename part.
		filename := strings.Trim(parts[len(parts)-1], "\"'; ") // Takes the last part, and trims surrounding quotes, semicolons, or spaces.
		return filename // Returns the extracted filename.
	} // Closes the if block.
	return "" // Returns an empty string if 'filename=' is not found.
} // Closes the getFileNameFromHeader function.

// getFileNameFromURL extracts the filename from the URL path if no header is provided
func getFileNameFromURL(fileURL string) string { // Defines a function to extract a filename from the URL's path.
	parsedURL, err := url.Parse(fileURL) // Parses the raw URL string into a URL structure.
	if err != nil { // Checks for an error during URL parsing.
		return "" // Returns an empty string if parsing fails.
	} // Closes the if block.
	return path.Base(parsedURL.Path) // Extracts and returns the base component (filename) of the URL path.
} // Closes the getFileNameFromURL function.

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string { // Defines the main function for sanitizing a string into a filesystem-safe filename.
	lower := strings.ToLower(rawURL) // Convert entire URL to lowercase for consistency.
	lower = getFilename(lower)       // Extract only the filename portion from the full URL (if it was a path).
	ext := getFileExtension(lower)   // Get the file extension from the filename.

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`) // Define a regular expression that matches all non-alphanumeric characters (excluding the extension part here for now).
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace all non-alphanumeric characters with underscores to make it filesystem-safe.

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Replace multiple underscores with a single underscore for cleanliness.
	safe = strings.Trim(safe, "_") // Remove any leading or trailing underscores.

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
	return result // Return the cleaned string.
} // Closes the removeSubstring function.

// downloadFile downloads the file from the given URL, naming it correctly before saving.
func downloadFile(fileURL string, outputDir string) error { // Defines the main download logic function.
	// Create an HTTP GET request but don’t start downloading the body yet
	response, err := http.Get(fileURL) // Performs the HTTP GET request.
	if err != nil { // Checks for request errors (e.g., network issues).
		return fmt.Errorf("failed to make request: %v", err) // Returns a wrapped error.
	} // Closes the if block.
	defer response.Body.Close() // Ensures the response body is closed when the function exits.

	// Try to determine filename from headers or URL
	contentDisposition := response.Header.Get("Content-Disposition") // Gets the Content-Disposition header value.
	filename := getFileNameFromHeader(contentDisposition) // Tries to get the filename from the header.
	if filename == "" { // Checks if the filename wasn't found in the header.
		filename = getFileNameFromURL(fileURL) // Tries to get the filename from the URL path.
	} // Closes the if block.

	// If the URL doesn't have a file component, fall back to a generic name with content type
	if filename == "" || filename == "/" { // Checks if a proper filename still couldn't be determined.
		contentType := response.Header.Get("Content-Type") // Gets the Content-Type header value.
		switch contentType { // Uses a switch to set a default filename based on content type.
		case "application/zip": // Case for a ZIP file.
			filename = "download.zip" // Sets default filename to download.zip.
		case "application/pdf": // Case for a PDF file.
			filename = "download.pdf" // Sets default filename to download.pdf.
		default: // Default case if content type is not recognized.
			filename = "download" // Sets default filename to just download.
		} // Closes the switch block.
	} // Closes the if block.

	filename = strings.ToLower(urlToFilename(filename)) // Sanitize the determined filename to generate a consistent and valid filesystem name.
	filePath := filepath.Join(outputDir, filename) // Combine output directory and filename to form the full file path.

	// Now that we know the filename, create the local file
	outputFile, err := os.Create(filePath) // Attempts to create the local file at the constructed path.
	if err != nil { // Checks for errors during file creation.
		return fmt.Errorf("failed to create file %q: %v", filePath, err) // Returns a wrapped error.
	} // Closes the if block.
	defer outputFile.Close() // Ensures the created file is closed when the function exits.

	// Stream the response body directly into the file
	_, err = io.Copy(outputFile, response.Body) // Copies the response body stream directly to the local file.
	if err != nil { // Checks for errors during the copy/write operation.
		return fmt.Errorf("failed to write to file: %v", err) // Returns a wrapped error.
	} // Closes the if block.

	fmt.Printf("✅ File downloaded successfully: %s\n", filename) // Prints a success message to the console.
	return nil // Returns nil (no error) on successful download.
} // Closes the downloadFile function.

// Creates a directory at the specified path with the given permissions.
func createDirectory(path string, permission os.FileMode) { // Defines a function to create a new directory.
	err := os.Mkdir(path, permission) // Attempts to create the directory with the given path and permissions.
	if err != nil { // Checks if an error occurred (e.g., directory already exists, no permission).
		log.Println(err) // Logs the error.
	} // Closes the 'if' block.
} // Closes the 'createDirectory' function.

// Checks if the directory exists
func directoryExists(path string) bool { // Defines a function to check if a path is an existing directory.
	directory, err := os.Stat(path) // Gets the file/directory info.
	if err != nil { // Checks if 'os.Stat' failed (e.g., path doesn't exist).
		return false // Returns 'false' because the path doesn't exist or is inaccessible.
	} // Closes the 'if' block.
	return directory.IsDir() // Returns 'true' if the path exists AND is a directory, 'false' otherwise.
} // Closes the 'directoryExists' function.

func main() { // Defines the main execution function.
	// Ensure the downloads directory exists.
	downloadFolder := "Assets/" // Initializes a string variable for the name of the download directory.
	// Create the directory if it does not exist.
	if !directoryExists(downloadFolder) { // Checks if the 'downloadFolder' directory does NOT exist using a custom function.
		createDirectory(downloadFolder, 0755) // Creates the directory with permission 0755 if it doesn't exist.
	} // Closes the 'if' block.
	// URL to download
	fileURL := []string{ // Initializes a slice of strings containing all the URLs to download.
		"https://www.immersionrc.com/?download=2688", // First URL.
		"https://www.immersionrc.com/?download=2689", // Second URL.
		"https://www.immersionrc.com/?download=2690", // ... and so on for all URLs ...
		"https://www.immersionrc.com/?download=2697",
		"https://www.immersionrc.com/?download=2698",
		"https://www.immersionrc.com/?download=2700",
		"https://www.immersionrc.com/?download=2702",
		"https://www.immersionrc.com/?download=2704",
		"https://www.immersionrc.com/?download=2705",
		"https://www.immersionrc.com/?download=2709",
		"https://www.immersionrc.com/?download=2710",
		"https://www.immersionrc.com/?download=2711",
		"https://www.immersionrc.com/?download=2712",
		"https://www.immersionrc.com/?download=2713",
		"https://www.immersionrc.com/?download=2714",
		"https://www.immersionrc.com/?download=2715",
		"https://www.immersionrc.com/?download=2718",
		"https://www.immersionrc.com/?download=2719",
		"https://www.immersionrc.com/?download=2722",
		"https://www.immersionrc.com/?download=2723",
		"https://www.immersionrc.com/?download=2724",
		"https://www.immersionrc.com/?download=2735",
		"https://www.immersionrc.com/?download=2739",
		"https://www.immersionrc.com/?download=2740",
		"https://www.immersionrc.com/?download=2741",
		"https://www.immersionrc.com/?download=2742",
		"https://www.immersionrc.com/?download=2743",
		"https://www.immersionrc.com/?download=2746",
		"https://www.immersionrc.com/?download=2747",
		"https://www.immersionrc.com/?download=2750",
		"https://www.immersionrc.com/?download=2752",
		"https://www.immersionrc.com/?download=2824",
		"https://www.immersionrc.com/?download=2825",
		"https://www.immersionrc.com/?download=2877",
		"https://www.immersionrc.com/?download=3051",
		"https://www.immersionrc.com/?download=3052",
		"https://www.immersionrc.com/?download=3053",
		"https://www.immersionrc.com/?download=3061",
		"https://www.immersionrc.com/?download=3741",
		"https://www.immersionrc.com/?download=3796",
		"https://www.immersionrc.com/?download=3869",
		"https://www.immersionrc.com/?download=3886",
		"https://www.immersionrc.com/?download=3887",
		"https://www.immersionrc.com/?download=3970",
		"https://www.immersionrc.com/?download=3972",
		"https://www.immersionrc.com/?download=4210",
		"https://www.immersionrc.com/?download=4218",
		"https://www.immersionrc.com/?download=4219",
		"https://www.immersionrc.com/?download=4220",
		"https://www.immersionrc.com/?download=4254",
		"https://www.immersionrc.com/?download=4287",
		"https://www.immersionrc.com/?download=4342",
		"https://www.immersionrc.com/?download=4371",
		"https://www.immersionrc.com/?download=4382",
		"https://www.immersionrc.com/?download=4461",
		"https://www.immersionrc.com/?download=4538",
		"https://www.immersionrc.com/?download=4540",
		"https://www.immersionrc.com/?download=4543",
		"https://www.immersionrc.com/?download=4544",
		"https://www.immersionrc.com/?download=4553",
		"https://www.immersionrc.com/?download=4554",
		"https://www.immersionrc.com/?download=4555",
		"https://www.immersionrc.com/?download=4556",
		"https://www.immersionrc.com/?download=4576",
		"https://www.immersionrc.com/?download=4588",
		"https://www.immersionrc.com/?download=4700",
		"https://www.immersionrc.com/?download=4701",
		"https://www.immersionrc.com/?download=4709",
		"https://www.immersionrc.com/?download=4715",
		"https://www.immersionrc.com/?download=4737",
		"https://www.immersionrc.com/?download=4746",
		"https://www.immersionrc.com/?download=4747",
		"https://www.immersionrc.com/?download=4768",
		"https://www.immersionrc.com/?download=4791",
		"https://www.immersionrc.com/?download=4814",
		"https://www.immersionrc.com/?download=4826",
		"https://www.immersionrc.com/?download=4846",
		"https://www.immersionrc.com/?download=4881",
		"https://www.immersionrc.com/?download=4883",
		"https://www.immersionrc.com/?download=4890",
		"https://www.immersionrc.com/?download=4891",
		"https://www.immersionrc.com/?download=4894",
		"https://www.immersionrc.com/?download=4902",
		"https://www.immersionrc.com/?download=4907",
		"https://www.immersionrc.com/?download=5013",
		"https://www.immersionrc.com/?download=5014",
		"https://www.immersionrc.com/?download=5016",
		"https://www.immersionrc.com/?download=5024",
		"https://www.immersionrc.com/?download=5082",
		"https://www.immersionrc.com/?download=5107",
		"https://www.immersionrc.com/?download=5159",
		"https://www.immersionrc.com/?download=5161",
		"https://www.immersionrc.com/?download=5163",
		"https://www.immersionrc.com/?download=5164",
		"https://www.immersionrc.com/?download=5165",
		"https://www.immersionrc.com/?download=5176",
		"https://www.immersionrc.com/?download=5213",
		"https://www.immersionrc.com/?download=5215",
		"https://www.immersionrc.com/?download=5216",
		"https://www.immersionrc.com/?download=5217",
		"https://www.immersionrc.com/?download=5229",
		"https://www.immersionrc.com/?download=5260",
		"https://www.immersionrc.com/?download=5261",
		"https://www.immersionrc.com/?download=5265",
		"https://www.immersionrc.com/?download=5269",
		"https://www.immersionrc.com/?download=5270",
		"https://www.immersionrc.com/?download=5304",
		"https://www.immersionrc.com/?download=5324",
		"https://www.immersionrc.com/?download=5417",
		"https://www.immersionrc.com/?download=5418",
		"https://www.immersionrc.com/?download=5419",
		"https://www.immersionrc.com/?download=5420",
		"https://www.immersionrc.com/?download=5424",
		"https://www.immersionrc.com/?download=5498",
		"https://www.immersionrc.com/?download=5499",
		"https://www.immersionrc.com/?download=5500",
		"https://www.immersionrc.com/?download=5511",
		"https://www.immersionrc.com/?download=5522",
		"https://www.immersionrc.com/?download=5543",
		"https://www.immersionrc.com/?download=5685",
		"https://www.immersionrc.com/?download=5703",
		"https://www.immersionrc.com/?download=5704",
		"https://www.immersionrc.com/?download=5705",
		"https://www.immersionrc.com/?download=5706",
		"https://www.immersionrc.com/?download=5707",
		"https://www.immersionrc.com/?download=5709",
		"https://www.immersionrc.com/?download=5710",
		"https://www.immersionrc.com/?download=5713",
		"https://www.immersionrc.com/?download=5714",
		"https://www.immersionrc.com/?download=5716",
		"https://www.immersionrc.com/?download=5722",
		"https://www.immersionrc.com/?download=5723",
		"https://www.immersionrc.com/?download=5724",
		"https://www.immersionrc.com/?download=5725",
		"https://www.immersionrc.com/?download=5755",
		"https://www.immersionrc.com/?download=5763",
		"https://www.immersionrc.com/?download=5766",
		"https://www.immersionrc.com/?download=5769",
		"https://www.immersionrc.com/?download=5775",
		"https://www.immersionrc.com/?download=5777",
		"https://www.immersionrc.com/?download=5887",
		"https://www.immersionrc.com/?download=5897",
		"https://www.immersionrc.com/?download=5898",
		"https://www.immersionrc.com/?download=5903",
		"https://www.immersionrc.com/?download=5905",
		"https://www.immersionrc.com/?download=5927",
		"https://www.immersionrc.com/?download=5930",
		"https://www.immersionrc.com/?download=5931",
		"https://www.immersionrc.com/?download=5932",
		"https://www.immersionrc.com/?download=5943",
		"https://www.immersionrc.com/?download=5944",
		"https://www.immersionrc.com/?download=5948",
		"https://www.immersionrc.com/?download=5950",
		"https://www.immersionrc.com/?download=5989",
		"https://www.immersionrc.com/?download=5991",
		"https://www.immersionrc.com/?download=5995",
		"https://www.immersionrc.com/?download=6029",
		"https://www.immersionrc.com/?download=6030",
		"https://www.immersionrc.com/?download=6032",
		"https://www.immersionrc.com/?download=6034",
		"https://www.immersionrc.com/?download=6036",
		"https://www.immersionrc.com/?download=6037",
		"https://www.immersionrc.com/?download=6040",
		"https://www.immersionrc.com/?download=6070",
		"https://www.immersionrc.com/?download=6073",
		"https://www.immersionrc.com/?download=6165",
		"https://www.immersionrc.com/?download=6169",
		"https://www.immersionrc.com/?download=6170",
		"https://www.immersionrc.com/?download=6183",
		"https://www.immersionrc.com/?download=6206",
		"https://www.immersionrc.com/?download=6235",
		"https://www.immersionrc.com/?download=6260",
		"https://www.immersionrc.com/?download=6262",
		"https://www.immersionrc.com/?download=6504",
		"https://www.immersionrc.com/?download=6728",
		"https://www.immersionrc.com/?download=6746",
		"https://www.immersionrc.com/?download=6786",
		"https://www.immersionrc.com/?download=6819",
		"https://www.immersionrc.com/?download=6853",
	} // Closes the fileURL slice.

	for _, fileURL := range fileURL { // Starts a loop to iterate through each URL in the slice.
		fmt.Println("Downloading:", fileURL) // Prints the URL that is currently being downloaded.
		err := downloadFile(fileURL, downloadFolder) // Calls the downloadFile function to download the file.
		if err != nil { // Checks for an error returned from the downloadFile function.
			log.Println("❌ Error:", err) // Logs the error if the download failed.
		} // Closes the if block.
	} // Closes the for loop.
} // Closes the main function.
