#!/usr/bin/env bash

# findreplace.sh - A simple find and replace utility for text files
# Usage: ./findreplace.sh [options] "search pattern" "replacement" files...

# Display help information
show_help() {
    echo "Usage: ./findreplace.sh [options] \"search pattern\" \"replacement\" files..."
    echo ""
    echo "Options:"
    echo "  -h, --help       Display this help message"
    echo "  -i, --ignore-case Perform case-insensitive search"
    echo "  -b, --backup     Create a backup of original files (.bak extension)"
    echo "  -r, --recursive  Recursively process directories"
    echo "  -d, --dry-run    Show what would be changed without making changes"
    echo ""
    echo "Examples:"
    echo "  ./findreplace.sh \"foo\" \"bar\" *.txt         # Replace 'foo' with 'bar' in all .txt files"
    echo "  ./findreplace.sh -i \"ERROR\" \"WARNING\" log.txt  # Case-insensitive replace"
    echo "  ./findreplace.sh -rb \"old\" \"new\" .            # Recursive with backup"
    exit 0
}

# Initialize variables
IGNORE_CASE=0
BACKUP=0
RECURSIVE=0
DRY_RUN=0

# Parse options
while [[ "$1" == -* ]]; do
    case "$1" in
        -h|--help)
            show_help
            ;;
        -i|--ignore-case)
            IGNORE_CASE=1
            shift
            ;;
        -b|--backup)
            BACKUP=1
            shift
            ;;
        -r|--recursive)
            RECURSIVE=1
            shift
            ;;
        -d|--dry-run)
            DRY_RUN=1
            shift
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            ;;
    esac
done

# Check if required arguments are provided
if [ $# -lt 3 ]; then
    echo "Error: Missing required arguments"
    show_help
fi

# Get search and replace patterns
SEARCH_PATTERN="$1"
REPLACEMENT="$2"
shift 2

# Build sed command
SED_CMD="sed"
if [ "$IGNORE_CASE" -eq 1 ]; then
    SED_CMD="$SED_CMD -i"
fi

if [ "$BACKUP" -eq 1 ]; then
    SED_CMD="$SED_CMD -i.bak"
else
    SED_CMD="$SED_CMD -i"
fi

# Function to process a single file
process_file() {
    local file="$1"
    
    # Skip if not a regular file
    if [ ! -f "$file" ]; then
        return
    fi
    
    # Skip binary files
    if file "$file" | grep -q "binary"; then
        echo "Skipping binary file: $file"
        return
    fi
    
    if [ "$DRY_RUN" -eq 1 ]; then
        # Show what would be changed
        echo "Would process: $file"
        if [ "$IGNORE_CASE" -eq 1 ]; then
            grep -i "$SEARCH_PATTERN" "$file" | while read -r line; do
                echo "  $line"
            done
        else
            grep "$SEARCH_PATTERN" "$file" | while read -r line; do
                echo "  $line"
            done
        fi
    else
        # Perform the replacement
        if [ "$IGNORE_CASE" -eq 1 ]; then
            $SED_CMD "s/$SEARCH_PATTERN/$REPLACEMENT/gI" "$file"
        else
            $SED_CMD "s/$SEARCH_PATTERN/$REPLACEMENT/g" "$file"
        fi
        echo "Processed: $file"
    fi
}

# Process files and directories
for target in "$@"; do
    if [ -f "$target" ]; then
        # Process a single file
        process_file "$target"
    elif [ -d "$target" ] && [ "$RECURSIVE" -eq 1 ]; then
        # Process directory recursively
        if [ "$IGNORE_CASE" -eq 1 ]; then
            find "$target" -type f -print0 | xargs -0 grep -l -i "$SEARCH_PATTERN" | while read -r file; do
                process_file "$file"
            done
        else
            find "$target" -type f -print0 | xargs -0 grep -l "$SEARCH_PATTERN" | while read -r file; do
                process_file "$file"
            done
        fi
    elif [ -d "$target" ] && [ "$RECURSIVE" -eq 0 ]; then
        echo "Skipping directory: $target (use -r for recursive processing)"
    else
        echo "Warning: $target does not exist"
    fi
done

echo "Find and replace operation completed"
