package plugin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	expectedJSON := `{"parameters": [{"number":123,"digest":"sha256:942ae2dfd73088b54d7151a3c3fd5af038a51c50029bfcfd21f1e650d9579967"},{"number":456,"digest":"sha256:224e68cc69566e5cbbb76034b3c42cd2ed57c1a66720396e1c257794cb7d68c1"}]}`
	token := "0bc57212c3cbbec69d20b34c507284bd300def5b"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		authHeader := r.Header.Get("Authorization")
		if authHeader != "Bearer "+token {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		_, err := w.Write([]byte(expectedJSON))
		if err != nil {
			assert.NoError(t, fmt.Errorf("Error Write %w", err))
		}
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	client, err := NewPluginService("plugin-test", ts.URL, token, 0)
	require.NoError(t, err)

	data, err := client.List(t.Context(), nil)
	require.NoError(t, err)

	var expectedData ServiceResponse
	err = json.Unmarshal([]byte(expectedJSON), &expectedData)
	require.NoError(t, err)
	assert.Equal(t, &expectedData, data)
}
