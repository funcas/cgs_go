package analysis

import (
	"fmt"
	"testing"

	"github.com/funcas/cgs/manager"
)

func TestAnalysisResult(t *testing.T) {
	input := `{
    "order_id": "1234567",
    "tracking_number": "1z9999999999999999",
    "items": [
        {
            "item_sku": "ab123",
            "item_price": 12.34,
            "number_purchased": 5
        },
        {
            "item_sku": "ck763-23",
            "item_price": 3.12,
            "number_purchased": 2
        }
    ]
}
`
	manager.LoadTransformer()
	ana := NewOmniAnalysis()
	ret := ana.AnalysisResult(input, "sample")
	fmt.Println(ret)
}
