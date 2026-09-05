package domain

type User struct {
	UserID   string `json:"userId" firestore:"userId"`
	Email    string `json:"email" firestore:"email"`
	Name     string `json:"name" firestore:"name"`
	Role     string `json:"role" firestore:"role"`
	FCMToken string `json:"fcmToken,omitempty" firestore:"fcmToken,omitempty"`
}
