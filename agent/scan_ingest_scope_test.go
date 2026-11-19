package agent

import (
	"testing"

	"github.com/Autumn-27/artex/db"
)

func TestRequireTaskAssetAllowlist(t *testing.T) {
	d := testDB(t)
	t.Cleanup(func() { _ = d.Close() })
	task, err := d.CreateTaskWithOptions("allowlist fixture", "test only", db.TaskCreateOptions{})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.DeleteTask(task.ID) })
	tsx := NewToolSet(d.Exploration(task.ExplorationID), "worker-1")
	tsx.SetTaskID(task.ID)
	tsx.SetAssetStore(d.Assets(), d.Companies())
	if err := tsx.requireTaskAssetAllowlist(); err == nil {
		t.Fatal("ingestion allowed without a task allow rule")
	}
	if _, err := d.Assets().CreateTaskInterceptRule(task.ID, "allow", "exact_domain", "allowed.example.com", "disabled fixture", false); err != nil {
		t.Fatal(err)
	}
	if err := tsx.requireTaskAssetAllowlist(); err == nil {
		t.Fatal("ingestion allowed with only disabled allow rules")
	}
	if _, err := d.Assets().CreateTaskInterceptRule(task.ID, "allow", "exact_domain", "allowed.example.com", "enabled fixture", true); err != nil {
		t.Fatal(err)
	}
	if err := tsx.requireTaskAssetAllowlist(); err != nil {
		t.Fatalf("ingestion rejected with an enabled allow rule: %v", err)
	}
}
