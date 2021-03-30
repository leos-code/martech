package updater

import (
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
	"time"
)

func TestManager(t *testing.T) {
	db, _ := orm.New(&orm.Option{
		Type: orm.DBTypeSQLite,
		DSN:  filepath.Join(t.TempDir(), "sqlite"),
	})
	db.Debug()
	_ = orm.Setup(db)

	stage := &types.ExperimentStage{
		Status:    types.ExperimentStageRunning,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(10 * time.Second),
	}

	rtaExp := &types.RtaExp{
		ID:             "rtaExp1",
		ExpirationTime: time.Now().Add(2 * time.Second),
		Status:         types.RtaExpValid,
	}

	assert.NoError(t, orm.UpsertExperimentStage(db, stage))
	assert.NoError(t, orm.InsertRtaExp(db, rtaExp))

	manager := NewTableUpdateHelperManager()

	expStageSync := NewExperimentStageTableSnapshotSync()
	expStageHelper := NewTableUpdateHelper(expStageSync, db)

	rtaExpSync := NewRtaExpTableSnapshotSync()
	rtaExpHelper := NewTableUpdateHelper(rtaExpSync, db)

	manager.AppendTableUpdateHelper(expStageHelper, rtaExpHelper)

	assert.NoError(t, manager.Run())

	time.Sleep(2 * time.Second)

	stage.EndTime = time.Now()
	assert.NoError(t, orm.UpsertExperimentStage(db, stage))

	time.Sleep(5 * time.Second)

	stage1, err := orm.GetExperimentStageById(db, stage.ID)
	assert.NoError(t, err)
	assert.Equal(t, types.ExperimentStageStop, stage1.Status)

	rtaExp1, err := orm.GetRtaExpById(db, rtaExp.ID)
	assert.NoError(t, err)
	assert.Equal(t, types.RtaExpExpire, rtaExp1.Status)
}
