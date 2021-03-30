package handler_test

import (
	"fmt"
	"github.com/tencentad/martech/api/types"
	"github.com/tencentad/martech/cmd/web/handler"
	"github.com/tencentad/martech/pkg/orm"
	"github.com/ahmetb/go-linq/v3"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"path/filepath"
	"testing"
)

func TestDiff(t *testing.T) {
	add, remove := handler.Diff([]interface{}{"1", "2"}, []interface{}{"2", "3"})
	assert.ElementsMatch(t, []string{"3"}, add)
	assert.ElementsMatch(t, []string{"1"}, remove)
}

func TestEnum(t *testing.T) {
	a := types.ObjectType("aa")
	t.Log(a)
}

func BenchmarkUserRole(b *testing.B) {
	db, e := userRoleBenchmarkInit(b, 10, 5)
	b.Run("id-1;filter-2;role-3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := id1Filter2Role3(db, e)
			if err != nil {
				b.Error(err)
			}
		}
	})
	b.Run("id-1;role-2;filter-3", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := id1Role2Filter3(db, e)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func userRoleBenchmarkInit(b *testing.B, uc int, rc int) (*gorm.DB, casbin.IEnforcer) {
	db, _ := getTestDB(b)
	_ = orm.Setup(db)

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		b.Fatal(err)
	}

	e, err := casbin.NewSyncedEnforcer("../../../configs/casbin_model.conf", adapter)
	if err != nil {
		b.Fatal(err)
	}

	for i := 1; i <= uc; i++ {
		user := &types.User{
			PhoneNumber: fmt.Sprintf("phone_%v", i),
			Email:       fmt.Sprintf("email_%v", i),
		}
		assert.NoError(b, orm.UpsertUser(db, user))
		for _, v := range []types.LoginType{types.LoginTypeRio, types.LoginTypeQQ} {
			lu := &types.LoginUser{
				OpenID:    fmt.Sprintf("openid_%v_%v", i, v),
				LoginType: v,
				UserID:    user.ID,
			}
			assert.NoError(b, orm.UpsertLoginUser(db, lu))
		}
	}

	for i := 1; i <= rc; i++ {
		role := &types.Role{
			Name:   fmt.Sprintf("role_%v", i),
			Object: fmt.Sprintf("object_%v", i%2),
		}
		assert.NoError(b, orm.UpsertRole(db, role))
	}

	_, _ = e.AddPolicies([][]string{
		{"role_1", "tenant", "object_1", "read"},
		{"role_2", "tenant", "object_0", "read"},
	})

	for i := 1; i <= uc; i++ {
		for j := 0; j*2 < rc; j++ {
			t := "tenant"
			u := fmt.Sprintf("user_%v", i)
			r := fmt.Sprintf("role_%v", j*2+2-(i%2))
			_, _ = e.AddRoleForUserInDomain(u, r, t)
		}
	}

	return db, e
}

func id1Filter2Role3(db *gorm.DB, e casbin.IEnforcer) (interface{}, error) {
	t := "tenant"
	rules := e.GetFilteredGroupingPolicy(2, t)
	var uid []uint64
	var rid []uint64
	for _, rule := range rules {
		uid = append(uid, decode2ID(rule[0]))
		rid = append(rid, decode2ID(rule[1]))
	}

	linq.From(uid).Distinct().ToSlice(&uid)
	linq.From(rid).Distinct().ToSlice(&rid)
	users, err := orm.ListUserByIdWithLoginUser(db, uid)
	if err != nil {
		return nil, err
	}

	roles, err := orm.ListRoleById(db, rid)
	if err != nil {
		return nil, err
	}
	// em := handler.NewEnforceCache(e)
	rm := map[uint64]*types.Role{}
	for _, v := range roles {
		if ok, _ := e.Enforce("user_1", t, v.Object, "read"); ok {
			rm[v.ID] = v
		}
	}

	for _, v := range users {
		rs := e.GetRolesForUserInDomain(fmt.Sprintf("user_%v", v.ID), t)
		for _, r := range rs {
			if role, ok := rm[decode2ID(r)]; ok {
				v.Role = append(v.Role, role)
			}
		}
	}

	return users, nil
}

func id1Role2Filter3(db *gorm.DB, e casbin.IEnforcer) (interface{}, error) {
	t := "tenant"
	rules := e.GetFilteredGroupingPolicy(2, t)
	var uid []uint64
	var rid []uint64
	rsm := map[uint64][]uint64{}
	for _, rule := range rules {
		ui := decode2ID(rule[0])
		ri := decode2ID(rule[1])
		uid = append(uid, ui)
		rid = append(rid, ri)
		rsm[ui] = append(rsm[ui], ri)
	}

	linq.From(uid).Distinct().ToSlice(&uid)
	linq.From(rid).Distinct().ToSlice(&rid)
	users, err := orm.ListUserByIdWithLoginUser(db, uid)
	if err != nil {
		return nil, err
	}

	roles, err := orm.ListRoleById(db, rid)
	if err != nil {
		return nil, err
	}
	rm := map[uint64]*types.Role{}
	for _, v := range roles {
		if ok, _ := e.Enforce("user_1", t, v.Object, "read"); ok {
			rm[v.ID] = v
		}
	}

	for _, v := range users {
		rs := rsm[v.ID]
		for _, r := range rs {
			if role, ok := rm[r]; ok {
				v.Role = append(v.Role, role)
			}
		}
	}

	return users, nil
}

func decode2ID(v string) uint64 {
	id := uint64(0)
	for _, x := range []string{"role", "user"} {
		format := fmt.Sprintf("%v_%%v", x)
		n, _ := fmt.Sscanf(v, format, &id)
		if n == 1 {
			break
		}
	}
	return id
}

func getTestDB(tb testing.TB) (*gorm.DB, error) {
	return orm.New(&orm.Option{
		Type: orm.DBTypeSQLite,
		DSN:  filepath.Join(tb.TempDir(), "sqlite"),
	})
}
