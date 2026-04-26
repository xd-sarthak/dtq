package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/xd-sarthak/dtq/internal/model"
	"github.com/xd-sarthak/dtq/internal/queue"
	"github.com/xd-sarthak/dtq/internal/storage"

)

type APIHandler struct {
	queue 	*queue.PriorityQueue
	delayed *queue.DelayedQueue
}