package analysis

import (
	"fmt"
	"testing"

	"github.com/funcas/cgs/manager"
)

func TestAnalysisResult(t *testing.T) {
	input := `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema">
   <soap:Body>
      <getEnCnTwoWayTranslatorResponse xmlns="http://WebXml.com.cn/">
         <getEnCnTwoWayTranslatorResult>
            <string>misty: [ 'misti ]</string>
            <string>a. 有雾的,模糊的,含糊的 | 
词形变化:副词:mistily 形容词比较级:mistier 最高级:mistiest 名词:mistiness  |</string>
         </getEnCnTwoWayTranslatorResult>
      </getEnCnTwoWayTranslatorResponse>
   </soap:Body>
</soap:Envelope>
`
	manager.LoadTransformer()
	ana := NewOmniAnalysis()
	ret := ana.AnalysisResult(input, "test4")
	fmt.Println(ret)
}
