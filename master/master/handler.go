package master

import (
	"context"
	"crontab/master/dao"
	"crontab/shared/model"
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	DB     *dao.Etcd
	Logger *zap.Logger
}

func (h *Handler) Register() {
	w := Wrapper{Logger: h.Logger}
	http.HandleFunc("/job/save", w.WrapErr(h.SaveJob))
	http.HandleFunc("/job/delete", w.WrapErr(h.DeleteJob))
	http.HandleFunc("/job/list", w.WrapErr(h.GetJobs))
	http.HandleFunc("/job/kill", w.WrapErr(h.KillJob))
}

func (h *Handler) SaveJob(w http.ResponseWriter, req *http.Request) error {
	j := req.PostFormValue("job")

	var job model.Job
	err := json.Unmarshal([]byte(j), &job)
	if err != nil {
		h.Logger.Error("cannot unmarshal post job", zap.Error(err))
		return err
	}

	oldJ, err := h.DB.CreateJob(context.Background(), &job)
	if err != nil {
		h.Logger.Error("cannot create job", zap.Error(err))
		return err
	}

	// 返回正常应答（{"errno": 0, "msg": "", "data": {...}}）
	err = response(w, oldJ)
	if err != nil {
		h.Logger.Error("fail to response", zap.Error(err))
		return err
	}

	return nil
}

func (h *Handler) GetJobs(w http.ResponseWriter, _ *http.Request) error {
	jobs, err := h.DB.GetJobs(context.Background())
	if err != nil {
		h.Logger.Error("cannot get jobs", zap.Error(err))
		return err
	}

	b, err := json.Marshal(respMsg{
		Errno: 0,
		Msg:   "success",
		Data:  jobs,
	})
	if err != nil {
		return fmt.Errorf("cannot marshal response: %v", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write response: %v", err)
	}

	return nil
}

func (h *Handler) DeleteJob(w http.ResponseWriter, req *http.Request) error {
	name := req.PostFormValue("name")

	oldJ, err := h.DB.DeleteJob(context.Background(), name)
	if err != nil {
		return err
	}

	err = response(w, oldJ)
	if err != nil {
		h.Logger.Error("fail to response", zap.Error(err))
		return err
	}

	return nil
}

func (h *Handler) KillJob(w http.ResponseWriter, req *http.Request) error {
	name := req.PostFormValue("name")

	err := h.DB.KillJob(context.Background(), name)
	if err != nil {
		return fmt.Errorf("cannot kill job: %v", err)
	}

	err = response(w, nil)
	if err != nil {
		return fmt.Errorf("fail to response: %v", err)
	}

	return nil
}

type respMsg struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

func response(w http.ResponseWriter, data *model.Job) error {
	b, err := json.Marshal(respMsg{
		Errno: 0,
		Msg:   "success",
		Data:  data,
	})
	if err != nil {
		return fmt.Errorf("cannot marshal response: %v", err)
	}

	_, err = w.Write(b)
	if err != nil {
		return fmt.Errorf("cannot write response: %v", err)
	}

	return nil
}
