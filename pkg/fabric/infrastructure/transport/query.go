package transport

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type fabricResponse struct {
	FabricId string  `json:"fabric_id"`
	Name     string  `json:"name"`
	Amount   float32 `json:"amount"`
	Cost     float32 `json:"cost"`
}

type fabricsResponse struct {
	Fabrics []fabricResponse `json:"fabrics"`
}

func (s *server) getFabrics(w http.ResponseWriter, _ *http.Request) {
	fabrics, err := s.fqs.GetFabrics()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var fabricsResponseList []fabricResponse
	for _, fabric := range fabrics {
		fabricsResponseList = append(fabricsResponseList, fabricResponse{
			FabricId: fabric.ID,
			Name:     fabric.Name,
			Amount:   fabric.Amount,
			Cost:     fabric.Cost,
		})
	}

	jsonFabrics, err := json.Marshal(fabricsResponse{
		Fabrics: fabricsResponseList,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if _, err = io.WriteString(w, string(jsonFabrics)); err != nil {
		log.WithField("err", err).Error("write response error")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
