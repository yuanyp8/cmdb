package impl

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/yuanyp8/cmdb/apps/host"
)

func (i *HostServiceImpl) save(ctx context.Context, ins *host.Host) error {
	var err error

	// 把数据写入到cmdb.resource 以及 cmdb.host
	//  一次需要往2个表录入数据, 我们需要2个操作 要么都成功，要么都失败, 事务的逻辑
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start tx error, %s", err)
	}

	// 通过defer处理事务提交的方式
	// 1. 无报错，则Commit事务
	// 2. 有报错，则Rollback事务
	defer func() {
		if err != nil {
			if err := tx.Rollback(); err != nil {
				i.l.With(logger.NewAny("func", "save")).Errorf("rollback error, %s", err)
			}
		} else {
			if err := tx.Commit(); err != nil {
				i.l.With(logger.NewAny("func", "save")).Errorf("commit error, %s", err)
			}
		}
	}()

	// 插入resource表的数据
	rstmt, err := tx.PrepareContext(ctx, InsertResourceSQL)
	if err != nil {
		return err
	}
	defer rstmt.Close()

	_, err = rstmt.ExecContext(ctx,
		ins.Id, ins.Vendor, ins.Region, ins.CreateAt, ins.ExpireAt, ins.Type,
		ins.Name, ins.Description, ins.Status, ins.UpdateAt, ins.SyncAt, ins.Account, ins.PublicIP,
		ins.PrivateIP,
	)
	if err != nil {
		return err
	}

	// 插入describe数据
	dstmt, err := tx.PrepareContext(ctx, InsertDescribeSQL)
	if err != nil {
		return err
	}
	defer dstmt.Close()

	_, err = dstmt.ExecContext(ctx,
		ins.Id, ins.CPU, ins.Memory, ins.GPUAmount, ins.GPUSpec,
		ins.OSType, ins.OSName, ins.SerialNumber,
	)
	if err != nil {
		return err
	}
	return nil
}
