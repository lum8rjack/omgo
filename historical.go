package omgo

import "context"

func (c Client) Historical(ctx context.Context, loc Location, opts *HistoricalOptions) (*Historical, error) {
	body, err := c.GetHistocial(ctx, loc, opts)
	if err != nil {
		return nil, err
	}

	return ParseHistoricalBody(body)
}
