package googleapi

type SearchRequest struct {
	Query    string
	Corpora  string
	DriveID  string
	PageSize int
	Fields   string
}

type FileMetaRequest struct {
	ID     string
	Fields string
}

type ExportRequest struct {
	ID   string
	MIME string
}

type DocTabsRequest struct {
	ID string
}
