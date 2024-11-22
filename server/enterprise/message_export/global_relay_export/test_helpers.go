// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.enterprise for license information.

package global_relay_export

import (
	"net/mail"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func AssertHeaderContains(t *testing.T, msg string, expected map[string]string) {
	m, err := mail.ReadMessage(strings.NewReader(msg))
	require.NoError(t, err)
	header := m.Header

	for k, v := range expected {
		assert.Equal(t, v, header.Get(k))
	}
}
