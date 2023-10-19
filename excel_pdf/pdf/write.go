package pdf

import (
    "fmt"
    "github.com/signintech/gopdf"
    log "github.com/sirupsen/logrus"
    "personal_work/utils"
)

func GeneratePdf() {
    pdf := gopdf.GoPdf{}
    pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
    pdf.AddPage()
    pathDir := utils.GetFilePath()
    fmt.Println("pathDir :", pathDir)
    err := pdf.AddTTFFont("wts11", pathDir+"/excel_pdf/pdf/NotoSansSC-Regular.ttf")
    if err != nil {
        log.Errorf("AddTTFFont Error:%v", err.Error())
        return
    }

    err = pdf.SetFont("wts11", "", 14)
    if err != nil {
        log.Errorf("SetFont Error:%v", err.Error())
        return
    }
    err = pdf.Cell(nil, "您好")
    if err != nil {
        log.Errorf("Cell Error:%v", err.Error())
        return
    }
    err = pdf.WritePdf(pathDir + "/excel_pdf/pdf/text.pdf")
    if err != nil {
        log.Errorf("WritePdf Error:%v", err.Error())
        return
    }
}
