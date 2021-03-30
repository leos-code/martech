package handler

import (
	"fmt"
	"math"
	"strings"

	"github.com/tencentad/martech/api/types"
	"github.com/ahmetb/go-linq/v3"
	"github.com/casbin/casbin/v2"
)

const (
	read             = "read"
	write            = "write"
	superAdminRole   = "superadmin"
	superAdminDomain = "superdomain"

	SuperAdminRoleID     = math.MaxInt32
	SuperAdminRoleName   = "超级管理员组"
	SuperAdminDomainID   = math.MaxInt32
	SuperAdminDomainName = "超级管理员域"
)

var (
	ErrEmptyID           = fmt.Errorf("empty id")
	ErrNoWritePermission = fmt.Errorf("no write permission")
	ErrNoAPIPermission   = fmt.Errorf("no api permission")
)

func Diff(source, target []interface{}) (add, remove []interface{}) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

type pData interface {
	GetObject() string
}

type writePermissionChecker struct {
	e       casbin.IEnforcer
	subject string
	domain  string
	item    pData
	take    func(uint64) (pData, error)
}

func (w *writePermissionChecker) toCheck(id ...uint64) error {
	if ok, _ := w.e.Enforce(w.subject, w.domain, w.item.GetObject(), write); !ok {
		return ErrNoWritePermission
	}

	for _, v := range id {
		if v == 0 {
			continue
		}

		toCheck, err := w.take(v)
		if err != nil {
			return err
		}

		if ok, _ := w.e.Enforce(w.subject, w.domain, toCheck.GetObject(), write); !ok {
			return ErrNoWritePermission
		}
	}

	return nil
}

func casbinUserEncode(user *types.User) string {
	return fmt.Sprintf("user_%v", user.ID)
}

func casbinUserDecode(user string) *types.User {
	id := uint64(0)
	_, _ = fmt.Sscanf(user, "user_%v", &id)
	return &types.User{ID: id}
}

func casbinUserIs(user string) bool {
	return strings.HasPrefix(user, "user_")
}

func casbinTenantEncode(tenant *types.Tenant) string {
	if tenant.ID == SuperAdminDomainID {
		return superAdminDomain
	}

	return fmt.Sprintf("tenant_%v", tenant.ID)
}

func casbinTenantDecode(tenant string) *types.Tenant {
	id := uint64(0)
	_, _ = fmt.Sscanf(tenant, "tenant_%v", &id)
	return &types.Tenant{ID: id}
}

func casbinTenantIs(tenant string) bool {
	return strings.HasPrefix(tenant, "tenant_") || tenant == superAdminDomain
}

func casbinRoleEncode(role *types.Role) string {
	if role.ID == SuperAdminRoleID {
		return superAdminRole
	}

	return fmt.Sprintf("role_%v", role.ID)
}

func casbinRoleDecode(role string) *types.Role {
	id := uint64(0)
	_, _ = fmt.Sscanf(role, "role_%v", &id)
	return &types.Role{ID: id}
}

func casbinRoleIs(role string) bool {
	return strings.HasPrefix(role, "role_") || role == superAdminRole
}

func casbinObjectEncode(object *types.Object) string {
	return fmt.Sprintf("object_%v", object.ID)
}

func casbinObjectDecode(object string) *types.Object {
	id := uint64(0)
	_, _ = fmt.Sscanf(object, "object_%v", &id)
	return &types.Object{ID: id}
}

func casbinObjectIs(object string) bool {
	return strings.HasPrefix(object, "object_")
}

func isSuperAdmin(e casbin.IEnforcer, user interface{}) bool {
	v := ""
	switch u := user.(type) {
	case *types.User:
		v = casbinUserEncode(u)
	case string:
		v = u
	default:
		v = ""
	}

	has, _ := e.HasRoleForUser(v, superAdminRole, superAdminDomain)
	return has
}
