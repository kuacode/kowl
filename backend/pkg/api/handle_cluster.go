package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/cloudhut/common/rest"
	"github.com/cloudhut/kowl/backend/pkg/kafka"
	"github.com/cloudhut/kowl/backend/pkg/owl"
)

func (api *API) handleDescribeCluster() http.HandlerFunc {
	type response struct {
		ClusterInfo *owl.ClusterInfo `json:"clusterInfo"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		clusterInfo, err := api.OwlSvc.GetClusterInfo(r.Context())
		if err != nil {
			restErr := &rest.Error{
				Err:      err,
				Status:   http.StatusInternalServerError,
				Message:  "Could not describe cluster",
				IsSilent: false,
			}
			rest.SendRESTError(w, r, api.Logger, restErr)
			return
		}

		response := response{
			ClusterInfo: clusterInfo,
		}
		rest.SendResponse(w, r, api.Logger, http.StatusOK, response)
	}
}

const clusterDataFile = "./cluster.json"

var curretCluster = ""

type clusterList struct {
	Selected string        `json:"selected"`
	Data     []clusterItem `json:"data"`
}

// 新增
type clusterItem struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Brokers  []string `json:"brokers"`
	Selected bool     `json:"selected"`
}

func clusterListFromFile() (clusterList, error) {
	var clusters clusterList
	data, err := os.ReadFile(clusterDataFile)
	if err != nil {
		return clusters, err
	}
	err = json.Unmarshal(data, &clusters)
	if err != nil {
		return clusters, err
	}
	return clusters, nil
}

func (api *API) handleClusterListData() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clusters, err := clusterListFromFile()
		if err == nil {
			clusters.Selected = curretCluster
			rest.SendResponse(w, r, api.Logger, http.StatusOK, clusters)
		} else {
			restErr := &rest.Error{
				Status:  http.StatusInternalServerError,
				Message: err.Error(),
			}
			rest.SendRESTError(w, r, api.Logger, restErr)
		}
	}
}

func (api *API) handleClusterChange() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		if api.kafkaServiceMap[id] == nil {
			clusters, err := clusterListFromFile()
			if err != nil {
				restErr := &rest.Error{
					Status:  http.StatusNotFound,
					Message: err.Error(),
				}
				rest.SendRESTError(w, r, api.Logger, restErr)
				return
			}
			// Find cluster by id
			for _, v := range clusters.Data {
				if v.Id == id {
					// Update config
					api.Cfg.Kafka.Id = id
					api.Cfg.Kafka.Brokers = v.Brokers
					break
				}
			}
			if api.Cfg.Kafka.Id != id {
				restErr := &rest.Error{
					Status:  http.StatusNotFound,
					Message: "kafka cluster not found",
				}
				rest.SendRESTError(w, r, api.Logger, restErr)
				return
			}
			// Create new kafka service
			ks, err := kafka.NewService(api.Cfg.Kafka, api.Logger, api.Cfg.MetricsNamespace)
			if err != nil {
				restErr := &rest.Error{
					Status:  http.StatusNotFound,
					Message: err.Error(),
				}
				rest.SendRESTError(w, r, api.Logger, restErr)
				return
			}
			api.kafkaServiceMap[id] = ks
		}
		curretCluster = id
		api.KafkaSvc = api.kafkaServiceMap[id]
		api.OwlSvc.UpdateKafkaService(api.kafkaServiceMap[id])
		rest.SendResponse(w, r, api.Logger, http.StatusOK, "OK")
	}
}

// func (api *API) handleAddCluster() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// 1. Parse and validate request
// 		var req clusterItem
// 		restErr := rest.Decode(w, r, &req)
// 		if restErr != nil {
// 			rest.SendRESTError(w, r, api.Logger, restErr)
// 			return
// 		}
// 		f, err := os.OpenFile("./cluster.json", os.O_RDWR|os.O_CREATE, 0644)
// 		if err != nil {
// 			restErr = &rest.Error{
// 				Status:  http.StatusNotFound,
// 				Message: err.Error(),
// 			}
// 			rest.SendRESTError(w, r, api.Logger, restErr)
// 			return
// 		}
// 		io.ReadAll(f)
// 	}
// }
