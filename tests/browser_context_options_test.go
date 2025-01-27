package tests

import (
	_ "embed"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/xk6-browser/common"
)

func TestBrowserContextOptionsDefaultValues(t *testing.T) {
	t.Parallel()

	opts := common.NewBrowserContextOptions()
	assert.False(t, opts.AcceptDownloads)
	assert.False(t, opts.BypassCSP)
	assert.Equal(t, common.ColorSchemeLight, opts.ColorScheme)
	assert.Equal(t, 1.0, opts.DeviceScaleFactor)
	assert.Empty(t, opts.ExtraHTTPHeaders)
	assert.Nil(t, opts.Geolocation)
	assert.False(t, opts.HasTouch)
	assert.Nil(t, opts.HttpCredentials)
	assert.False(t, opts.IgnoreHTTPSErrors)
	assert.False(t, opts.IsMobile)
	assert.True(t, opts.JavaScriptEnabled)
	assert.Equal(t, common.DefaultLocale, opts.Locale)
	assert.False(t, opts.Offline)
	assert.Empty(t, opts.Permissions)
	assert.Equal(t, common.ReducedMotionNoPreference, opts.ReducedMotion)
	assert.Equal(t, &common.Screen{Width: common.DefaultScreenWidth, Height: common.DefaultScreenHeight}, opts.Screen)
	assert.Equal(t, "", opts.TimezoneID)
	assert.Equal(t, "", opts.UserAgent)
	assert.Equal(t, &common.Viewport{Width: common.DefaultScreenWidth, Height: common.DefaultScreenHeight}, opts.Viewport)
}

func TestBrowserContextOptionsDefaultViewport(t *testing.T) {
	p := newTestBrowser(t).NewPage(nil)

	viewportSize := p.ViewportSize()
	assert.Equal(t, float64(common.DefaultScreenWidth), viewportSize["width"])
	assert.Equal(t, float64(common.DefaultScreenHeight), viewportSize["height"])
}

func TestBrowserContextOptionsSetViewport(t *testing.T) {
	tb := newTestBrowser(t)
	bctx, err := tb.NewContext(tb.toGojaValue(struct {
		Viewport common.Viewport `js:"viewport"`
	}{
		Viewport: common.Viewport{
			Width:  800,
			Height: 600,
		},
	}))
	require.NoError(t, err)
	t.Cleanup(bctx.Close)
	p, err := bctx.NewPage()
	require.NoError(t, err)

	viewportSize := p.ViewportSize()
	assert.Equal(t, float64(800), viewportSize["width"])
	assert.Equal(t, float64(600), viewportSize["height"])
}

func TestBrowserContextOptionsExtraHTTPHeaders(t *testing.T) {
	tb := newTestBrowser(t, withHTTPServer())
	bctx, err := tb.NewContext(tb.toGojaValue(struct {
		ExtraHTTPHeaders map[string]string `js:"extraHTTPHeaders"`
	}{
		ExtraHTTPHeaders: map[string]string{
			"Some-Header": "Some-Value",
		},
	}))
	require.NoError(t, err)
	t.Cleanup(bctx.Close)
	p, err := bctx.NewPage()
	require.NoError(t, err)

	err = tb.awaitWithTimeout(time.Second*5, func() error {
		resp, err := p.Goto(tb.URL("/get"), nil)
		if err != nil {
			return err
		}
		require.NotNil(t, resp)
		var body struct{ Headers map[string][]string }
		require.NoError(t, json.Unmarshal(resp.Body().Bytes(), &body))
		h := body.Headers["Some-Header"]
		require.NotEmpty(t, h)
		assert.Equal(t, "Some-Value", h[0])
		return nil
	})
	require.NoError(t, err)
}
