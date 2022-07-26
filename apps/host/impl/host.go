package impl

import (
	"context"
	"fmt"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/sqlbuilder"
	"github.com/yuanyp8/cmdb/apps/host"
)

func (i *HostServiceImpl) CreateHost(ctx context.Context, ins *host.Host) (*host.Host, error) {
	l := i.l.With(logger.NewAny("func", "CreateHost"))
	l.Debugf("create host %s", ins.Name)

	// 校验数据合法性
	if err := ins.Validate(); err != nil {
		l.Errorf("Host: %s struct validated error", ins.Name)
		return nil, err
	}

	// 注入默认数据
	ins.InjectDefault()

	// Dao层将数据入库
	if err := i.save(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) QueryHost(ctx context.Context, req *host.QueryHostRequest) (*host.HostSet, error) {
	// 基于sqlbuilder生成query语句
	query := sqlbuilder.NewQuery(queryHostSQL)
	if req.Keywords != "" {
		query.Where("r.name LIKE ? OR r.description LIKE ? OR r.private_ip LIKE ? OR r.public_ip LIKE ?",
			"%"+req.Keywords+"%",
			"%"+req.Keywords+"%",
			req.Keywords+"%",
			req.Keywords+"%",
		)
	}
	query.Limit(req.OffSet(), req.GetPageSize())

	// build 查询语句
	sqlStr, args := query.Build()
	i.l.Debugf("sql: %s, args: %v", sqlStr, args)

	// Prepare
	stmt, err := i.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, fmt.Errorf("prepare query host sql error, %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, args...)
	if err != nil {
		return nil, fmt.Errorf("stmt query error, %s", err)
	}
	defer rows.Close()

	// 初始化需要返回的对象
	set := host.NewHostSet()

	// 迭代查询表里的数据
	for rows.Next() {
		ins := host.NewHost()
		if err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
			&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
			&ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.Id, &ins.CPU, &ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName, &ins.SerialNumber,
		); err != nil {
			return nil, err
		}
		// 将查询到的每条ins数据写入HostSet
		set.Add(ins)
	}

	// total统计
	countStr, args := query.BuildCount()
	i.l.Debugf("count sql: %s, args: %v", countStr, args)

	countStmt, err := i.db.PrepareContext(ctx, countStr)
	if err != nil {
		return nil, fmt.Errorf("prepare count stmt error, %s", err)
	}
	defer countStmt.Close()

	if err := countStmt.QueryRowContext(ctx, args...).Scan(&set.Total); err != nil {
		return nil, fmt.Errorf("count stmt query error, %s", err)
	}

	return set, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	query := sqlbuilder.NewQuery(queryHostSQL)
	query.Where("r.id = ?", req.Id)

	queryStr, args := query.Build()
	i.l.With(logger.NewAny("func", "DescribeHost")).Debugf("describe sql: %s, args: %v", queryStr, args)

	// query stmt,构建prepare
	stmt, err := i.db.PrepareContext(ctx, queryStr)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	ins := host.NewHost()
	if err := stmt.QueryRowContext(ctx, args...).Scan(
		&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
		&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
		&ins.Account, &ins.PublicIP, &ins.PrivateIP,
		&ins.Id, &ins.CPU, &ins.Memory, &ins.GPUAmount, &ins.GPUSpec, &ins.OSType, &ins.OSName, &ins.SerialNumber,
	); err != nil {
		return nil, err
	}
	return ins, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	// 获取已有对象
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(req.Id))
	if err != nil {
		i.l.With(logger.NewAny("func", "UpdateHost")).Errorf("describe host details failed, %s", err)
		return nil, err
	}
	// 根据更新模式来更新主机
	switch req.UpdateMode {
	case host.UPDATE_MODE_PUT:
		if err := ins.Put(req.Host); err != nil {
			return nil, err
		}
	case host.UPDATE_MODE_PATCH:
		// 整个对象的局部更新
		if err := ins.Patch(req.Host); err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("update_mode only requred put/patch")
	}

	// 检查更新的数据是否合法
	if err := ins.Validate(); err != nil {
		return nil, err
	}

	// 更新数据库里的数据
	if err := i.update(ctx, ins); err != nil {
		return nil, err
	}

	// 返回更新后的对象
	return ins, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	ins, err := i.DescribeHost(ctx, host.NewDescribeHostRequestWithId(req.Id))
	if err != nil {
		return nil, err
	}

	if err := i.delete(ctx, ins); err != nil {
		return nil, err
	}
	return ins, nil
}
