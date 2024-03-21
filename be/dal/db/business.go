package db

import (
	"context"
	dbmodel "sy_spatio-temporal_big_data_platform/db_model"
	"time"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
)

func SearchFiles(ctx context.Context, fileName string, createTime time.Time, visibility string, creator string, page int32, size int32) ([]*dbmodel.File, error) {
	files := make([]*dbmodel.File, 0)
	dbExec := db.Debug().Model(&dbmodel.File{}).Offset(page - 1).Limit(size)
	if fileName != "" {
		dbExec = dbExec.Where("file_name = ?", fileName)
	}
	var t time.Time
	if createTime != t {
		dbExec = dbExec.Where("create_time = ?", createTime)
	}
	// if creator==""{
	// 	if visibility!=""&&visibility=="1"{
	// 		dbExec=dbExec.Where("visibility = 1")
	// 	}else{
	// 		dbExec=dbExec.Where("visibility = 1")
	// 	}
	// }
	err := dbExec.Find(&files).Error
	if err != nil {
		logs.Error("db SearchFiles err: %v", err)
		return nil, err
	}
	return files, nil
}

func SaveFile(ctx context.Context, fileName string, fileOriginalName string, filePath string, fileSize int64, creatorId int64, extractPath string, dataSetSatus int, visibility int) error {
	file := &dbmodel.File{
		FileName:         fileName,
		FileOriginalName: fileOriginalName,
		FilePath:         filePath,
		FileSize:         fileSize,
		CreatorId:        creatorId,
		ExtractPath:      extractPath,
		DatasetStatus:    int8(dataSetSatus),
		Visibility:       int8(visibility),
		CreateTime:       time.Now(),
		UpdateTime:       time.Now(),
	}
	err := db.Debug().Model(&dbmodel.File{}).Create(file).Error
	if err != nil {
		logs.Error("db SaveFile err: %v", err)
		return err
	}
	return nil
}

func UpdateStatusByFileName(ctx context.Context, status int, fileName string) error {

	err := db.Debug().Model(&dbmodel.File{}).Where("file_name = ?", fileName).
		Updates(dbmodel.File{DatasetStatus: int8(status), UpdateTime: time.Now()}).Error
	if err != nil {
		logs.Error("db UpdateStatusByFileName err: %v", err)
		return err
	}
	return nil
}

func UpdateStatusAndBackgroundIdByFileId(ctx context.Context, status int, backgroundId int8, fileId int64) error {
	err := db.Debug().Model(&dbmodel.File{}).Where("id = ?", fileId).
		Updates(dbmodel.File{DatasetStatus: int8(status), BackgroundId: backgroundId, UpdateTime: time.Now()}).Error
	if err != nil {
		logs.Error("db UpdateStatusAndBackgroundIdByFileId err: %v", err)
		return err
	}
	return nil
}

func GetAllFiles(ctx context.Context) ([]*dbmodel.File, error) {
	files := make([]*dbmodel.File, 0)
	err := db.Debug().Model(&dbmodel.File{}).Find(&files).Error
	if err != nil {
		logs.Error("db SearchFiles err: %v", err)
		return nil, err
	}
	return files, nil
}

func CheckTaskExists(ctx context.Context, taskName string) (bool, int64, error) {
	task := dbmodel.Task{}
	res := db.Debug().Model(&dbmodel.Task{}).Where("task_name = ?", taskName).Find(&task)
	if res.RowsAffected > 0 {
		return true, task.Id, nil
	} else {
		return false, 0, nil
	}
}

func GetFileById(ctx context.Context, fileId int64) (dbmodel.File, error) {
	file := dbmodel.File{}
	err := db.Debug().Model(&dbmodel.File{}).Where("id = ?", fileId).Find(&file).Error
	if err != nil {
		logs.Error("db GetFileById err: %v", err)
		return file, err
	}
	return file, nil
}
