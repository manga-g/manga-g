package app

import (
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
)

//	func CreatePdf() {
//	   m := pdf.NewMaroto(consts.Portrait, consts.A4)
//	   m.SetPageMargins(20, 10, 20)
//	   err := m.OutputFileAndClose("pdf/test.pdf")
//	   if err != nil {
//	       fmt.Println("failed to create pdf", err)
//	   }
//
// }
func GetImageList() []string {
	var imageList []string
	//	var n :=0

	return imageList
}

func ImagesToPdf() {
	m := pdf.NewMaroto(consts.Portrait, consts.A4)
	m.SetPageMargins(20, 10, 20)

}
