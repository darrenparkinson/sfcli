package salesforce

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// BulkType represents the type of bulk operation required for bulk operations
// for example ingest or query
type BulkType int

const (
	// BulkTypeIngest represents the ingest bulk type
	BulkTypeIngest BulkType = iota + 1
	// BulkTypeQuery represents the query bulk type
	BulkTypeQuery
)

func (b BulkType) String() string {
	return [...]string{"ingest", "query"}[b-1]
}

// EnumIndex provides the index for a given bulk type
func (b BulkType) EnumIndex() int {
	return int(b)
}

// JobInfo represents a specific Job
type JobInfo struct {
	ID                     string  `json:"id"`
	Operation              string  `json:"operation"`
	Object                 string  `json:"object"`
	CreatedByID            string  `json:"createdById"`
	CreatedDate            string  `json:"createdDate"`
	SystemModstamp         string  `json:"systemModstamp"`
	State                  string  `json:"state"`
	ConcurrencyMode        string  `json:"concurrencyMode"`
	ContentType            string  `json:"contentType"`
	APIVersion             float64 `json:"apiVersion"`
	JobType                string  `json:"jobType"`
	LineEnding             string  `json:"lineEnding"`
	ColumnDelimiter        string  `json:"columnDelimiter"`
	NumberRecordsProcessed int     `json:"numberRecordsProcessed"`
	NumberRecordsFailed    int     `json:"numberRecordsFailed"`
	Retries                int     `json:"retries"`
}

// BulkRequest represents the object required to send when creating a Job Request
type BulkRequest struct {
	Object              string `json:"object,omitempty"`              // e.g. Account
	ContentType         string `json:"contentType,omitempty"`         // e.g. CSV
	Operation           string `json:"operation,omitempty"`           // e.g. insert,upsert,query
	LineEnding          string `json:"lineEnding,omitempty"`          // e.g. CRLF (windows). Default is LF
	ColumnDelimiter     string `json:"columnDelimiter,omitempty"`     // e.g. SEMICOLON. Default is COMMA
	ExternalIDFieldName string `json:"externalIdFieldName,omitempty"` // required only for Upserts
	Query               string `json:"query,omitempty"`               // required only for query operations
}

// BulkListResponse is the response from Salesforce listing the jobs
type BulkListResponse struct {
	Done           bool      `json:"done"`
	NextRecordsURL string    `json:"nextRecordsUrl"`
	Records        []JobInfo `json:"records"`
	// TotalSize      int       `json:"totalSize"`
}

// ListJobs lists the first 1000 jobs of type BulkType
// TODO: Add option to list all?
// https://developer.salesforce.com/docs/atlas.en-us.api_asynch.meta/api_asynch/get_all_jobs.htm
func (s *BulkService) ListJobs(ctx context.Context, jobType BulkType) (*BulkListResponse, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/%s", s.client.BaseURL, s.client.Version, jobType)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var blr BulkListResponse
	if err := s.client.makeRequest(ctx, req, &blr); err != nil {
		return nil, err
	}
	return &blr, nil
}

// GetJob allows you to get details for a specific job id
func (s *BulkService) GetJob(ctx context.Context, jobType BulkType, id string) (*JobInfo, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/%s/%s", s.client.BaseURL, s.client.Version, jobType, id)
	req, err := http.NewRequest("GET", sfurl, nil)
	if err != nil {
		return nil, err
	}
	var job JobInfo
	if err := s.client.makeRequest(ctx, req, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// CreateJob allows you to create a new bulk job request
func (s *BulkService) CreateJob(ctx context.Context, br BulkRequest) (*JobInfo, error) {
	// calculate job type, e.g. ingest or query based on the operation type
	jobType := BulkTypeIngest
	if strings.ToLower(br.Operation) == "query" {
		jobType = BulkTypeQuery
	}
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/%s", s.client.BaseURL, s.client.Version, jobType)
	payload, err := json.Marshal(br)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", sfurl, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	var job JobInfo
	if err := s.client.makeRequest(ctx, req, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// CancelJob allows you to cancel an existing job request
func (s *BulkService) CancelJob(ctx context.Context, jobType BulkType, id string) (*JobInfo, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/%s/%s", s.client.BaseURL, s.client.Version, jobType, id)
	payload := strings.NewReader(`{ "state" : "Aborted" }`)
	req, err := http.NewRequest("PATCH", sfurl, payload)
	if err != nil {
		return nil, err
	}
	var job JobInfo
	if err := s.client.makeRequest(ctx, req, &job); err != nil {
		return nil, err
	}
	return &job, nil
}

// UploadCSV will upload CSV data from the provided io.Reader to the provided job id
// You must remember to begin processing the job and then check it for success/errors.
// Note that only a single upload/batch is currently supported.
func (s *BulkService) UploadCSV(ctx context.Context, id string, payload io.Reader) error {
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/ingest/%s/batches", s.client.BaseURL, s.client.Version, id)
	req, err := http.NewRequest("PUT", sfurl, payload)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/csv")
	if err := s.client.makeRequest(ctx, req, nil); err != nil {
		return err
	}
	return nil
}

// ProcessJob marks a job as UploadComplete and begins processing
func (s *BulkService) ProcessJob(ctx context.Context, jobType BulkType, id string) (*JobInfo, error) {
	sfurl := fmt.Sprintf("%s/services/data/%s/jobs/%s/%s", s.client.BaseURL, s.client.Version, jobType, id)
	payload := strings.NewReader(`{ "state" : "UploadComplete" }`)
	req, err := http.NewRequest("PATCH", sfurl, payload)
	if err != nil {
		return nil, err
	}
	var job JobInfo
	if err := s.client.makeRequest(ctx, req, &job); err != nil {
		return nil, err
	}
	return &job, nil
}
