package rest

import (
	"context"
	"net/url"
	"strconv"

	"github.com/amiwrpremium/go-thalex/types"
)

// NotificationsInbox retrieves inbox notifications.
func (c *Client) NotificationsInbox(ctx context.Context, limit *int) (types.NotificationsResult, error) {
	q := url.Values{}
	if limit != nil {
		q.Set("limit", strconv.Itoa(*limit))
	}
	var result types.NotificationsResult
	err := c.transport.DoPrivateGET(ctx, "/private/notifications_inbox", q, &result)
	return result, err
}

// MarkNotificationAsRead marks a notification as read or unread.
func (c *Client) MarkNotificationAsRead(ctx context.Context, notificationID string, read bool) error {
	body := struct {
		NotificationID string `json:"notification_id"`
		Read           bool   `json:"read"`
	}{NotificationID: notificationID, Read: read}
	return c.transport.DoPrivatePOST(ctx, "/private/mark_inbox_notification_as_read", body, nil)
}
