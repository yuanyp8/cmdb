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
	query := sqlbuilder.NewQuery(queryHostSQL).Order("creat_at").Desc().Limit(int64(req.OffSet()), uint(req.PageSize))
	// build 查询语句
	sqlStr, args := query.BuildQuery()
	i.l.Debugf("sql: %s, args: %v", sqlStr, args)

	/* Dao */
	// Prepare
	stmt, err := i.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		return nil, fmt.Errorf("prepare query host sql error, %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, fmt.Errorf("stmt query error, %s", err)
	}

	// 初始化需要返回的对象
	set := host.NewHostSet()

	// 迭代查询表里的数据
	for rows.Next() {
		ins := host.NewHost()
		if err := rows.Scan(
			&ins.Id, &ins.Vendor, &ins.Region, &ins.CreateAt, &ins.ExpireAt,
			&ins.Type, &ins.Name, &ins.Description, &ins.Status, &ins.UpdateAt, &ins.SyncAt,
			&ins.Account, &ins.PublicIP, &ins.PrivateIP,
			&ins.CPU, &ins.Memory, &ins.GPUSpec, &ins.GPUAmount, &ins.OSType, &ins.OSName, &ins.SerialNumber,
		); err != nil {
			return nil, err
		}
		set.Add(ins)
	}

	// total统计

	return nil, nil
}

func (i *HostServiceImpl) DescribeHost(ctx context.Context, req *host.DescribeHostRequest) (*host.Host, error) {
	// TODO
	return nil, nil
}

func (i *HostServiceImpl) UpdateHost(ctx context.Context, req *host.UpdateHostRequest) (*host.Host, error) {
	// TODO
	return nil, nil
}

func (i *HostServiceImpl) DeleteHost(ctx context.Context, req *host.DeleteHostRequest) (*host.Host, error) {
	return nil, nil
}
