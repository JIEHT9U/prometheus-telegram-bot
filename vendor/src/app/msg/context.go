package msg

import (
	"context"
	"errors"
)

type key int

const (
	alert key = 616841354
)

//SetAlertsFromContext return context with Alerts msg
func SetAlertsFromContext(ctx context.Context, a Alerts) context.Context {
	return context.WithValue(ctx, alert, a)
}

//GetAlertsFromContext extract  Alerts in  context
func GetAlertsFromContext(ctx context.Context) (Alerts, error) {
	if alerts, ok := ctx.Value(alert).(Alerts); ok {
		return alerts, nil
	}

	return Alerts{}, errors.New("Error get Img from context")
}
