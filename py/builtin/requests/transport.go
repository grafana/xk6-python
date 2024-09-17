package requests

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.k6.io/k6/js/modules"
	khttp "go.k6.io/k6/js/modules/k6/http"
	"go.k6.io/k6/lib/netext/httpext"
)

type tripware struct {
	vu modules.VU
}

var _ http.RoundTripper = (*tripware)(nil)

func newTransport(vu modules.VU) http.RoundTripper {
	return &tripware{vu: vu}
}

func (t *tripware) toParsedRequest(req *http.Request) (*httpext.ParsedHTTPRequest, error) {
	state := t.vu.State()
	if state == nil {
		return nil, khttp.ErrHTTPForbiddenInInitContext
	}

	u := req.URL.String()

	url, err := httpext.NewURL(u, u)
	if err != nil {
		return nil, err
	}

	preq := &httpext.ParsedHTTPRequest{ // nolint:exhaustruct
		URL:         &url,
		Req:         req,
		Timeout:     httpTimeout,
		Throw:       state.Options.Throw.Bool,
		Redirects:   state.Options.MaxRedirects,
		Cookies:     make(map[string]*httpext.HTTPRequestCookie),
		TagsAndMeta: t.vu.State().Tags.GetCurrentValues(),
	}

	if state.Options.DiscardResponseBodies.Bool {
		preq.ResponseType = httpext.ResponseTypeNone
	} else {
		preq.ResponseType = httpext.ResponseTypeBinary
	}

	if req.Body != nil {
		data, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}

		preq.Body = bytes.NewBuffer(data)

		req.Body.Close()
	}

	preq.Req.Header.Set("User-Agent", state.Options.UserAgent.String)

	return preq, nil
}

func (t *tripware) toResponse(req *http.Request, eresp *httpext.Response) (*http.Response, error) {
	resp := new(http.Response)

	resp.Request = req

	resp.Header = http.Header{}

	for k, v := range eresp.Headers {
		resp.Header.Set(k, v)
	}

	resp.Proto = eresp.Proto
	fields := strings.SplitN(resp.Proto, "/", protoFields)

	if major, err := strconv.Atoi(fields[0]); err == nil {
		resp.ProtoMajor = major
	}

	if len(fields) > 1 {
		if minor, err := strconv.Atoi(fields[1]); err == nil {
			resp.ProtoMinor = minor
		}
	}

	if eresp.Body != nil {
		if data, ok := eresp.Body.([]byte); ok {
			resp.Body = io.NopCloser(bytes.NewBuffer(data))
		}
	}

	resp.Status = eresp.StatusText
	resp.StatusCode = eresp.Status

	return resp, nil
}

func (t *tripware) RoundTrip(req *http.Request) (*http.Response, error) {
	preq, err := t.toParsedRequest(req)
	if err != nil {
		return nil, err
	}

	resp, err := httpext.MakeRequest(t.vu.Context(), t.vu.State(), preq)
	if err != nil {
		return nil, err
	}

	return t.toResponse(req, resp)
}

const (
	httpTimeout = 60 * time.Second
	protoFields = 2
)
