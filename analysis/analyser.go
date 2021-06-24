package analysis

type Analyser interface {
	AnalysisResult(result, transCode string) string
}
