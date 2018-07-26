package repository

import (
	"context"
	"github.com/crusttech/crust/sam/types"
	"github.com/titpetric/factory"
	"time"
)

const (
	sqlAttachmentScope = "deleted_at IS NULL"

	ErrAttachmentNotFound = repositoryError("AttachmentNotFound")
)

type (
	attachment struct{}
)

func Attachment() attachment {
	return attachment{}
}

func (r attachment) FindByID(ctx context.Context, id uint64) (*types.Attachment, error) {
	db := factory.Database.MustGet()

	mod := &types.Attachment{}
	if err := db.GetContext(ctx, mod, "SELECT * FROM attachments WHERE id = ? AND "+sqlAttachmentScope, id); err != nil {
		return nil, err
	} else if mod.ID == 0 {
		return nil, ErrAttachmentNotFound
	} else {
		return mod, nil
	}
}

func (r attachment) FindByRange(ctx context.Context, channelID, fromAttachmentID, toAttachmentID uint64) ([]*types.Attachment, error) {
	db := factory.Database.MustGet()

	sql := `
		SELECT * 
          FROM attachments
         WHERE id BETWEEN ? AND ?
           AND rel_channel = ?
           AND deleted_at IS NULL`

	rval := make([]*types.Attachment, 0)
	if err := db.SelectContext(ctx, &rval, sql, fromAttachmentID, toAttachmentID, channelID); err != nil {
		return nil, err
	}

	return rval, nil
}

func (r attachment) Create(ctx context.Context, mod *types.Attachment) (*types.Attachment, error) {
	db := factory.Database.MustGet()

	mod.SetID(factory.Sonyflake.NextID())
	mod.SetCreatedAt(time.Now())

	if mod.Attachment == nil {
		mod.SetAttachment([]byte("{}"))
	}

	if err := db.Insert("attachments", mod); err != nil {
		return nil, err
	} else {
		return mod, nil
	}
}

func (r attachment) Update(ctx context.Context, mod *types.Attachment) (*types.Attachment, error) {
	db := factory.Database.MustGet()

	now := time.Now()
	mod.SetUpdatedAt(&now)

	if err := db.Replace("attachments", mod); err != nil {
		return nil, err
	} else {
		return mod, nil
	}
}

func (r attachment) Delete(ctx context.Context, id uint64) error {
	return simpleDelete(ctx, "attachments", id)
}