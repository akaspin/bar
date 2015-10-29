package transport_test

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"path"
	"testing"
)

func Test_CloneURL(t *testing.T) {
	ep := "http://example.com/bard/v1"
	u1, err := url.Parse(ep)
	assert.NoError(t, err)

	var u2 *url.URL = new(url.URL)
	*u2 = *u1
	u2.Path = path.Join(u2.Path, "/upload")
	assert.NotEqual(t, u1.String(), u2.String())
}
