package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func createSearchFilesHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		query, err := request.RequireString("query")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'query' is required"), nil
		}

		maxResults := mcp.ParseInt(request, "maxResults", 10)

		// Execute Google Drive search
		files, err := driveService.SearchFiles(ctx, query, maxResults)
		if err != nil {
			return mcp.NewToolResultError("Failed to search files: " + err.Error()), nil
		}

		// Convert result to JSON
		result := map[string]any{
			"files": files,
			"count": len(files),
		}

		resultData, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
		}

		return mcp.NewToolResultText(string(resultData)), nil
	}
}

func createListFilesHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		folderID := mcp.ParseString(request, "folderId", "")
		maxResults := mcp.ParseInt(request, "maxResults", 10)

		// Execute Google Drive list
		files, err := driveService.ListFiles(ctx, folderID, maxResults)
		if err != nil {
			return mcp.NewToolResultError("Failed to list files: " + err.Error()), nil
		}

		// Convert result to JSON
		result := map[string]any{
			"files": files,
			"count": len(files),
		}

		resultData, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
		}

		return mcp.NewToolResultText(string(resultData)), nil
	}
}

func createGetDocumentHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		documentID, err := request.RequireString("documentId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'documentId' is required"), nil
		}

		// Get document content
		content, err := driveService.GetDocumentContent(ctx, documentID)
		if err != nil {
			return mcp.NewToolResultError("Failed to get document content: " + err.Error()), nil
		}

		return mcp.NewToolResultText(content), nil
	}
}

func createUpdateDocumentHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		documentID, err := request.RequireString("documentId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'documentId' is required"), nil
		}

		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'content' is required"), nil
		}

		// Update document content
		err = driveService.UpdateDocumentContent(ctx, documentID, content)
		if err != nil {
			return mcp.NewToolResultError("Failed to update document: " + err.Error()), nil
		}

		return mcp.NewToolResultText("Document updated successfully"), nil
	}
}

func createGetPresentationHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		presentationID, err := request.RequireString("presentationId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'presentationId' is required"), nil
		}

		// Get presentation content
		content, err := driveService.GetPresentationContent(ctx, presentationID)
		if err != nil {
			return mcp.NewToolResultError("Failed to get presentation content: " + err.Error()), nil
		}

		return mcp.NewToolResultText(content), nil
	}
}

func createUpdatePresentationHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		presentationID, err := request.RequireString("presentationId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'presentationId' is required"), nil
		}

		slideIndex := mcp.ParseInt(request, "slideIndex", 0)

		title, err := request.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'title' is required"), nil
		}

		content, err := request.RequireString("content")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'content' is required"), nil
		}

		// Update presentation slide
		err = driveService.UpdatePresentationSlide(ctx, presentationID, slideIndex, title, content)
		if err != nil {
			return mcp.NewToolResultError("Failed to update presentation: " + err.Error()), nil
		}

		return mcp.NewToolResultText("Presentation slide updated successfully"), nil
	}
}

func createGetSpreadsheetHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		spreadsheetID, err := request.RequireString("spreadsheetId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'spreadsheetId' is required"), nil
		}

		rangeName, err := request.RequireString("range")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'range' is required"), nil
		}

		// Get spreadsheet values
		values, err := driveService.GetSpreadsheetValues(ctx, spreadsheetID, rangeName)
		if err != nil {
			return mcp.NewToolResultError("Failed to get spreadsheet values: " + err.Error()), nil
		}

		// Convert result to JSON
		result := map[string]any{
			"values": values,
			"range":  rangeName,
		}

		resultData, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultError("Failed to serialize result: " + err.Error()), nil
		}

		return mcp.NewToolResultText(string(resultData)), nil
	}
}

func createUpdateSpreadsheetHandler(driveService *DriveService) func(context.Context, mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Get parameters
		spreadsheetID, err := request.RequireString("spreadsheetId")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'spreadsheetId' is required"), nil
		}

		rangeName, err := request.RequireString("range")
		if err != nil {
			return mcp.NewToolResultError("Parameter 'range' is required"), nil
		}

		valuesParam := request.Params["values"]
		if valuesParam == nil {
			return mcp.NewToolResultError("Parameter 'values' is required"), nil
		}

		// Convert values to [][]interface{}
		var values [][]interface{}
		if valuesSlice, ok := valuesParam.([]interface{}); ok {
			for _, row := range valuesSlice {
				if rowSlice, ok := row.([]interface{}); ok {
					values = append(values, rowSlice)
				} else {
					return mcp.NewToolResultError("Invalid values format: each row must be an array"), nil
				}
			}
		} else {
			return mcp.NewToolResultError("Invalid values format: values must be a 2D array"), nil
		}

		// Update spreadsheet values
		err = driveService.UpdateSpreadsheetValues(ctx, spreadsheetID, rangeName, values)
		if err != nil {
			return mcp.NewToolResultError("Failed to update spreadsheet: " + err.Error()), nil
		}

		return mcp.NewToolResultText("Spreadsheet updated successfully"), nil
	}
}

