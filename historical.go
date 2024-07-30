package omgo

import (
	"context"
)

// Historical retreives the historical weather data for the provided location.
//
// Use `HistoricalOptions` to specify the dates and which metrics to retrieve. The response is
// a Historical struct that will contains the historical weather, all requested hourly predictions
// and all requested daily predictions
func (c Client) Historical(ctx context.Context, loc Location, opts *HistoricalOptions) (*Historical, error) {
	body, err := c.GetHistorical(ctx, loc, opts)
	if err != nil {
		return nil, err
	}

	return ParseHistoricalBody(body)
}
