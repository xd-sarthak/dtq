package api

import (
	"encoding/json"
	"net/http"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

)

type APIHandler struct {
	queue 	*queue.PriorityQueue
	delayed *queue.DelayedQueue
}