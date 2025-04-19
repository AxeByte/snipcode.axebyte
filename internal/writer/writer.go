package writer

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// generateTree creates a string representation of the file tree based on the provided file paths.
// It mimics the output of the 'tree' command in an OS-agnostic way.
func generateTree(files []string) string {
	var treeBuilder strings.Builder
	treeBuilder.WriteString(".\n") // Start with the root directory representation

	// nodes maps directory paths to a sorted list of their children's base names.
	nodes := make(map[string][]string)
	// isDir tracks whether a given full path corresponds to a directory.
	isDir := make(map[string]bool)

	// Initialize root node
	nodes["."] = []string{}
	isDir["."] = true

	// Populate nodes and isDir maps from the flat list of file paths.
	for _, file := range files {
		// Normalize path separators for consistency
		file = filepath.Clean(file)
		parts := strings.Split(file, string(filepath.Separator))
		currentPath := "." // Start from root relative path

		for i, part := range parts {
			// Skip empty parts resulting from splitting (e.g., leading slash)
			if part == "" {
				continue
			}

			// Construct the full path of the child (directory or file)
			childPath := filepath.Join(currentPath, part)

			// Ensure the parent directory exists in the nodes map.
			if _, exists := nodes[currentPath]; !exists {
				nodes[currentPath] = []string{}
				isDir[currentPath] = true // If it wasn't known before but has children, it's a dir.
			}

			// Add the current part (child) to its parent's list if not already present.
			// This builds the parent-child relationships.
			found := false
			for _, existingChild := range nodes[currentPath] {
				if existingChild == part {
					found = true
					break
				}
			}
			if !found {
				nodes[currentPath] = append(nodes[currentPath], part)
			}

			// Mark the childPath as a directory if it's not the last part of the original file path.
			if i < len(parts)-1 {
				isDir[childPath] = true
			}

			// Move down to the next level for the next iteration.
			currentPath = childPath
		}
	}

	// Sort the children list for each directory node alphabetically.
	for k := range nodes {
		sort.Strings(nodes[k])
	}

	// Recursive function to actually build the tree string with prefixes.
	var buildLevel func(dirPath string, prefix string)
	buildLevel = func(dirPath string, prefix string) {
		// Get the sorted list of children (base names) for the current directory.
		children, ok := nodes[dirPath]
		if !ok || len(children) == 0 {
			return // No children or directory not explicitly listed (shouldn't happen for populated dirs)
		}

		for i, childBaseName := range children {
			isLast := i == len(children)-1                 // Check if this is the last child in the list.
			childFullPath := filepath.Join(dirPath, childBaseName) // Get the full path of the child.
			childIsDir := isDir[childFullPath]              // Check if the child is a directory.

			// Determine the correct prefix and connector based on position.
			connector := "├── "
			nextPrefix := prefix + "│   " // Prefix for children of this child (if it's a dir).
			if isLast {
				connector = "└── "         // Use different connector for the last child.
				nextPrefix = prefix + "    " // No vertical line needed in the prefix for children of the last child.
			}

			// Append the formatted line for this child to the tree string.
			treeBuilder.WriteString(prefix + connector + childBaseName + "\n")

			// If the child is a directory, recursively call buildLevel for its children.
			if childIsDir {
				buildLevel(childFullPath, nextPrefix)
			}
		}
	}

	// Start the recursive tree building process from the root directory ".".
	buildLevel(".", "")

	return treeBuilder.String()
}

// Write processes the list of files, formats them using a fixed heredoc identifier,
// optionally generates and appends a file tree, and writes the result to the output file.
// It logs progress to standard output.
func Write(out string, files []string, withTree bool) error {
	var b strings.Builder // Use a strings.Builder for efficient string concatenation.
	total := 0            // Keep track of the total characters written from file contents.

	// --- Write File Contents ---
	for _, f := range files {
		// Get the relative path from the current directory.
		rel, err := filepath.Rel(".", f)
		if err != nil {
			// Handle error getting relative path, though unlikely if collector worked.
			fmt.Fprintf(os.Stderr, "Error getting relative path for %s: %v\n", f, err)
			rel = f // Fallback to using the original path.
		}

		// Write the header and the start of the cat command with the fixed heredoc identifier.
		// Added extra newlines for better separation as per user request.
		b.WriteString(fmt.Sprintf("## %s\n\ncat <<SNIPCODE_HEREDOC\n", rel))

		// Read the content of the current file.
		data, err := os.ReadFile(f)
		if err != nil {
			// Log the error and return, failing the whole process if one file can't be read.
			fmt.Fprintf(os.Stderr, "Error reading file %s: %v\n", f, err)
			return fmt.Errorf("failed to read file %s: %w", f, err)
		}
		// Write the raw file content.
		b.Write(data)

		// Ensure the content block ends with a newline before the closing heredoc identifier.
		if len(data) > 0 && data[len(data)-1] != '\n' {
			b.WriteString("\n")
		}

		// Write the closing heredoc identifier and the separator.
		// Added extra newlines for better separation.
		b.WriteString("SNIPCODE_HEREDOC\n\n---\n\n")

		// Update total size and log progress for the included file.
		size := len(data)
		total += size
		fmt.Printf("➡️ Included %s (%d chars)\n", rel, size)
	}

	// --- Write File Tree ---
	if withTree {
		fmt.Println("🌳 Generating file tree...")
		b.WriteString("## File Tree\n\n") // Add header for the tree section.

		// Generate the tree string using the OS-agnostic function.
		// Assumes 'files' is sorted appropriately by the collector. If not, sort here:
		// sort.Strings(files)
		treeString := generateTree(files)
		b.WriteString(treeString)

		// Ensure the output ends with a newline if the tree is not empty.
		if len(treeString) > 0 && treeString[len(treeString)-1] != '\n' {
			b.WriteString("\n")
		}
		fmt.Println("🌳 Appended file tree")
	}

	// --- Write Output File ---
	// Write the entire accumulated string to the specified output file.
	if err := os.WriteFile(out, []byte(b.String()), 0644); err != nil {
		// Log error if writing the final output fails.
		fmt.Fprintf(os.Stderr, "Error writing output file %s: %v\n", out, err)
		return fmt.Errorf("failed to write output file %s: %w", out, err)
	}

	// Log the final summary.
	fmt.Printf("📥 Wrote %s: %d files, %d total content chars\n", out, len(files), total)
	return nil // Indicate success.
}
