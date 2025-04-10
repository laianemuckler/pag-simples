package notification

type Notification struct {
    UserID  int    `json:"user_id"`
    Status  string `json:"status"`
		Data	 Data   `json:"data"`
}

type Data struct {
		Message string `json:"message"` 
}

type NotificationRequest struct {
	Email string `json:"phone_number"`
	Message     string `json:"message"`
}
