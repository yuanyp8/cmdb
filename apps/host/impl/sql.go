package impl

const (
	InsertResourceSQL = `
	INSERT INTO resource (
		id,
		vendor,
		region,
		create_at,
		expire_at,
		type,
		name,
		description,
		status,
		update_at,
		sync_at,
		account,
		public_ip,
		private_ip
	)
	VALUES
		(?,?,?,?,?,?,?,?,?,?,?,?,?,?);
	`

	InsertDescribeSQL = `
	INSERT INTO host ( resource_id, cpu, memory, gpu_amount, gpu_spec, os_type, os_name, serial_number )
	VALUES
		( ?,?,?,?,?,?,?,? );
	`

	queryHostSQL = `SELECT * FROM resource as r LEFT JOIN host h ON r.id=h.resource_id`

	updateResourceSQL = `UPDATE resource SET vendor=?,region=?,expire_at=?,name=?,description=? WHERE id = ?`
	updateHostSQL     = `UPDATE host SET cpu=?, memory=? WHERE resource_id=?`

	deleteResourceSQL = `DELETE FROM resource WHERE id = ?`
	deleteHostSQL     = `DELETE FROM host WHERE resource_id = ?`
)
