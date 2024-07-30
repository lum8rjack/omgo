package omgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCurrentWeather(t *testing.T) {
	c, err := NewClient()
	require.NoError(t, err)

	loc, err := NewLocation(52.3738, 4.8910) // Amsterdam
	require.NoError(t, err)

	res, err := c.CurrentWeather(context.Background(), loc, nil)
	require.NoError(t, err)
	require.False(t, res.Time.IsZero())
}
