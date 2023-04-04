package mainWindow

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

const windowWidth = 10000
const windowHeight = 10000

func Run(host string, port int, uiPrefix string, log *logrus.Logger) {
	log.Infof("Starting main window\n")

	//view := webview.New(true)
	//defer view.Destroy()
	//view.SetTitle("Go Desktop Application")
	//view.SetSize(windowWidth, windowHeight, webview.HintNone)
	////view.Navigate(fmt.Sprintf("http://%s:%d%s", host, port, uiPrefix))
	//view.Navigate(fmt.Sprintf("http://ibm.com"))
	//view.Run()

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		//chromedp.Flag("app", true),
		//chromedp.Flag("app", "data:text/html,<html><body><script>window.moveTo(580,240);window.resizeTo(800,600);window.location='http://ibm.com';</script></body></html>"),
		//chromedp.Flag("app", "data:text/html,<title>TEST</title>"),
		chromedp.Flag("app", "data:text/html,<title>App</title><style>html{background: #000000;}</style>"),
		//chromedp.Flag("app", "data:text/html,<title>&lrm;</title>"),
		chromedp.Flag("window-size", "7800,7600"),
		chromedp.Flag("window-position", "100,200"),
	}

	allocatorContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(
		allocatorContext,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("http://%s:%d%s", host, port, uiPrefix)),
	}); err != nil {
		log.Fatal(err)
	}

	<-chromedp.FromContext(ctx).Browser.LostConnection

	log.Infof("Stopping main window\n")
}
