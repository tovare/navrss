package feed

import (
	"context"
	"testing"
)

//
// export GOOGLE_APPLICATION_CREDENTIALS="NAV Stillinger-3100c6cb4edc.json"

func TestUpdateAndStore(t *testing.T) {
	UpdateAndStore(context.Background(), PubSubMessage{[]byte("hello")})
}
