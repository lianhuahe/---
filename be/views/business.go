package views

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sy_spatio-temporal_big_data_platform/dal/db"
	dbmodel "sy_spatio-temporal_big_data_platform/db_model"
	"time"

	"github.com/cloudwego/hertz/cmd/hz/util/logs"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	DATASET_EXAMPLE_PATH    = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/sample/T_DRIVE_SMALL.zip"
	TASK_PARAM_EXAMPLE_PATH = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/sample/config.json"
	ADMIN_FRONT_HTML_PATH   = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/fe/dist/"
	DATASET_PATH            = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/dataset/"
	SCRIPTS_PATH            = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/scripts/"
	// TASK_MODEL_PATH = "/home/lianhua/goproject/src/shenyang/sy_spatio-temporal_big_data_platform/be/scripts/" LIBCITY_PATH + os.sep + 'libcity' + os.sep + 'config' + os.sep + 'model'
)

// 数据集状态枚举类
const (
	ERROR               = -1
	PROCESSING          = 0
	PROCESSING_COMPLETE = 1 // 用来表示完成geojson文件生成
	SUCCESS             = 2 // 用来表示完成html文件生成
	UN_PROCESS          = 3 // 已经生成geojson但是没有进行可视化处理
	CHECK               = 4 // 正在生成geojson（检查是否可以进行可视化处理）
	SUCCESS_stat        = 5 // 用来表示完成html文件生成
)

type SearchFilesReq struct {
	FileName   string    `json:"file_name" form:"file_name"`
	CreatTime  time.Time `json:"create_time" form:"create_time"`
	Visibility string    `json:"visibility" form:"visibility"`
	Creator    string    `json:"creator" form:"creator"`
	Page       int32     `json:"page" form:"page"`
	Size       int32     `json:"size" form:"size"`
}

type SearchFilesResp struct {
	Results []*dbmodel.File `json:"results"`
	Count   int             `json:"count"`
}

func SearchFiles(c context.Context, ctx *app.RequestContext) {
	req := &SearchFilesReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("SearchFiles Bind err: %v", err)
		return
	}

	files, err := db.SearchFiles(c, req.FileName, req.CreatTime, req.Visibility, req.Creator, req.Page, req.Size)
	if err != nil {
		logs.Error("SearchFiles db SearchFiles err: %v", err)
		return
	}

	ctx.JSON(consts.StatusOK, utils.H{"data": &SearchFilesResp{Results: files, Count: len(files)}})
}

func Download(c context.Context, ctx *app.RequestContext) {
	ctx.File(DATASET_EXAMPLE_PATH)
}

type UploadFileReq struct {
	IsPublic bool   `json:"isPublic" form:"isPublic"`
	DataSet  []byte `json:"dataset" form:"dataset"`
}

func UploadFile(c context.Context, ctx *app.RequestContext) {
	req := &UploadFileReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("SearchFiles Bind err: %v", err)
		return
	}

	//allowedSuffix := []string{".geo", ".usr", ".rel", ".dyna", ".ext", ".json", ".grid", ".gridod", ".od"}
	file, _ := ctx.FormFile("dataset")
	fileSize := file.Size
	fileName := strings.Split(file.Filename, ".")[0]
	extractPath := fmt.Sprintf("%s%s", DATASET_PATH, fileName)
	os.Mkdir(extractPath, os.ModePerm)
	zipPath := fmt.Sprintf("%s%s", DATASET_PATH, file.Filename)
	err = ctx.SaveUploadedFile(file, zipPath)
	fmt.Println("err1: ", err)
	err = db.SaveFile(c, fileName, file.Filename, zipPath, fileSize, 3, extractPath, CHECK, 1)
	if err != nil {
		logs.Error("UploadFile db SaveFile err: %v", err)
		ctx.JSON(consts.StatusInternalServerError, utils.H{"data": ""})
		return
	}
	_, err = os.Create(ADMIN_FRONT_HTML_PATH + fileName + ".html")
	fmt.Println("err2: ", err)
	go ExecuteGeojsonThread(c, zipPath, extractPath, fileName)
}

