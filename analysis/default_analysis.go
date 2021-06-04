package analysis

const DefaultAnalysisName = "defaultAnalyser"

type DefaultAnalysis struct {
}

func NewDefaultAnalysis() *DefaultAnalysis {
	return &DefaultAnalysis{}
}

func (o DefaultAnalysis) AnalysisResult(result, transCode string) string {
	return result
}
