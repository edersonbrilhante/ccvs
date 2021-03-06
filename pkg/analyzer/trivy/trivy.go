package trivy

import (
	"context"
	"strings"
	"time"

	"net/url"

	"github.com/edersonbrilhante/vilicus/pkg/types"
	"github.com/edersonbrilhante/vilicus/pkg/util/config"

	"golang.org/x/xerrors"

	aimage "github.com/aquasecurity/fanal/artifact/image"
	"github.com/aquasecurity/fanal/image"
	"github.com/aquasecurity/trivy/pkg/cache"
	"github.com/aquasecurity/trivy/pkg/log"
	"github.com/aquasecurity/trivy/pkg/report"
	"github.com/aquasecurity/trivy/pkg/rpc/client"
	"github.com/aquasecurity/trivy/pkg/scanner"
	ttrivy "github.com/aquasecurity/trivy/pkg/types"
)

// Trivy stores pointers to config and results from scan
type Trivy struct {
	Config    *config.Trivy
	resultRaw *report.Results
	analysis  *types.Analysis
}

// RemoteURL to hold remote host
type RemoteURL string

// Analyzer runs an analysis and stores the results in Trivy.resultRaw
func (t *Trivy) Analyzer(al *types.Analysis) error {
	ctx := context.Background()
	t.analysis = al
	scanner, cleanup, err := t.dockerScanner(ctx)
	if err != nil {
		return xerrors.Errorf("error in image scan: %w", err)
	}
	defer cleanup()

	scanOptions := ttrivy.ScanOptions{
		VulnType:            []string{"os", "library"},
		ScanRemovedPackages: false,
	}

	vulns, err := scanner.ScanArtifact(ctx, scanOptions)
	if err != nil {
		return xerrors.Errorf("error in image scan: %w", err)
	}
	t.resultRaw = &vulns

	return nil
}

// Parser parses Trivy.resultRaw and store the final data into a type Analysis
func (t *Trivy) Parser() error {
	if t.resultRaw == nil {
		return xerrors.Errorf("Result is empty")
	}

	r := types.VendorResults{}
	for _, res := range *t.resultRaw {
		for _, v := range res.Vulnerabilities {

			vuln := types.Vuln{
				Fix:            v.FixedVersion,
				URL:            filterValidURLs(v.References),
				Name:           v.VulnerabilityID,
				Vendor:         "Trivy",
				Severity:       strings.Title(strings.ToLower(v.Severity)),
				PackageName:    v.PkgName,
				PackageVersion: v.InstalledVersion,
			}
			switch vuln.Severity {
			case "Unknown":
				r.UnknownVulns = append(r.UnknownVulns, vuln)
			case "Negligible":
				r.NegligibleVulns = append(r.NegligibleVulns, vuln)
			case "Low":
				r.LowVulns = append(r.LowVulns, vuln)
			case "Medium":
				r.MediumVulns = append(r.MediumVulns, vuln)
			case "High":
				r.HighVulns = append(r.HighVulns, vuln)
			case "Critical":
				r.CriticalVulns = append(r.CriticalVulns, vuln)
			}
		}

	}

	t.analysis.Results.TrivyResult = r

	return nil

}

func (t *Trivy) dockerScanner(ctx context.Context) (scanner.Scanner, func(), error) {
	err := log.InitLogger(false, true)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	customheaders := client.CustomHeaders{}
	scannerScanner := client.NewProtobufClient(client.RemoteURL(t.Config.URL))
	clientScanner := client.NewScanner(customheaders, scannerScanner)
	dockerOption, err := ttrivy.GetDockerOption(time.Duration(t.Config.Timeout) * time.Second)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	imageImage, cleanup, err := image.NewDockerImage(ctx, t.analysis.Image, dockerOption)
	if err != nil {
		return scanner.Scanner{}, nil, err
	}
	artifactCache := cache.NewRemoteCache(cache.RemoteURL(t.Config.URL), nil)
	artifact := aimage.NewArtifact(imageImage, artifactCache)
	s := scanner.NewScanner(clientScanner, artifact)
	return s, cleanup, nil
}

func filterValidURLs(urls []string) []string {
	validURLs := []string{}
	for _, ur := range urls {
		u, err := url.Parse(ur)
		if err == nil && u.Scheme != "" && u.Host != "" {
			validURLs = append(validURLs, u.String())
		}
	}
	return validURLs
}
