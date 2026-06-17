package v2

import (
	auditlog "github.com/larsartmann/samber-do-auditlog"
)

// AuditLogServiceByName returns the first auditlog ServiceInfo matching the name.
// Returns nil if audit logging is not enabled.
func AuditLogServiceByName[T any](cli *CLI[T], name string) *auditlog.ServiceInfo {
	if cli.auditLog == nil {
		return nil
	}

	return cli.auditLog.Report().ServiceByName(name)
}

// AuditLogFailedServices returns all services with invocation or shutdown errors.
// Returns nil if audit logging is not enabled.
func AuditLogFailedServices[T any](cli *CLI[T]) []auditlog.ServiceInfo {
	if cli.auditLog == nil {
		return nil
	}

	return cli.auditLog.Report().FailedServices()
}
