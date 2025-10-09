package pdf

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func (p *PDF) Generate(cookies []*http.Cookie, reqUrl string) error {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	ctx, cancelTimeout := context.WithTimeout(ctx, 60*time.Second)
	defer cancelTimeout()

	var err error
	p.buffer, err = getBuffer(ctx, cookies, reqUrl)
	if err != nil {
		return err
	}

	return save(p.buffer, p.path)
}

func getBuffer(ctx context.Context, cookies []*http.Cookie, reqUrl string) ([]byte, error) {
	var buffer []byte
	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}

	return buffer, chromedp.Run(ctx,
		network.Enable(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			for _, c := range cookies {
				if c.Domain == "" {
					c.Domain = u.Hostname()
				}
				if err := network.SetCookie(c.Name, c.Value).
					WithDomain(c.Domain).
					WithPath(c.Path).
					Do(ctx); err != nil {
					return err
				}
			}
			return nil
		}),
		chromedp.Navigate(reqUrl),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buffer, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperWidth(8.27).   // Format A4
				WithPaperHeight(11.67). // Format A4
				WithMarginTop(0.4).
				WithMarginBottom(0.4).
				WithMarginLeft(0.4).
				WithMarginRight(0.4).
				Do(ctx)
			return err
		}),
	)
}

func save(buffer []byte, path string) error {
	return os.WriteFile(path, buffer, 0644)
}
