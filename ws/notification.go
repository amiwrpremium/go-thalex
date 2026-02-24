package ws

import (
	"context"

	"github.com/amiwrpremium/go-thalex/types"
)

// NotificationsInbox retrieves inbox notifications via WebSocket.
func (ws *Client) NotificationsInbox(ctx context.Context, limit *int) (types.NotificationsResult, error) {
	params := map[string]any{}
	if limit != nil {
		params["limit"] = *limit
	}
	var result types.NotificationsResult
	err := ws.call(ctx, "private/notifications_inbox", params, &result)
	return result, err
}

// MarkNotificationAsRead marks a notification as read via WebSocket.
func (ws *Client) MarkNotificationAsRead(ctx context.Context, notificationID string, read bool) error {
	return ws.callNoResult(ctx, "private/mark_inbox_notification_as_read", map[string]any{
		"notification_id": notificationID,
		"read":            read,
	})
}
