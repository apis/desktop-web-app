package mainWindow

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/sirupsen/logrus"
)

const windowLeft = 0
const windowTop = 0
const windowWidth = 800
const windowHeight = 600

const jsCode = `
// Disable common Chrome shortcuts like Ctrl+N, Ctrl+T, Ctrl+F etc.
window.addEventListener("keydown", (event) => {
  if (event.ctrlKey) {
    event.stopPropagation();
    event.preventDefault();
  }
});

// Disable Chrome context menu
window.addEventListener("contextmenu", (event) => {
  event.stopPropagation();
  event.preventDefault();
});
`

func Run(host string, port int, uiPrefix string, log *logrus.Logger) {
	log.Infof("Starting main window\n")

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		//chromedp.Flag("app", true),
		//chromedp.Flag("app", "data:text/html,<html><body><script>window.moveTo(580,240);window.resizeTo(800,600);window.location='http://ibm.com';</script></body></html>"),
		//chromedp.Flag("app", "data:text/html,<title>TEST</title>"),
		//chromedp.Flag("app", "data:text/html,<title>App</title><style>html{background: #0000FF;}</style>"),
		chromedp.Flag("app", "data:text/html,<title>App</title>"),
		//chromedp.Flag("app", "data:text/html,<title>&lrm;</title>"),
		//	chromedp.Flag("window-size", fmt.Sprintf("%d,%d", windowWidth, windowHeight)),
		//	chromedp.Flag("window-position", fmt.Sprintf("%d,%d", windowLeft, windowTop)),
		chromedp.Flag("start-maximized", true),
	}

	allocatorContext, _ := chromedp.NewExecAllocator(context.Background(), opts...)

	ctx, cancel := chromedp.NewContext(
		allocatorContext,
		//chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(fmt.Sprintf("http://%s:%d%s", host, port, uiPrefix)),
		chromedp.Evaluate(jsCode, nil),
	}); err != nil {
		log.Fatalf("chromedp.Run() failed: %s\n",
			fmt.Errorf("%w", err))
	}

	<-chromedp.FromContext(ctx).Browser.LostConnection

	log.Infof("Stopping main window\n")
}
