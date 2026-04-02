package updater

import (
	"testing"

	"rime-wanxiang-updater/internal/types"
)

func TestResetPlannedUpdateInfoClearsOnlyComponentsMarkedForUpdate(t *testing.T) {
	combined := &CombinedUpdater{
		SchemeUpdater: &SchemeUpdater{UpdateInfo: &types.UpdateInfo{Tag: "scheme-stale"}},
		DictUpdater:   &DictUpdater{UpdateInfo: &types.UpdateInfo{Tag: "dict-stale"}},
		ModelUpdater:  &ModelUpdater{UpdateInfo: &types.UpdateInfo{Tag: "model-stale"}},
	}

	combined.resetPlannedUpdateInfo(true, false, true)

	if combined.SchemeUpdater.UpdateInfo != nil {
		t.Fatal("SchemeUpdater.UpdateInfo != nil, want nil when scheme refresh is required")
	}
	if combined.DictUpdater.UpdateInfo == nil {
		t.Fatal("DictUpdater.UpdateInfo = nil, want original value when dict refresh is not required")
	}
	if combined.ModelUpdater.UpdateInfo != nil {
		t.Fatal("ModelUpdater.UpdateInfo != nil, want nil when model refresh is required")
	}
}
