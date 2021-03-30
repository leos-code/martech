package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"sort"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/orm"
	"gorm.io/gorm"
)

func checkExperimentGroup(db *gorm.DB, group *types.ExperimentGroup) error {
	if err := checkRtaAccountId(group); err != nil {
		return err
	}
	if err := checkExperimentStageMetadata(db, group.Draft); err != nil {
		return err
	}
	return nil
}

func checkRtaAccountId(group *types.ExperimentGroup) error {
	if group.RtaAccountID == 0 {
		return fmt.Errorf("rta account id can not be empty")
	}
	return nil
}

func checkExperimentStageMetadata(db *gorm.DB, stage *types.ExperimentStage) error {
	for _, item := range stage.ExperimentItem {
		if item.ExperimentMetadata != nil {
			for _, metaData := range item.ExperimentMetadata {
				_, err := orm.GetExperimentParameterById(db, metaData.ExperimentParameterID)
				if err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("invalid operation, choose an existed parameter")
					} else {
						return err
					}
				}
			}
		}
	}
	return nil
}

func updateGroupAll(db *gorm.DB, group *types.ExperimentGroup) error {
	// 修改草稿的版本
	draft := group.Draft
	latestVersion, e := orm.GetLatestVersionOfGroup(db, group.ID)
	if e != nil {
		return e
	}
	draft.Version = latestVersion + 1
	// 将草稿添加到group的ExperimentStage列表才能全局更新
	group.ExperimentStage = []*types.ExperimentStage{group.Draft}
	// 全局更新
	if err := db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(group).Error; err != nil {
		return err
	}
	// 更新user关联关系
	if err := updateGroupUser(db, group); err != nil {
		return err
	}
	// 更新itemList
	if err := updateStageExperimentItemList(db, draft); err != nil {
		return err
	}
	// 循环每一个item，更新其rtaExp和metaData
	for _, item := range draft.ExperimentItem {
		if err := updateItemRtaExp(db, item); err != nil {
			return err
		}
		if err := updateItemMetadata(db, item); err != nil {
			return err
		}
	}
	return nil
}

func updateGroupUser(db *gorm.DB, group *types.ExperimentGroup) error {
	return db.Model(group).Association("User").Replace(group.User)
}

