package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-microservice/internal/utils"
	"time"
)

type IntegrationService struct {
	storage *utils.MinioClient
}

func NewIntegrationService(storage *utils.MinioClient) *IntegrationService {
	return &IntegrationService{storage: storage}
}

func (s *IntegrationService) SaveAudit(event utils.AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	objectName := fmt.Sprintf(
		"audit/%s/%d.json",
		time.Now().Format("2006-01-02"),
		time.Now().UnixNano(),
	)

	return s.storage.PutObject(context.Background(), objectName, data)
}
