package app

import (
	"net/http"

	"github.com/Alechan/gdrivectl/internal/auth"
	"github.com/Alechan/gdrivectl/internal/googleapi"
	"github.com/Alechan/gdrivectl/internal/output"
	"github.com/Alechan/gdrivectl/internal/service"
)

type Services struct {
	Doctor    *service.DoctorService
	Search    *service.SearchService
	FileMeta  *service.FileMetaService
	DocTabs   *service.DocTabsService
	DocExport *service.DocExportService
	Output    *output.Writer
}

func NewServices(cfg Config) Services {
	httpClient := &http.Client{Timeout: cfg.Timeout}
	tokenProvider := auth.NewGcloudTokenProvider(cfg.GcloudBin)
	driveClient := googleapi.NewDriveClient(httpClient)
	docsClient := googleapi.NewDocsClient(httpClient)

	return Services{
		Doctor:    service.NewDoctorService(cfg.GcloudBin, tokenProvider, driveClient, docsClient),
		Search:    service.NewSearchService(tokenProvider, driveClient),
		FileMeta:  service.NewFileMetaService(tokenProvider, driveClient),
		DocTabs:   service.NewDocTabsService(tokenProvider, docsClient),
		DocExport: service.NewDocExportService(tokenProvider, driveClient),
		Output:    output.NewWriter(),
	}
}
