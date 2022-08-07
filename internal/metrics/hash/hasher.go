package hash

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"eridiumdev/yandex-praktikum-go-devops/internal/common/logger"
	"eridiumdev/yandex-praktikum-go-devops/internal/metrics/domain"
)

type hasher struct {
	hashKey []byte
}

func NewHasher(hashKey string) *hasher {
	return &hasher{
		hashKey: []byte(hashKey),
	}
}

func (h *hasher) Hash(ctx context.Context, metric domain.Metric) string {
	var payload string
	switch metric.Type {
	case domain.TypeCounter:
		payload = fmt.Sprintf("%s:counter:%d", metric.Name, metric.Counter)
	case domain.TypeGauge:
		payload = fmt.Sprintf("%s:gauge:%f", metric.Name, metric.Gauge)
	}

	hash := hmac.New(sha256.New, h.hashKey)
	_, err := hash.Write([]byte(payload))
	if err != nil {
		logger.New(ctx).Errorf("[metrics hasher] error when calculating metric hash: %s", err.Error())
		return ""
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func (h *hasher) Check(ctx context.Context, metric domain.Metric, hash string) bool {
	return hash == h.Hash(ctx, metric)
}
