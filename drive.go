package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"google.golang.org/api/docs/v1"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"google.golang.org/api/slides/v1"
)

// DriveFile represents information about a Google Drive file
type DriveFile struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"mimeType"`
}

// DriveService manages Google Drive, Docs, Slides, and Sheets API services
type DriveService struct {
	driveService  *drive.Service
	docsService   *docs.Service
	slidesService *slides.Service
	sheetsService *sheets.Service
}

// NewDriveService creates a new DriveService
func NewDriveService(ctx context.Context) (*DriveService, error) {
	// Use gcloud application-default credentials
	options := []option.ClientOption{
		option.WithScopes(drive.DriveScope, docs.DocumentsScope, slides.PresentationsScope, sheets.SpreadsheetsScope),
	}

	// Use quota project if set in environment variable
	if quotaProject := os.Getenv("GOOGLE_CLOUD_QUOTA_PROJECT_ID"); quotaProject != "" {
		options = append(options, option.WithQuotaProject(quotaProject))
	}

	driveService, err := drive.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create drive service: %w", err)
	}

	docsService, err := docs.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create docs service: %w", err)
	}

	slidesService, err := slides.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create slides service: %w", err)
	}

	sheetsService, err := sheets.NewService(ctx, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create sheets service: %w", err)
	}

	return &DriveService{
		driveService:  driveService,
		docsService:   docsService,
		slidesService: slidesService,
		sheetsService: sheetsService,
	}, nil
}

// SearchFiles searches for files in Google Drive (DriveService method)
func (ds *DriveService) SearchFiles(ctx context.Context, query string, maxResults int) ([]DriveFile, error) {
	if query == "" {
		return nil, errors.New("search query is empty")
	}

	// Execute search with Google Drive API
	searchQuery := fmt.Sprintf("name contains '%s'", query)
	r, err := ds.driveService.Files.List().
		Q(searchQuery).
		PageSize(int64(maxResults)).
		Fields("nextPageToken, files(id, name, mimeType)").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to search files: %w", err)
	}

	var files []DriveFile
	for _, file := range r.Files {
		files = append(files, DriveFile{
			ID:   file.Id,
			Name: file.Name,
			Type: file.MimeType,
		})
	}

	return files, nil
}

// ListFiles lists files in a Google Drive folder
func (ds *DriveService) ListFiles(ctx context.Context, folderID string, maxResults int) ([]DriveFile, error) {
	// Build query for listing files in folder
	var query string
	if folderID == "" {
		// List files in root folder (My Drive)
		query = "'root' in parents and trashed = false"
	} else {
		// List files in specific folder
		query = fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	}

	// Execute list with Google Drive API
	r, err := ds.driveService.Files.List().
		Q(query).
		PageSize(int64(maxResults)).
		Fields("nextPageToken, files(id, name, mimeType)").
		Context(ctx).
		Do()
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var files []DriveFile
	for _, file := range r.Files {
		files = append(files, DriveFile{
			ID:   file.Id,
			Name: file.Name,
			Type: file.MimeType,
		})
	}

	return files, nil
}

// GetDocumentContent retrieves the content of a Google Document
func (ds *DriveService) GetDocumentContent(ctx context.Context, documentID string) (string, error) {
	if documentID == "" {
		return "", errors.New("document ID is empty")
	}

	doc, err := ds.docsService.Documents.Get(documentID).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get document: %w", err)
	}

	var content string
	for _, element := range doc.Body.Content {
		if element.Paragraph != nil {
			for _, elem := range element.Paragraph.Elements {
				if elem.TextRun != nil {
					content += elem.TextRun.Content
				}
			}
		}
	}

	return content, nil
}

// UpdateDocumentContent updates the content of a Google Document
func (ds *DriveService) UpdateDocumentContent(ctx context.Context, documentID, content string) error {
	if documentID == "" {
		return errors.New("document ID is empty")
	}

	// First, get the current document to determine the end index
	doc, err := ds.docsService.Documents.Get(documentID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Calculate the end index of the document content
	endIndex := int64(1)
	for _, element := range doc.Body.Content {
		if element.EndIndex > endIndex {
			endIndex = element.EndIndex
		}
	}

	// Create batch update requests
	requests := []*docs.Request{
		// Delete all existing content (except the last character which is always a newline)
		{
			DeleteContentRange: &docs.DeleteContentRangeRequest{
				Range: &docs.Range{
					StartIndex: 1,
					EndIndex:   endIndex - 1,
				},
			},
		},
		// Insert new content
		{
			InsertText: &docs.InsertTextRequest{
				Location: &docs.Location{
					Index: 1,
				},
				Text: content,
			},
		},
	}

	// Execute the batch update
	batchUpdateRequest := &docs.BatchUpdateDocumentRequest{
		Requests: requests,
	}

	_, err = ds.docsService.Documents.BatchUpdate(documentID, batchUpdateRequest).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// GetPresentationContent retrieves the content of a Google Slides presentation
func (ds *DriveService) GetPresentationContent(ctx context.Context, presentationID string) (string, error) {
	if presentationID == "" {
		return "", errors.New("presentation ID is empty")
	}

	presentation, err := ds.slidesService.Presentations.Get(presentationID).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("failed to get presentation: %w", err)
	}

	var content string
	content += fmt.Sprintf("Title: %s\n\n", presentation.Title)

	for i, slide := range presentation.Slides {
		content += fmt.Sprintf("--- Slide %d ---\n", i+1)
		
		for _, element := range slide.PageElements {
			if element.Shape != nil && element.Shape.Text != nil {
				for _, textElement := range element.Shape.Text.TextElements {
					if textElement.TextRun != nil {
						content += textElement.TextRun.Content
					}
				}
				content += "\n"
			}
		}
		content += "\n"
	}

	return content, nil
}

// UpdatePresentationSlide updates a specific slide in a Google Slides presentation
func (ds *DriveService) UpdatePresentationSlide(ctx context.Context, presentationID string, slideIndex int, title, content string) error {
	if presentationID == "" {
		return errors.New("presentation ID is empty")
	}

	presentation, err := ds.slidesService.Presentations.Get(presentationID).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("failed to get presentation: %w", err)
	}

	if slideIndex < 0 || slideIndex >= len(presentation.Slides) {
		return fmt.Errorf("slide index %d is out of range (0-%d)", slideIndex, len(presentation.Slides)-1)
	}

	slide := presentation.Slides[slideIndex]
	var requests []*slides.Request

	// Clear existing text elements
	for _, element := range slide.PageElements {
		if element.Shape != nil && element.Shape.Text != nil {
			requests = append(requests, &slides.Request{
				DeleteText: &slides.DeleteTextRequest{
					ObjectId: element.ObjectId,
					TextRange: &slides.Range{
						Type: "ALL",
					},
				},
			})
		}
	}

	// Find title and content text boxes, or create new ones if needed
	var titleObjectId, contentObjectId string
	
	for _, element := range slide.PageElements {
		if element.Shape != nil {
			// Assume first text box is title, second is content
			if titleObjectId == "" {
				titleObjectId = element.ObjectId
			} else if contentObjectId == "" {
				contentObjectId = element.ObjectId
				break
			}
		}
	}

	// Insert new content
	if titleObjectId != "" && title != "" {
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: titleObjectId,
				Text:     title,
				InsertionIndex: 0,
			},
		})
	}

	if contentObjectId != "" && content != "" {
		requests = append(requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: contentObjectId,
				Text:     content,
				InsertionIndex: 0,
			},
		})
	}

	if len(requests) > 0 {
		batchUpdateRequest := &slides.BatchUpdatePresentationRequest{
			Requests: requests,
		}

		_, err = ds.slidesService.Presentations.BatchUpdate(presentationID, batchUpdateRequest).Context(ctx).Do()
		if err != nil {
			return fmt.Errorf("failed to update presentation: %w", err)
		}
	}

	return nil
}

// GetSpreadsheetValues retrieves values from a Google Spreadsheet
func (ds *DriveService) GetSpreadsheetValues(ctx context.Context, spreadsheetID, rangeName string) ([][]interface{}, error) {
	if spreadsheetID == "" {
		return nil, errors.New("spreadsheet ID is empty")
	}
	if rangeName == "" {
		return nil, errors.New("range name is empty")
	}

	resp, err := ds.sheetsService.Spreadsheets.Values.Get(spreadsheetID, rangeName).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to get spreadsheet values: %w", err)
	}

	return resp.Values, nil
}

// UpdateSpreadsheetValues updates values in a Google Spreadsheet
func (ds *DriveService) UpdateSpreadsheetValues(ctx context.Context, spreadsheetID, rangeName string, values [][]interface{}) error {
	if spreadsheetID == "" {
		return errors.New("spreadsheet ID is empty")
	}
	if rangeName == "" {
		return errors.New("range name is empty")
	}

	valueRange := &sheets.ValueRange{
		Values: values,
	}

	_, err := ds.sheetsService.Spreadsheets.Values.Update(spreadsheetID, rangeName, valueRange).
		ValueInputOption("USER_ENTERED").
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("failed to update spreadsheet values: %w", err)
	}

	return nil
}
