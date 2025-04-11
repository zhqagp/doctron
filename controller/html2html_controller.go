package controller

import (
	"context"
	"github.com/chromedp/chromedp"
	"time"

	"github.com/gorilla/schema"
	uuid "github.com/iris-contrib/go.uuid"
	"github.com/kataras/iris/v12"
	"github.com/lampnick/doctron/common"
	"github.com/lampnick/doctron/conf"
	"github.com/lampnick/doctron/converter"
	"github.com/lampnick/doctron/converter/doctron_core"
	"gopkg.in/go-playground/validator.v9"
)

func getHTMLContent(url string) (string, error) {
	// 创建一个带有取消功能的上下文
	ctx, cancel := chromedp.NewContext(context.Background())
	// 确保在函数结束时取消上下文以释放资源
	defer cancel()

	var htmlContent string
	// 定义要执行的动作序列
	actions := chromedp.Tasks{
		// 导航到指定的 URL
		chromedp.Navigate(url),
		// 提取整个 HTML 页面的内容
		chromedp.OuterHTML("html", &htmlContent),
	}

	// 执行动作序列
	err := chromedp.Run(ctx, actions)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}

func Html2HtmlHandler(ctx iris.Context) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(conf.LoadedConfig.Doctron.ConvertTimeout)*time.Second)
	defer cancel()
	traceId, _ := uuid.NewV4()
	outputDTO := common.NewDefaultOutputDTO(nil)
	doctronConfig, err := initHtml2HtmlConfig(ctx)
	if err != nil {
		outputDTO.Code = common.InvalidParams
		outputDTO.Message = err.Error()
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}
	doctronConfig.TraceId = traceId
	doctronConfig.Ctx = ctxTimeout
	doctronConfig.IrisCtx = ctx
	doctronConfig.DoctronType = doctron_core.DoctronHtml2Pdf

	doctronOutputDTO, err := getHTMLContent(doctronConfig.Url)
	if err != nil {
		outputDTO.Code = common.ConvertPdfFailed
		outputDTO.Message = "worker run process failed." + err.Error()
		_, _ = common.NewJsonOutput(ctx, outputDTO)
		return
	}

	outputDTO.Data = doctronOutputDTO
	_, _ = common.NewJsonOutput(ctx, outputDTO)
	return
}

func initHtml2HtmlConfig(ctx iris.Context) (converter.DoctronConfig, error) {
	result := converter.DoctronConfig{}
	requestDTO := newDefaultHtml2HtmlRequestDTO()
	// decode query params
	var decoder = schema.NewDecoder()
	err := decoder.Decode(requestDTO, ctx.Request().URL.Query())
	if err != nil {
		log(ctx, "[%s][%s]", ctx.Request().RequestURI, err)
	}
	// validate
	validate := validator.New()
	err = validate.Struct(requestDTO)
	if err != nil {
		return result, err
	}
	result.Url = requestDTO.Url
	result.UploadKey = requestDTO.UploadKey
	result.Params = convertToHtmlParams(requestDTO)
	return result, nil
}

func convertToHtmlParams(requestDTO *Html2HtmlRequestDTO) doctron_core.PDFParams {
	params := doctron_core.NewDefaultPDFParams()
	params.Landscape = requestDTO.Landscape
	params.DisplayHeaderFooter = requestDTO.DisplayHeaderFooter
	params.PrintBackground = requestDTO.PrintBackground
	params.Scale = requestDTO.Scale
	params.PaperWidth = requestDTO.PaperWidth
	params.PaperHeight = requestDTO.PaperHeight
	params.MarginTop = requestDTO.MarginTop
	params.MarginBottom = requestDTO.MarginBottom
	params.MarginLeft = requestDTO.MarginLeft
	params.MarginRight = requestDTO.MarginRight
	params.PageRanges = requestDTO.PageRanges
	params.IgnoreInvalidPageRanges = requestDTO.IgnoreInvalidPageRanges
	params.HeaderTemplate = requestDTO.HeaderTemplate
	params.FooterTemplate = requestDTO.FooterTemplate
	params.PreferCSSPageSize = requestDTO.PreferCSSPageSize
	params.WaitingTime = requestDTO.WaitingTime
	return params
}

func newDefaultHtml2HtmlRequestDTO() *Html2HtmlRequestDTO {
	return &Html2HtmlRequestDTO{
		Landscape:               doctron_core.DefaultLandscape,
		DisplayHeaderFooter:     doctron_core.DefaultDisplayHeaderFooter,
		PrintBackground:         doctron_core.DefaultPrintBackground,
		Scale:                   doctron_core.DefaultScale,
		PaperWidth:              doctron_core.DefaultPaperWidth,
		PaperHeight:             doctron_core.DefaultPaperHeight,
		MarginTop:               doctron_core.DefaultMarginTop,
		MarginBottom:            doctron_core.DefaultMarginBottom,
		MarginLeft:              doctron_core.DefaultMarginLeft,
		MarginRight:             doctron_core.DefaultMarginRight,
		PageRanges:              doctron_core.DefaultPageRanges,
		IgnoreInvalidPageRanges: doctron_core.DefaultIgnoreInvalidPageRanges,
		HeaderTemplate:          doctron_core.DefaultHeaderTemplate,
		FooterTemplate:          doctron_core.DefaultFooterTemplate,
		PreferCSSPageSize:       doctron_core.DefaultPreferCSSPageSize,
		WaitingTime:             doctron_core.DefaultWaitingTime,
	}
}
