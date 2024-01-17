package driver

import (
	"context"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

type Chrome struct {
	ua      string
	options []chromedp.ExecAllocatorOption
	cookie  string
	timeout time.Duration
}

func NewChrome(headless bool, ua string, timeout time.Duration) *Chrome {
	options := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("ignore-certificate-errors", "1"),
		chromedp.Flag("enable-automation", false),
		chromedp.UserAgent(ua),
	)
	return &Chrome{
		options: options,
		ua:      ua,
		timeout: timeout,
	}
}

func (c *Chrome) Close() {

}

func (c *Chrome) Browse(url string, cookie string) error {
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), c.options...)
	defer cancel()
	// create context
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, c.timeout)
	defer cancel()

	err := chromedp.Run(ctx,
		//设置webdriver检测反爬
		chromedp.ActionFunc(func(cxt context.Context) error {
			_, err := page.AddScriptToEvaluateOnNewDocument("Object.defineProperty(navigator, 'webdriver', { get: () => false, });").Do(cxt)
			return err
		}),
		chromedp.Navigate(url),
		chromedp.WaitReady("#content"),
		chromedp.SendKeys("//*[@id=\"kw\"]", "hello world"),
		chromedp.Submit("//*[@id=\"su\"]"),
		chromedp.Sleep(time.Second),
	)
	return err

}