func updateStageExperimentItemList(db *gorm.DB, stage *types.ExperimentStage) error {
	var experimentItemList []*types.ExperimentItem
	if err := db.Model(stage).Association("ExperimentItem").Find(&experimentItemList); err != nil {
		return err
	}
	currentItemList := stage.ExperimentItem
	itemMap := make(map[uint64]*types.ExperimentItem)
	for _, item := range currentItemList {
		if err := orm.UpdateOuterId(db, item); err != nil {
			return err
		}
		itemMap[item.ID] = item
	}
	for _, item := range experimentItemList {
		_, found := itemMap[item.ID]
		if !found {
			if err := db.Model(item).UpdateColumn("experiment_stage_id", nil).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func updateItemRtaExp(db *gorm.DB, item *types.ExperimentItem) error {
	return db.Model(item).Association("RtaExp").Replace(item.RtaExp)
}

func updateItemMetadata(db *gorm.DB, item *types.ExperimentItem) error {
	return db.Model(item).Association("ExperimentMetadata").Replace(item.ExperimentMetadata)
}

func stopCurrentStage(tx *gorm.DB, stage *types.ExperimentStage) error {
	stage.Status = types.ExperimentStageStop
	stage.EndTime = time.Now()
	return orm.UpsertExperimentStage(tx, stage)
}

func insertExperimentGroupAll(db *gorm.DB, group *types.ExperimentGroup) error {
	if err := orm.UpsertExperimentGroup(db, group); err != nil {
		return err
	}
	for _, stage := range group.ExperimentStage {
		for _, item := range stage.ExperimentItem {
			if err := orm.UpdateOuterId(db, item); err != nil {
				return err
			}
		}
	}
	return nil
}

func createExperimentStageFromDraft(db *gorm.DB, draft *types.ExperimentStage) (*types.ExperimentStage, error) {
	newExpStage := &types.ExperimentStage{}
	bytes, _ := json.Marshal(draft)
	_ = json.Unmarshal(bytes, newExpStage)

	newExpStage.ID = 0
	newExpStage.StartTime = time.Now()
	newExpStage.Status = types.ExperimentStageRunning

	for _, expItem := range newExpStage.ExperimentItem {
		expItem.ID = 0
		for _, metadata := range expItem.ExperimentMetadata {
			metadata.ID = 0
		}
	}

	if err := orm.UpsertExperimentStage(db, newExpStage); err != nil {
		return nil, err
	}

	if err := orm.FillExperimentStage(db, newExpStage); err != nil {
		return nil, err
	}
	return newExpStage, nil
}

func unbindExperimentStage(tx *gorm.DB, stage *types.ExperimentStage) error {
	return unbindStageExpItem(tx, stage.ExperimentItem)
}

func bindExperimentStage(tx *gorm.DB, stage *types.ExperimentStage) error {
	return bindStageExpItem(tx, stage.ExperimentItem)
}

func unbindStageExpItem(tx *gorm.DB, itemList []*types.ExperimentItem) error {
	for _, oldItem := range itemList {
		for _, rtaExp := range oldItem.RtaExp {
			if err := unbindExperimentItemAndRtaExp(tx, rtaExp); err != nil {
				return err
			}
		}
	}
	return nil
}

func bindStageExpItem(tx *gorm.DB, itemList []*types.ExperimentItem) error {
	for _, newItem := range itemList {
		for _, rtaExp := range newItem.RtaExp {
			if err := bindExperimentItemAndRtaExp(tx, newItem, rtaExp); err != nil {
				return err
			}
		}
	}
	return nil
}

func unbindExperimentItemAndRtaExp(tx *gorm.DB, rtaExp *types.RtaExp) error {
	return tx.Set("gorm:query_option", "FOR UPDATE").
		Model(rtaExp).
		Where("bind_status = ?", types.RtaExpBusy).
		Updates(map[string]interface{}{"bind_status": types.RtaExpIdle, "bind_experiment_item_id": nil}).
		Error
}

func bindExperimentItemAndRtaExp(tx *gorm.DB, experimentItem *types.ExperimentItem, rtaExp *types.RtaExp) error {
	rtaExp.BindExperimentItemID = experimentItem.ID
	rtaExp.BindStatus = types.RtaExpBusy
	rowsAffected := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("bind_status = ?", types.RtaExpIdle).
		Updates(rtaExp).RowsAffected
	if rowsAffected > 0 {
		return nil
	}
	return errors.New("chosen rta experiments are already used by others")
}

func newExperimentStage(tx *gorm.DB, current *types.ExperimentStage, draft *types.ExperimentStage) error {
	if current != nil {
		if err := unbindExperimentStage(tx, current); err != nil {
			return err
		}

		if err := stopCurrentStage(tx, current); err != nil {
			return err
		}
	}

	target, err := createExperimentStageFromDraft(tx, draft)
	if err != nil {
		return err
	}

	return bindExperimentStage(tx, target)
}

func sortMetadataByParameterId(metadataList []*types.ExperimentMetadata) {
	sort.Slice(metadataList, func(i, j int) bool {
		return metadataList[i].ExperimentParameterID < metadataList[j].ExperimentParameterID
	})
}

func modifyExperimentStage(tx *gorm.DB, current *types.ExperimentStage, draft *types.ExperimentStage) error {
	target := &types.ExperimentStage{}
	bytes, _ := json.Marshal(draft)
	_ = json.Unmarshal(bytes, target)

	target.ID = current.ID
	target.Status = types.ExperimentStageRunning

	return modifyExperimentStageAll(tx, target, current)
}

func modifyExperimentStageAll(tx *gorm.DB, target *types.ExperimentStage, current *types.ExperimentStage) error {
	for i, item := range target.ExperimentItem {
		currentItem := current.ExperimentItem[i]
		item.ID = currentItem.ID

		if err := modifyItemExperimentMetadata(tx, item, currentItem); err != nil {
			return err
		}
		for j, rtaExp := range item.RtaExp {
			rtaExp.ID = currentItem.RtaExp[j].ID
			if err := orm.UpdateRtaExp(tx, rtaExp); err != nil {
				return err
			}
		}
		if err := orm.UpsertExperimentItem(tx, item); err != nil {
			return err
		}
	}
	return orm.UpsertExperimentStage(tx, target)
}

func modifyItemExperimentMetadata(db *gorm.DB,
	targetItem *types.ExperimentItem,
	currentItem *types.ExperimentItem) error {

	targetMetadataList := targetItem.ExperimentMetadata
	currentMetadataList := currentItem.ExperimentMetadata

	sortMetadataByParameterId(targetMetadataList)
	sortMetadataByParameterId(currentMetadataList)

	t, c := 0, 0
	for t < len(targetMetadataList) && c < len(currentMetadataList) {
		targetMetadataList[t].ID = 0
		targetParaId := targetMetadataList[t].ExperimentParameterID
		currentParaId := currentMetadataList[c].ExperimentParameterID
		if targetParaId == currentParaId {
			targetMetadataList[t].ID = currentMetadataList[c].ID
			if err := orm.UpsertExperimentMetadata(db, targetMetadataList[t]); err != nil {
				return err
			}
			t++
			c++
		} else if targetParaId > currentParaId {
			if err := orm.UnbindExperimentMetadata(db, currentMetadataList[c]); err != nil {
				return err
			}
			c++
		} else {
			t++
		}
	}
	for ; t < len(targetMetadataList); t++ {
		targetMetadataList[t].ID = 0
		if err := orm.UpsertExperimentMetadata(db, targetMetadataList[t]); err != nil {
			return err
		}
	}
	for ; c < len(currentMetadataList); c++ {
		if err := orm.UnbindExperimentMetadata(db, currentMetadataList[c]); err != nil {
			return err
		}
	}

	return nil
}

// EditExperimentGroup: 编辑实验
func EditExperimentGroup(db *gorm.DB, group *types.ExperimentGroup) error {
	if group == nil {
		return fmt.Errorf("nil group")
	}

	if group.Draft == nil {
		return fmt.Errorf("there is nil draft in group")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		group.Draft.Status = types.ExperimentStageDraft
		if group.Draft.StartTime.IsZero() {
			group.Draft.StartTime = time.Now()
		}
		if err := checkExperimentGroup(db, group); err != nil {
			return err
		}
		if group.ID == 0 {
			group.Draft.Version = 1
			group.ExperimentStage = []*types.ExperimentStage{group.Draft}
			return insertExperimentGroupAll(tx, group)
		}
		return updateGroupAll(tx, group)
	})
}

// PromoteExperimentGroup: 提升实验
func PromoteExperimentGroup(db *gorm.DB, group *types.ExperimentGroup, c *gin.Context) error {
	return db.Transaction(func(tx *gorm.DB) error {
		draft, err := orm.GetDraftWithAllData(tx, group)
		if err != nil {
			return err
		}
		current, err := orm.GetCurrentWithAllData(tx, group)
		if err != nil {
			return err
		}

		if len(draft.ExperimentItem) == 0 {
			return fmt.Errorf("can't promote a experiment without item")
		}

		fn := modifyExperimentStage
		if !draft.StageEqual(current) {
			fn = newExperimentStage
		}

		if err := fn(tx, current, draft); err != nil {
			return err
		}

		user, err := userSessionGet(c)
		if err != nil {
			return err
		}

		if err := orm.InsertExperimentModifyRecordAfterPromote(tx, group.ID, user.ID); err != nil {
			return err
		}

		return nil
	})
}

// StopExperimentGroup: 停实验
func StopExperimentGroup(db *gorm.DB, group *types.ExperimentGroup, c *gin.Context) error {
	return db.Transaction(func(tx *gorm.DB) error {
		current, err := orm.GetCurrentWithAllData(tx, group)
		if err != nil {
			return err
		}
		currentStageId := current.ID

		if err := unbindExperimentStage(tx, current); err != nil {
			return err
		}

		if err := stopCurrentStage(tx, current); err != nil {
			return err
		}

		user, err := userSessionGet(c)
		if err != nil {
			return err
		}

		if err := orm.InsertExperimentModifyRecordAfterStop(tx, group.ID, currentStageId, user.ID); err != nil {
			return err
		}

		return nil
	})
}