func UnZip(zipPath string, extractPath string) {
	dst := extractPath
	archive, err := zip.OpenReader(zipPath)
	if err != nil {
		panic(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)
		fmt.Println("unzipping file ", filePath)

		if !strings.HasPrefix(filePath, filepath.Clean(dst)+string(os.PathSeparator)) {
			fmt.Println("invalid file path")
			return
		}
		if f.FileInfo().IsDir() {
			fmt.Println("creating directory...")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			panic(err)
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			panic(err)
		}

		fileInArchive, err := f.Open()
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			panic(err)
		}

		dstFile.Close()
		fileInArchive.Close()
	}
}

func getOutputDirectly(name string, args ...string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output() // 等到命令执行完, 一次性获取输出
	if err != nil {
		logs.Error("getOutputDirectly exec err: %v", err)
		panic(err)
	}
	output, err = simplifiedchinese.GB18030.NewDecoder().Bytes(output)
	if err != nil {
		logs.Error("getOutputDirectly simplifiedchinese err: %v", err)
		panic(err)
	}
	return string(output)
}

func ExecuteGeojsonThread(ctx context.Context, zipPath string, extractPath string, fileName string) {
	fileViewStatus := UN_PROCESS
	UnZip(zipPath, extractPath)
	fmt.Println("python3", SCRIPTS_PATH+"save_geojson.py", "get_geo_json", fileName, extractPath+"_geo_json")
	res := getOutputDirectly("python3", SCRIPTS_PATH+"save_geojson.py", "get_geo_json", fileName, extractPath+"_geo_json")
	fmt.Println("res: ", res)
	content, _ := os.ReadFile(extractPath + "/status")
	fileFormStatus, _ := strconv.Atoi(string(content))
	if fileFormStatus == PROCESSING_COMPLETE {
		db.UpdateStatusByFileName(ctx, fileFormStatus, fileName)
		logs.Info("%v geojson文件生成完毕", fileName)
	} else {
		db.UpdateStatusByFileName(ctx, fileViewStatus, fileName)
		logs.Info("%v 无法生成geojson文件", fileName)
	}
}

type GenerateGisViewReq struct {
	FileId     int64 `path:"id"`
	Background int   `form:"background"`
}

func GenerateGisView(c context.Context, ctx *app.RequestContext) {
	req := &GenerateGisViewReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("GenerateGisView Bind err: %v", err)
		return
	}

	err = db.UpdateStatusAndBackgroundIdByFileId(c, PROCESSING, int8(req.Background), req.FileId)
	if err != nil {
		logs.Error("GenerateGisView db UpdateStatusAndBackgroundIdByFileId err: %v", err)
		return
	}

	file, err := db.GetFileById(c, req.FileId)
	if err != nil {
		logs.Error("GenerateGisView db GetFileById: %v", err)
		return
	}

	go ExecuteGeoViewThread(c, file.ExtractPath, file.FileName, req.Background)

	ctx.JSON(consts.StatusOK, utils.H{})
}

func ExecuteGeoViewThread(ctx context.Context, extractPath string, fileName string, backgroundId int) {
	fileViewStatus := PROCESSING
	_, err := os.Stat(extractPath + "_geo_json")
	fmt.Println("python3", SCRIPTS_PATH+"save_geojson.py", "transfer_geo_json",
		extractPath+"_geo_json", fileName, fmt.Sprintf("%d", backgroundId))
	if err == nil {
		res := getOutputDirectly("python3", SCRIPTS_PATH+"save_geojson.py", "transfer_geo_json",
			extractPath+"_geo_json", fileName, fmt.Sprintf("%d", backgroundId))
		fmt.Println("res: ", res)
		logs.Info("res: %v", res)
		content, _ := os.ReadFile(extractPath + "/status")
		fmt.Println(content)
		fileViewStatus, _ = strconv.Atoi(string(content))
	} else {
		fmt.Println(123456789)
	}
	if fileViewStatus == SUCCESS || fileViewStatus == SUCCESS_stat {
		logs.Info("%v 数据可视化处理完毕", fileName)
	} else {
		logs.Info("%v 数据可视化处理失败", fileName)
	}
	db.UpdateStatusByFileName(ctx, fileViewStatus, fileName)
}

