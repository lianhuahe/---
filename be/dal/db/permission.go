package db

import (
	"context"
	"fmt"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
)

func GetAllPermissions(ctx context.Context, accountId int32) ([]string, error) {
	permissionCodes := make([]string, 0)
	query_sql := `
                SELECT tp.code
                FROM tb_account a
                    LEFT JOIN tb_account_roles tar ON a.id = tar.account_id
                    LEFT JOIN tb_role_permissions trp ON tar.role_id = trp.role_id
                    LEFT JOIN tb_permission tp ON trp.permission_id = tp.id
                WHERE a.id = ?
                `
	err := db.Debug().Raw(query_sql, accountId).Scan(&permissionCodes).Error
	if err != nil {
		logs.Error("db GetAllPermissions err: %v", err)
		return nil, err
	}
	fmt.Println(permissionCodes)
	return permissionCodes, nil
}
