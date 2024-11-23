#!/bin/bash

# Script to gather all files with a specific extension in the current directory
# and its subdirectories, and output them to specific text files for each top-level folder.
# It requires a file extension as the first argument, and optionally a list of exclusions.

# Function to display usage information
show_usage() {
    echo "Usage: confiles <file-extension> [exclusions]"
    echo ""
    echo "Description: This script collects files with the specified extension and creates output files for each top-level folder that contains these files."
    echo ""
    echo "Arguments:"
    echo "  <file-extension> : The file extension to search for (e.g., .go, .py, .js)."
    echo "  [exclusions]     : (Optional) List of folders or files to exclude. Can be one or more."
    echo ""
    echo "Examples:"
    echo "  Single File Exclusion:"
    echo "    confiles .go 'main.go'"
    echo "    - This will find all .go files except for the 'main.go' file."
    echo ""
    echo "  Multiple Folder/File Exclusions:"
    echo "    confiles .go 'handlers' 'utils/test.go'"
    echo "    - This will exclude the entire 'handlers' folder and the specific 'utils/test.go' file."
    exit 1
}

# Check if the user provided a file extension
if [ -z "$1" ]; then
    show_usage
fi

# Get the current directory
current_dir=$(pwd)

# File extension to search for
file_ext="$1"

# Optional exclusions (provided as subsequent arguments)
shift
exclusions=("$@")

# Function to check if a folder or file is in the exclusions list
is_excluded() {
    local target="$1"
    for exclusion in "${exclusions[@]}"; do
        if [[ "$target" == *"$exclusion"* ]]; then
            return 0  # It is excluded
        fi
    done
    return 1  # Not excluded
}

# Traverse through top-level folders and process only folders that contain the specified file extension
find "$current_dir" -maxdepth 1 -type d | while read -r folder; do
    folder_name=$(basename "$folder")
    
    # Skip if the folder is in the exclusion list
    if is_excluded "$folder_name"; then
        echo "Skipping excluded folder: $folder_name"
        continue
    fi
    
    # Check if the folder or its subfolders contain any files with the given extension
    files=$(find "$folder" -type f -name "*${file_ext}")
    
    if [ -n "$files" ]; then
        # Create a unique output file for this folder
        output_file="${current_dir}/output_${folder_name}.txt"
        echo "Creating $output_file"
        
        # Process each file within the folder and its subfolders
        echo "$files" | while read -r file; do
            # Skip if the file is in the exclusion list
            if is_excluded "$file"; then
                echo "Skipping excluded file: $file"
                continue
            fi

            # Append the file path and content to the output file using grouped output redirection
            {
                printf "\n%s:\n\`\`\`%s\n" "$file" "${file_ext#*.}"
                cat "$file"
                printf "\n\`\`\`\n---\n"
            } >> "$output_file"
        done
    fi
done
