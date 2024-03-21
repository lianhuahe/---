package dbmodel

import "time"

type File struct {
	Id               int64     `json:"id"`
	FileName         string    `json:"file_name"`
	FileSize         int64     `json:"file_size"`
	FilePath         string    `json:"file_path"`
	ExtractPath      string    `json:"extract_path"`
	DatasetStatus    int8      `json:"dataset_status"`
	BackgroundId     int8      `json:"background_id"`
	CreatorId        int64     `json:"creator"`
	Visibility       int8      `json:"visibility"`
	FileOriginalName string    `json:"file_original_name"`
	ErrorMessage     string    `json:"error_message"`
	CreateTime       time.Time `json:"create_time"`
	UpdateTime       time.Time `json:"update_time"`
}

func (t File) TableName() string {
	return "tb_file"
}

type Task struct {
	Id              int64     `json:"id"`
	CreateTime      time.Time `json:"create_time"`
	UpdateTime      time.Time `json:"update_time"`
	TaskName        string    `json:"task_name"`
	TaskDescription string    `json:"task_description"`
	TaskStatus      int8      `json:"task_status"`
	Task            string    `json:"task"`
	Model           string    `json:"model"`
	Dataset         string    `json:"dataset"`
	ConfigFile      string    `json:"config_file"`
	SavedModel      int8      `json:"saved_model"`
	Train           int8      `json:"train"`
	BatchSize       int64     `json:"batch_size"`
	TrainRate       float64   `json:"train_rate"`
	EvalRate        float64   `json:"eval_rate"`
	MaxEpoch        int       `json:"max_epoch"`
	GPU             int8      `json:"gpu"`
	GPUId           int       `json:"gpu_id"`
	CreatorId       int64     `json:"creator_id"`
	LearningRate    float64   `json:"learning_rate"`
	ExecuteTime     time.Time `json:"execute_time"`
	ExecuteEndTime  time.Time `json:"execute_end_time"`
	ExecuteMsg      string    `json:"execute_msg"`
	ExpId           int       `json:"exp_id"`
	LogFileName     string    `json:"log_file_name"`
	Visibility      int8      `json:"visibility"`
	TaskNameShow    string    `json:"task_name_show"`
}

func (t Task) TableName() string {
	return "tb_task"
}
