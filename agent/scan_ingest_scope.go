package agent

import "errors"

// requireTaskAssetAllowlist prevents scanner-supplied asset candidates from
// entering a task graph unless the task has an explicit, enabled allow rule.
// Asset-level matching still happens in insert_assets.
func (t *ToolSet) requireTaskAssetAllowlist() error {
	if t == nil || t.as == nil || t.taskID <= 0 {
		return errors.New("scanner asset ingestion requires a task context and an enabled allow rule")
	}
	_, allowRules, err := t.as.TaskInterceptRulesSplit(t.taskID)
	if err != nil {
		return errors.New("scanner asset ingestion blocked because the task allow-list could not be verified")
	}
	for _, rule := range allowRules {
		if rule.Enabled {
			return nil
		}
	}
	return errors.New("scanner asset ingestion blocked: add an enabled task allow rule before importing scanner output")
}
