package v2

import "errors"

// Audit-log-related sentinel errors.
var (
	// ErrAuditLogNotEnabled indicates an audit-log command was requested but audit logging is not enabled.
	ErrAuditLogNotEnabled = errors.New("audit log not enabled")

	// ErrInvalidOutputFormat indicates an unsupported audit-log output format was requested.
	ErrInvalidOutputFormat = errors.New("invalid output format")
)
