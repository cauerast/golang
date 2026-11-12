package stage6

import (
	"context"
	"net/http"
)

// Fetch performs a GET request to url and returns the status code and body.
//
//   - Use the given client and attach ctx to the request, so cancelling ctx or
//     hitting its deadline aborts the request (the returned error must satisfy
//     errors.Is(err, context.DeadlineExceeded) on a deadline).
//   - Read at most 1 MiB of the body.
//   - A non-2xx status returns the status AND body together with a non-nil error.
//
// TODO: http.NewRequestWithContext, client.Do, and ALWAYS defer resp.Body.Close()
// (otherwise connections leak). Use io.LimitReader for the size cap.
func Fetch(ctx context.Context, client *http.Client, url string) (int, string, error) {
	return 0, "", nil
}

// GetJSON performs a GET request and decodes the JSON response body into out
// (a pointer). Non-2xx statuses return an error without decoding.
//
// TODO: Check the status first, then json.NewDecoder(resp.Body).Decode(out).
func GetJSON(ctx context.Context, client *http.Client, url string, out any) error {
	return nil
}
