package orm

import (
	"encoding/json"
	"time"

	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/pkg/rta/ams/exp/api"
	"gorm.io/gorm"
)

// GetAllRtaAccount 从数据库中获取所有的Rta账户的信息
func GetAllRtaAccount(db *gorm.DB, option *StatementOption) ([]*types.RtaAccount, error) {
	var rtaAccounts []*types.RtaAccount
	if err := getDBWithOption(db, option).Find(&rtaAccounts).Error; err != nil {
		return nil, err
	}
	return rtaAccounts, nil
}

// GetRtaAccountById 根据id从数据库中获得Rta账户的信息
func GetRtaAccountById(db *gorm.DB, id uint64) (*types.RtaAccount, error) {
	rtaAccount := &types.RtaAccount{}
	if err := db.Take(rtaAccount, id).Error; err != nil {
		return nil, err
	}
	return rtaAccount, nil
}

// LoadAccountRtaExp 加载Rta账户管理的Rta实验
func LoadAccountRtaExp(db *gorm.DB, account *types.RtaAccount) error {
	return db.Model(account).Association("RtaExp").Find(&account.RtaExp)
}

// UpsertRtaAccount 插入或者更新Rta账户的信息
func UpsertRtaAccount(db *gorm.DB, rtaAccount *types.RtaAccount) error {
	if rtaAccount.ID == 0 {
		return insertOrRecover(db, rtaAccount)
	}
	return db.Updates(rtaAccount).Error
}

// DeleteRtaAccountById 根据id从数据库中软删除Rta账户
func DeleteRtaAccountById(db *gorm.DB, id uint64) error {
	return db.Delete(&types.RtaAccount{}, id).Error
}

// SyncRtaAccountExp 同步RTA账号绑定的实验
func SyncRtaAccountExp(db *gorm.DB, account *types.RtaAccount) error {
	currentRtaExp, err := GetRtaExp(account)
	if err != nil {
		return err
	}

	if err = LoadAccountRtaExp(db, account); err != nil {
		return err
	}

	if len(account.RtaExp) == 0 {
		return db.Model(account).Update("RtaExp", currentRtaExp).Error
	} else {
		return mergeRtaExp(db, account, currentRtaExp)
	}
}

func mergeRtaExp(db *gorm.DB, account *types.RtaAccount, current []*types.RtaExp) error {
	original := account.RtaExp
	for _, exp := range current {
		exp.RtaAccountID = account.ID
	}

	currentIdx := make(map[string]struct{})
	for _, cExp := range current {
		currentIdx[cExp.ID] = struct{}{}
	}

	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	originalIdx := make(map[string]struct{})
	for _, oExp := range original {
		originalIdx[oExp.ID] = struct{}{}

		if _, ok := currentIdx[oExp.ID]; !ok {
			if err := DeleteRtaExpById(tx, oExp.ID); err != nil {
				return err
			}
		}
	}

	for _, cExp := range current {
		if _, ok := originalIdx[cExp.ID]; ok {
			if err := UpdateRtaExp(tx, cExp); err != nil {
				return err
			}
		} else {
			if err := InsertRtaExp(tx, cExp); err != nil {
				return err
			}
		}
	}

	tx.Commit()
	if err := tx.Error; err != nil {
		tx.Rollback()
		return err
	}

	account.RtaExp = current
	return nil
}

// GetRtaExp 获取当前账号下的实验
func GetRtaExp(account *types.RtaAccount) ([]*types.RtaExp, error) {
	// 当前只支持ams
	return GetAmsRtaExp(account)
}

// GetAmsRtaExp 获取当前账号AMS的实验
func GetAmsRtaExp(account *types.RtaAccount) ([]*types.RtaExp, error) {
	req := &api.RtaExpListRequest{
		RtaExpRequestBase: &api.RtaExpRequestBase{
			RtaID: account.RtaID,
			Token: account.Token,
		},
	}
	resp, err := api.GetRtaExpList(req)
	if err != nil {
		return nil, err
	}

	ret := make([]*types.RtaExp, 0, len(resp.Data))
	for _, exp := range resp.Data {
		if dbExp, err := convertAmsRtaExp(exp); err != nil {
			return nil, err
		} else {
			ret = append(ret, dbExp)
		}
	}

	return ret, nil
}

var amsRtaExpStatusMap = map[int]types.RtaExpStatus{
	1: types.RtaExpValid,
	2: types.RtaExpPause,
	3: types.RtaExpExpire,
}

func convertAmsRtaExp(exp *api.RtaExp) (*types.RtaExp, error) {
	prop := make(map[string]interface{})
	prop["site_set"] = exp.SiteSet
	propBytes, _ := json.Marshal(prop)

	return &types.RtaExp{
		ID:             api.FormatAmsExpId(exp.ExpId),
		FlowRate:       float64(exp.FlowRate),
		ExpirationTime: time.Time(exp.EndTime),
		JsonProperty:   propBytes,
		Status:         amsRtaExpStatusMap[exp.Status],
	}, nil
}
