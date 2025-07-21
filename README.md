# Drive MCP

A Go implementation of MCP (Model Context Protocol) server for Google Drive, Google Docs, Google Slides, and Google Sheets operations.

## Features

- Search Google Drive files
- List files in Google Drive folders
- Read Google Document content
- Update Google Document content
- Read Google Slides presentation content
- Update Google Slides presentation slides
- Read Google Sheets values
- Update Google Sheets values
- Authentication using gcloud application-default credentials

## Setup

### Prerequisites

- Go 1.21 or later
- Google Cloud CLI (`gcloud`)
- GCP project with Google Drive API, Google Docs API, Google Slides API, and Google Sheets API enabled

### Authentication Setup

1. Enable Google Drive API, Google Docs API, Google Slides API, and Google Sheets API
    * https://console.cloud.google.com/apis/library/drive.googleapis.com
    * https://console.cloud.google.com/apis/library/docs.googleapis.com
    * https://console.cloud.google.com/apis/library/slides.googleapis.com
    * https://console.cloud.google.com/apis/library/sheets.googleapis.com
2. Run gcloud authentication:

```bash
gcloud auth application-default login --scopes=https://www.googleapis.com/auth/cloud-platform,https://www.googleapis.com/auth/drive
```

3. Set quota project environment variable if needed:

```bash
export GOOGLE_CLOUD_QUOTA_PROJECT_ID=your-project-id
```

### Installation

```bash
go mod download
```

## Usage

### Running the MCP Server

```bash
go build -o drive-mcp
./drive-mcp
```

### Available Tools

#### search_files

Search for files in Google Drive.

**Parameters:**
- `query` (required): File name or keyword to search
- `maxResults` (optional, default: 10): Maximum number of files to retrieve

**Example:**
```json
{
  "name": "search_files",
  "arguments": {
    "query": "meeting notes",
    "maxResults": 5
  }
}
```

#### list_files

List files in a Google Drive folder.

**Parameters:**
- `folderId` (optional): The ID of the folder to list files from. If empty, lists files in My Drive root
- `maxResults` (optional, default: 10): Maximum number of files to retrieve

**Example:**
```json
{
  "name": "list_files",
  "arguments": {
    "folderId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "maxResults": 20
  }
}
```

**Example (My Drive root):**
```json
{
  "name": "list_files",
  "arguments": {
    "maxResults": 10
  }
}
```

#### get_document

Get the content of a Google Document.

**Parameters:**
- `documentId` (required): The ID of the Google Document

**Example:**
```json
{
  "name": "get_document",
  "arguments": {
    "documentId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
  }
}
```

#### update_document

Update the content of a Google Document.

**Parameters:**
- `documentId` (required): The ID of the Google Document
- `content` (required): The new content for the document

**Example:**
```json
{
  "name": "update_document",
  "arguments": {
    "documentId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "content": "This is the new content for the document."
  }
}
```

#### get_presentation

Get the content of a Google Slides presentation.

**Parameters:**
- `presentationId` (required): The ID of the Google Slides presentation

**Example:**
```json
{
  "name": "get_presentation",
  "arguments": {
    "presentationId": "1EAYk18WDjIG-zp_0vLm3CsfQh_i8eXc67Jo2O9C6Vuc"
  }
}
```

#### update_presentation

Update a specific slide in a Google Slides presentation.

**Parameters:**
- `presentationId` (required): The ID of the Google Slides presentation
- `slideIndex` (optional, default: 0): The index of the slide to update (0-based)
- `title` (required): The title for the slide
- `content` (required): The content for the slide

**Example:**
```json
{
  "name": "update_presentation",
  "arguments": {
    "presentationId": "1EAYk18WDjIG-zp_0vLm3CsfQh_i8eXc67Jo2O9C6Vuc",
    "slideIndex": 0,
    "title": "New Slide Title",
    "content": "This is the new content for the slide."
  }
}
```

#### get_spreadsheet

Get values from a Google Spreadsheet.

**Parameters:**
- `spreadsheetId` (required): The ID of the Google Spreadsheet
- `range` (required): The range to retrieve (e.g., 'Sheet1!A1:C10')

**Example:**
```json
{
  "name": "get_spreadsheet",
  "arguments": {
    "spreadsheetId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "range": "Sheet1!A1:C10"
  }
}
```

#### update_spreadsheet

Update values in a Google Spreadsheet.

**Parameters:**
- `spreadsheetId` (required): The ID of the Google Spreadsheet
- `range` (required): The range to update (e.g., 'Sheet1!A1:C10')
- `values` (required): 2D array of values to write

**Example:**
```json
{
  "name": "update_spreadsheet",
  "arguments": {
    "spreadsheetId": "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms",
    "range": "Sheet1!A1:C2",
    "values": [
      ["Name", "Age", "City"],
      ["John", "30", "Tokyo"]
    ]
  }
}
```

## Testing

```bash
go test -v
```

## Structure

- `drive.go` - Google Drive, Docs, Slides, and Sheets API operations implementation
- `main.go` - MCP server entry point with tool handlers

## License

MIT