func main() {
	// Initialize Drive service once
	ctx := context.Background()
	driveService, err := NewDriveService(ctx)
	if err != nil {
		log.Fatal("Failed to initialize Drive service:", err)
	}

	s := server.NewMCPServer("Google Drive MCP", "1.0.0", server.WithToolCapabilities(true))

	// Define file search tool
	searchFilesTool := mcp.NewTool(
		"search_files",
		mcp.WithDescription("Search files in Google Drive"),
		mcp.WithString("query", mcp.Description("File name or keyword to search"), mcp.Required()),
		mcp.WithNumber("maxResults", mcp.Description("Maximum number of files to retrieve (default: 10)"), mcp.DefaultNumber(10)),
	)

	// Define list files tool
	listFilesTool := mcp.NewTool(
		"list_files",
		mcp.WithDescription("List files in a Google Drive folder"),
		mcp.WithString("folderId", mcp.Description("The ID of the folder to list files from. If empty, lists files in My Drive root")),
		mcp.WithNumber("maxResults", mcp.Description("Maximum number of files to retrieve (default: 10)"), mcp.DefaultNumber(10)),
	)

	// Define get document tool
	getDocumentTool := mcp.NewTool(
		"get_document",
		mcp.WithDescription("Get the content of a Google Document"),
		mcp.WithString("documentId", mcp.Description("The ID of the Google Document"), mcp.Required()),
	)

	// Define update document tool
	updateDocumentTool := mcp.NewTool(
		"update_document",
		mcp.WithDescription("Update the content of a Google Document"),
		mcp.WithString("documentId", mcp.Description("The ID of the Google Document"), mcp.Required()),
		mcp.WithString("content", mcp.Description("The new content for the document"), mcp.Required()),
	)

	// Define get presentation tool
	getPresentationTool := mcp.NewTool(
		"get_presentation",
		mcp.WithDescription("Get the content of a Google Slides presentation"),
		mcp.WithString("presentationId", mcp.Description("The ID of the Google Slides presentation"), mcp.Required()),
	)

	// Define update presentation tool
	updatePresentationTool := mcp.NewTool(
		"update_presentation",
		mcp.WithDescription("Update a specific slide in a Google Slides presentation"),
		mcp.WithString("presentationId", mcp.Description("The ID of the Google Slides presentation"), mcp.Required()),
		mcp.WithNumber("slideIndex", mcp.Description("The index of the slide to update (0-based, default: 0)"), mcp.DefaultNumber(0)),
		mcp.WithString("title", mcp.Description("The title for the slide"), mcp.Required()),
		mcp.WithString("content", mcp.Description("The content for the slide"), mcp.Required()),
	)

	// Define get spreadsheet tool
	getSpreadsheetTool := mcp.NewTool(
		"get_spreadsheet",
		mcp.WithDescription("Get values from a Google Spreadsheet"),
		mcp.WithString("spreadsheetId", mcp.Description("The ID of the Google Spreadsheet"), mcp.Required()),
		mcp.WithString("range", mcp.Description("The range to retrieve (e.g., 'Sheet1!A1:C10')"), mcp.Required()),
	)

	// Define update spreadsheet tool
	updateSpreadsheetTool := mcp.NewTool(
		"update_spreadsheet",
		mcp.WithDescription("Update values in a Google Spreadsheet"),
		mcp.WithString("spreadsheetId", mcp.Description("The ID of the Google Spreadsheet"), mcp.Required()),
		mcp.WithString("range", mcp.Description("The range to update (e.g., 'Sheet1!A1:C10')"), mcp.Required()),
		mcp.WithAny("values", mcp.Description("2D array of values to write"), mcp.Required()),
	)

	// Register tool handlers
	s.AddTool(searchFilesTool, createSearchFilesHandler(driveService))
	s.AddTool(listFilesTool, createListFilesHandler(driveService))
	s.AddTool(getDocumentTool, createGetDocumentHandler(driveService))
	s.AddTool(updateDocumentTool, createUpdateDocumentHandler(driveService))
	s.AddTool(getPresentationTool, createGetPresentationHandler(driveService))
	s.AddTool(updatePresentationTool, createUpdatePresentationHandler(driveService))
	s.AddTool(getSpreadsheetTool, createGetSpreadsheetHandler(driveService))
	s.AddTool(updateSpreadsheetTool, createUpdateSpreadsheetHandler(driveService))

	// Start server
	if err := server.ServeStdio(s); err != nil {
		log.Fatal("Failed to start MCP server:", err)
	}
}
