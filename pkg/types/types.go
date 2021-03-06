package types

import (
	"fmt"
	"strings"
	"time"
)

// Analyses is a list of Analysis
type Analyses []Analysis

// Analysis is the struct that stores all data from analysis performed.
type Analysis struct {
	ID        string    `json:"id,omitempty" pg:"id,type:uuid,pk,default:uuid_generate_v4()"`
	Image     string    `json:"image,omitempty" pg:"image" validate:"required"`
	Status    string    `json:"status,omitempty" pg:"status"`
	CreatedAt time.Time `json:"created_at,omitempty" pg:"created_at,notnull,default:now()"`
	UpdatedAt time.Time `json:"updated_at,omitempty" pg:"updated_at,notnull,default:now()"`
	Result    string    `json:"result,omitempty" pg:"result"`
	Errors    []string  `json:"errors,omitempty" pg:"errors"`
	Results   Results   `json:"vilicus_results,omitempty" pg:"vilicus_results"`
}

// Results is a struct that represents vilicus scan results.
type Results struct {
	ClairResult         VendorResults `json:"clair,omitempty"`
	AnchoreEngineResult VendorResults `json:"anchore_engine,omitempty"`
	TrivyResult         VendorResults `json:"trivy,omitempty"`
}

// VulnList is a function that return a merged list with all vulns
func (r Results) VulnList() []Vuln {
	vendorVulns := [][]Vuln{
		r.AnchoreEngineResult.CriticalVulns,
		r.ClairResult.CriticalVulns,
		r.TrivyResult.CriticalVulns,
		r.AnchoreEngineResult.HighVulns,
		r.ClairResult.HighVulns,
		r.TrivyResult.HighVulns,
		r.AnchoreEngineResult.MediumVulns,
		r.ClairResult.MediumVulns,
		r.TrivyResult.MediumVulns,
		r.AnchoreEngineResult.LowVulns,
		r.ClairResult.LowVulns,
		r.TrivyResult.LowVulns,
		r.AnchoreEngineResult.NegligibleVulns,
		r.ClairResult.NegligibleVulns,
		r.TrivyResult.NegligibleVulns,
		r.AnchoreEngineResult.UnknownVulns,
		r.ClairResult.UnknownVulns,
		r.TrivyResult.UnknownVulns,
	}

	var vulns []Vuln
	for _, v := range vendorVulns {
		vulns = append(vulns, v...)
	}
	return vulns
}

func (r Results) String() string {
	clairResult := strings.ReplaceAll(r.ClairResult.String(), "VendorResults", "Clair")
	anchoreEngineResult := strings.ReplaceAll(r.AnchoreEngineResult.String(), "VendorResults", "AnchoreEngine")
	trivyResult := strings.ReplaceAll(r.TrivyResult.String(), "VendorResults", "Trivy")

	return fmt.Sprintf("%s %s %s", clairResult, anchoreEngineResult, trivyResult)
}

// VendorResults stores all Unknown, Negligible Low, Medium, High and Critical vulnerabilities for a vendor
type VendorResults struct {
	UnknownVulns    []Vuln `json:"unknown_vulns,omitempty"`
	NegligibleVulns []Vuln `json:"negligible_vulns,omitempty"`
	LowVulns        []Vuln `json:"low_vulns,omitempty"`
	MediumVulns     []Vuln `json:"medium_vulns,omitempty"`
	HighVulns       []Vuln `json:"high_vulns,omitempty"`
	CriticalVulns   []Vuln `json:"critical_vulns,omitempty"`
}

func (v VendorResults) String() string {
	return fmt.Sprintf(
		"VendorResults<Unknown:%d | Negligible:%d | Low:%d | Medium:%d | High:%d | Critical%d>",
		len(v.UnknownVulns), len(v.NegligibleVulns), len(v.LowVulns),
		len(v.MediumVulns), len(v.HighVulns), len(v.CriticalVulns),
	)
}

// Vuln is the struct that stores vulnerability information.
type Vuln struct {
	Fix            string   `json:"fix"`
	URL            []string `json:"urls"`
	Name           string   `json:"name"`
	Severity       string   `json:"severity"`
	Vendor         string   `json:"vendor"`
	PackageName    string   `json:"package_name"`
	PackageVersion string   `json:"package_version"`
}

func (v Vuln) String() string {
	return fmt.Sprintf(
		"Vuln<[%s][%s][%s] %s(%s) - Fix:%s>",
		v.Vendor, v.Severity, v.Name, v.PackageName, v.PackageVersion, v.Fix,
	)
}