type GetFileStatusReq struct {
	FileId int64 `path:"id"`
}

func GetFileStatus(c context.Context, ctx *app.RequestContext) {
	req := &GetFileStatusReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("GetFileStatus Bind err: %v", err)
		return
	}

	file, err := db.GetFileById(c, req.FileId)
	if err != nil {
		logs.Error("GetFileStatus db GetFileById: %v", err)
		return
	}
	fmt.Println(file.DatasetStatus)
	if file.DatasetStatus == SUCCESS || file.DatasetStatus == SUCCESS_stat {
		res_data := struct {
			FileName         string
			OriginalFileName string
			DatasetStatus    int8
		}{
			FileName:         file.FileName,
			OriginalFileName: file.FileOriginalName,
			DatasetStatus:    file.DatasetStatus,
		}
		ctx.JSON(consts.StatusOK, utils.H{"data": res_data})
		return
	}
	ctx.JSON(consts.StatusAccepted, utils.H{})
}

type GetGisViewReq struct {
	FileName int64 `path:"file_name"`
}

func GetGisView(c context.Context, ctx *app.RequestContext) {
	req := &GetGisViewReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("GetGisView Bind err: %v", err)
		return
	}

	// file, err := db.GetFileById(c, req.FileId)
	// if err != nil {
	// 	logs.Error("GetFileStatus db GetFileById: %v", err)
	// 	return
	// }
	// (ADMIN_FRONT_HTML_PATH + fileName + ".html")
}

type FileGetAllResp struct {
	Data []*dbmodel.File `json:"data"`
}

func FileGetAll(c context.Context, ctx *app.RequestContext) {
	files, _ := db.GetAllFiles(c)
	ctx.JSON(consts.StatusOK, utils.H{"data": files})
}

type GetTaskModelDictResp struct {
	ResultDict map[string]interface{} `json:"result_dict"`
}

func GetTaskModelDict(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"data": &GetTaskModelDictResp{ResultDict: map[string]interface{}{}}})
}

type CheckTaskExistsReq struct {
	TaskName string `json:"task_name"`
}

type CheckTaskExistsResp struct {
	Id  int64  `json:"id"`
	Msg string `json:"msg"`
}

func CheckTaskExists(c context.Context, ctx *app.RequestContext) {
	req := &CheckTaskExistsReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("CheckTaskExists Bind err: %v", err)
		return
	}

	exests, taskId, err := db.CheckTaskExists(c, req.TaskName)
	if err != nil {
		logs.Error("CheckTaskExists db CheckTaskExists: %v", err)
		return
	}

	if exests {
		ctx.JSON(consts.StatusBadRequest, utils.H{"data": &CheckTaskExistsResp{Id: taskId, Msg: "任务已存在"}})
	} else {
		ctx.JSON(consts.StatusOK, utils.H{"data": &CheckTaskExistsResp{Id: taskId, Msg: ""}})
	}
}

type CreateTaskReq struct {
	TaskName        string `json:"task_name"`
	TaskDescription string `json:"task_description"`
	Task            string `json:"task"`
	Model           string `json:"model"`
	Dataset         string `json:"dataset"`
	SavedModel      bool   `json:"saved_model"`
	Train           bool   `json:"train"`
	MaxEpoch        int    `json:"max_epoch"`
	GPU             bool   `json:"gpu"`
	GPUId           int    `json:"gpu_id"`
	Visibility      int8   `json:"visibility"`
	TaskNameShow    string `json:"task_name_show"`
}

func CreateTask(c context.Context, ctx *app.RequestContext) {
	req := &CreateTaskReq{}
	err := ctx.Bind(req)
	if err != nil {
		logs.Error("CreateTask Bind err: %v", err)
		return
	}

	ctx.JSON(consts.StatusCreated, utils.H{"data": 1, "code": 201})
}

func SearchTasks(c context.Context, ctx *app.RequestContext) {
	ctx.JSON(consts.StatusOK, utils.H{"data": nil})
}
