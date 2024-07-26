#!/bin/bash

# Check if correct number of arguments are provided
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <old_app_name> <new_app_name>"
    exit 1
fi

OLD_NAME=$1
NEW_NAME=$2

# Function to replace text in files
replace_in_files() {
    local old_name=$1
    local new_name=$2
    find . -type f -exec bash -c '
        for file; do
            LC_CTYPE=C sed -i "" "s/${0}/${1}/g" "$file"
        done
    ' "$old_name" "$new_name" {} +
}

# Function to rename directories and files
rename_items() {
    local old_name=$1
    local new_name=$2
    find . -depth -name "*${old_name}*" -exec bash -c '
        for item; do
            newitem=$(dirname "$item")/"${item##*/}"
            newitem="${newitem//'"${old_name}"'/'"${new_name}"'}"
            if [ "$item" != "$newitem" ]; then
                mv "$item" "$newitem"
            fi
        done
    ' bash {} +
}

# Replace occurrences in file contents
replace_in_files "$OLD_NAME" "$NEW_NAME"

# Rename directories
rename_items "$OLD_NAME" "$NEW_NAME"

# Rename files
rename_items "$OLD_NAME" "$NEW_NAME"

echo "Renaming completed successfully."