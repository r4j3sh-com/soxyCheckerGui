#!/bin/bash

# Enhanced copyright header script with more file types and better handling

# Copyright header for Go files
GO_HEADER='/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */'

# Copyright header for JS/CSS/TS files
JS_HEADER='/*
 * SoxyChecker GUI - A powerful proxy checker application
 * Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
 * 
 * This software is licensed under the MIT License.
 * See the LICENSE file in the project root for full license information.
 */'

# Copyright header for HTML files
HTML_HEADER='<!--
  SoxyChecker GUI - A powerful proxy checker application
  Copyright (c) 2025 Rajesh Mondal (r4j3sh.com)
  
  This software is licensed under the MIT License.
  See the LICENSE file in the project root for full license information.
-->'

# Function to add header to file
add_header() {
    local file="$1"
    local header="$2"
    local temp_file="temp_$(basename "$file")"
    
    echo "Adding header to: $file"
    echo "$header" > "$temp_file"
    echo "" >> "$temp_file"
    cat "$file" >> "$temp_file"
    mv "$temp_file" "$file"
}

# Function to check if file already has copyright
has_copyright() {
    local file="$1"
    grep -q "Copyright (c)" "$file" 2>/dev/null
}

echo "Adding copyright headers to project files..."

# Add headers to Go files
echo "Processing Go files..."
find . -name "*.go" -not -path "./vendor/*" -not -path "./node_modules/*" -not -path "./.git/*" | while read -r file; do
    if ! has_copyright "$file"; then
        add_header "$file" "$GO_HEADER"
    else
        echo "Skipping $file (already has copyright)"
    fi
done

# Add headers to JS/JSX/TS/TSX files
echo "Processing JavaScript/TypeScript files..."
find . \( -name "*.js" -o -name "*.jsx" -o -name "*.ts" -o -name "*.tsx" \) \
    -not -path "./node_modules/*" -not -path "./.git/*" -not -path "./dist/*" -not -path "./build/*" | while read -r file; do
    if ! has_copyright "$file"; then
        add_header "$file" "$JS_HEADER"
    else
        echo "Skipping $file (already has copyright)"
    fi
done

echo "Copyright headers added successfully!"
