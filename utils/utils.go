package utils

import (
	"HelpingPixl/models"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log/slog"
	"net/http"
)

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}
func ChunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
func FetchToStruct[V models.ResponseStructs](url string) (V, error) {
	json := jsoniter.ConfigCompatibleWithStandardLibrary
	var result V

	slog.Debug("Request: ", url)

	resp, err := http.Get(url)
	if err != nil {
		return result, err
	}
	slog.Debug("Response: ", resp.StatusCode)

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
func FetchResponseCode(url string, method string, body io.Reader) (int, error) {
	slog.Debug("Request: ", url)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}

	return resp.StatusCode, nil
}
