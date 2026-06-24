package authctx

// Shared context keys used by auth middleware and HTTP handlers.
const (
	ContextKeyUID       = "uid"
	ContextKeyEmail     = "email"
	ContextKeyFirstName = "first_name"
	ContextKeyLastName  = "last_name"
	ContextKeyIsAdmin   = "is_admin"
	ContextKeySessionID = "session_id"
)
