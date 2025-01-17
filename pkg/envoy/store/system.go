package store

import (
	"context"

	"github.com/cortezaproject/corteza-server/pkg/envoy"
	"github.com/cortezaproject/corteza-server/pkg/rbac"
	"github.com/cortezaproject/corteza-server/store"
	"github.com/cortezaproject/corteza-server/system/types"
)

type (
	roleFilter        types.RoleFilter
	userFilter        types.UserFilter
	templateFilter    types.TemplateFilter
	applicationFilter types.ApplicationFilter
	settingFilter     types.SettingsFilter
	rbacFilter        struct {
		rbac.RuleFilter
		// This will help us determine what rules for what resources we are able to export
		resourceID map[uint64]bool
	}

	systemStore interface {
		store.Actionlogs
		store.Applications
		store.Attachments
		store.RbacRules
		store.Roles
		store.Settings
		store.Users
		store.Templates
	}

	systemDecoder struct {
		resourceID []uint64
	}
)

func newSystemDecoder() *systemDecoder {
	return &systemDecoder{
		resourceID: make([]uint64, 0, 200),
	}
}

func (d *systemDecoder) decodeRoles(ctx context.Context, s systemStore, ff []*roleFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.RoleSet
	var fn types.RoleFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchRoles(ctx, types.RoleFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newRole(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeUsers(ctx context.Context, s systemStore, ff []*userFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.UserSet
	var fn types.UserFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchUsers(ctx, types.UserFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newUser(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeTemplates(ctx context.Context, s systemStore, ff []*templateFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.TemplateSet
	var fn types.TemplateFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchTemplates(ctx, types.TemplateFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newTemplate(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

func (d *systemDecoder) decodeApplications(ctx context.Context, s systemStore, ff []*applicationFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.ApplicationSet
	var fn types.ApplicationFilter
	var err error

	for _, f := range ff {
		aux := *f

		if aux.Limit == 0 {
			aux.Limit = 1000
		}

		for {
			nn, fn, err = s.SearchApplications(ctx, types.ApplicationFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newApplication(n))
				d.resourceID = append(d.resourceID, n.ID)
			}

			if fn.NextPage != nil {
				aux.PageCursor = fn.NextPage
			} else {
				break
			}
		}
	}

	return &auxRsp{
		mm: mm,
	}
}
func (d *systemDecoder) decodeSettings(ctx context.Context, s systemStore, ff []*settingFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn types.SettingValueSet
	var err error

	for _, f := range ff {
		aux := *f

		for {
			nn, _, err = s.SearchSettings(ctx, types.SettingsFilter(aux))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				mm = append(mm, newSetting(n))
			}
			// mm = append(mm, NewSettings(nn))

			break
		}
	}

	return &auxRsp{
		mm: mm,
	}
}
func (d *systemDecoder) decodeRbac(ctx context.Context, s systemStore, ff []*rbacFilter) *auxRsp {
	mm := make([]envoy.Marshaller, 0, 100)
	if ff == nil {
		return &auxRsp{
			mm: mm,
		}
	}

	var nn rbac.RuleSet
	var err error

	for _, f := range ff {
		aux := *f

		for {
			nn, _, err = s.SearchRbacRules(ctx, rbac.RuleFilter(aux.RuleFilter))
			if err != nil {
				return &auxRsp{
					err: err,
				}
			}

			for _, n := range nn {
				// If not wildcard or is a system rule; check if resource is allowed
				if n.Resource.HasWildcard() || !n.Resource.IsAppendable() {
					mm = append(mm, newRbacRule(n))
				} else {
					id, err := n.Resource.GetID()
					if err != nil {
						return &auxRsp{
							err: err,
						}
					}
					if f.resourceID[id] {
						mm = append(mm, newRbacRule(n))
					}
				}

			}

			break
		}
	}

	return &auxRsp{
		mm: mm,
	}
}

// Roles adds a new RoleFilter
func (df *DecodeFilter) Roles(f *types.RoleFilter) *DecodeFilter {
	if df.roles == nil {
		df.roles = make([]*roleFilter, 0, 1)
	}
	df.roles = append(df.roles, (*roleFilter)(f))
	return df
}

// Users adds a new UserFilter
func (df *DecodeFilter) Users(f *types.UserFilter) *DecodeFilter {
	if df.users == nil {
		df.users = make([]*userFilter, 0, 1)
	}
	df.users = append(df.users, (*userFilter)(f))
	return df
}

// Templates adds a new TemplateFilter
func (df *DecodeFilter) Templates(f *types.TemplateFilter) *DecodeFilter {
	if df.templates == nil {
		df.templates = make([]*templateFilter, 0, 1)
	}
	df.templates = append(df.templates, (*templateFilter)(f))
	return df
}

// Applications adds a new ApplicationFilter
func (df *DecodeFilter) Applications(f *types.ApplicationFilter) *DecodeFilter {
	if df.applications == nil {
		df.applications = make([]*applicationFilter, 0, 1)
	}
	df.applications = append(df.applications, (*applicationFilter)(f))
	return df
}

// Settings adds a new SettingsFilter
func (df *DecodeFilter) Settings(f *types.SettingsFilter) *DecodeFilter {
	if df.settings == nil {
		df.settings = make([]*settingFilter, 0, 1)
	}
	df.settings = append(df.settings, (*settingFilter)(f))
	return df
}

// Rbac adds a new RuleFilter
func (df *DecodeFilter) Rbac(f *rbac.RuleFilter) *DecodeFilter {
	if df.rbac == nil {
		df.rbac = make([]*rbacFilter, 0, 1)
	} else {
		// There can only be a single rbac filter
		// since it makes no sense to have multiple of
		return df
	}

	df.rbac = append(df.rbac, &rbacFilter{RuleFilter: *f})
	return df
}

// allowRbacResource adds a new resource identifier to supported resource rules
func (df *DecodeFilter) allowRbacResource(id ...uint64) {
	if df.rbac == nil || len(df.rbac) == 0 {
		return
	}
	rf := df.rbac[0]

	if rf.resourceID == nil {
		rf.resourceID = make(map[uint64]bool)
	}
	for _, i := range id {
		rf.resourceID[i] = true
	}
}
