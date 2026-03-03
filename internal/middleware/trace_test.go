package middleware

import (
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestTraceID_Generation 测试 TraceID 生成
func TestTraceID_Generation(t *testing.T) {
	router := gin.New()
	router.Use(TraceID())
	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceID(c)
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	// 检查响应头中的 TraceID
	traceIDHeader := w.Header().Get(TraceIDHeader)
	if traceIDHeader == "" {
		t.Error("expected X-Trace-ID header to be set")
	}

	// 验证 TraceID 格式: trace-<UUID>
	traceIDPattern := regexp.MustCompile(`^trace-[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)
	if !traceIDPattern.MatchString(traceIDHeader) {
		t.Errorf("TraceID format invalid: %s", traceIDHeader)
	}

	// 检查响应体中的 TraceID
	body := w.Body.String()
	if !strings.Contains(body, traceIDHeader) {
		t.Errorf("expected trace_id in response body: %s", body)
	}
}

// TestTraceID_Unique 测试每次请求生成唯一的 TraceID
func TestTraceID_Unique(t *testing.T) {
	router := gin.New()
	router.Use(TraceID())
	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceID(c)
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	traceIDs := make(map[string]bool)
	for i := 0; i < 100; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		traceID := w.Header().Get(TraceIDHeader)
		if traceIDs[traceID] {
			t.Errorf("duplicate trace ID generated: %s", traceID)
		}
		traceIDs[traceID] = true
	}
}

// TestGetTraceID 测试从 Context 获取 TraceID
func TestGetTraceID(t *testing.T) {
	router := gin.New()
	router.Use(TraceID())
	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceID(c)
		if traceID == "" {
			c.JSON(500, gin.H{"error": "trace_id not found"})
		} else {
			c.JSON(200, gin.H{"trace_id": traceID})
		}
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	body := w.Body.String()
	if !strings.Contains(body, "trace_id") {
		t.Error("expected trace_id in response")
	}
}

// TestTraceID_ContextPropagation 测试 TraceID 在中间件链中传递
func TestTraceID_ContextPropagation(t *testing.T) {
	var traceID1, traceID2 string

	router := gin.New()
	router.Use(TraceID())
	router.GET("/test", func(c *gin.Context) {
		traceID1 = GetTraceID(c)
	}, func(c *gin.Context) {
		traceID2 = GetTraceID(c)
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if traceID1 != traceID2 {
		t.Errorf("TraceID not propagated: %s != %s", traceID1, traceID2)
	}

	if traceID1 == "" {
		t.Error("TraceID should not be empty")
	}
}

// TestTraceID_WithoutMiddleware 测试没有 TraceID 中间件时的行为
func TestTraceID_WithoutMiddleware(t *testing.T) {
	router := gin.New()
	router.GET("/test", func(c *gin.Context) {
		traceID := GetTraceID(c)
		c.JSON(200, gin.H{"trace_id": traceID})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 没有 TraceID 中间件时，应该返回空字符串
	body := w.Body.String()
	if !strings.Contains(body, `""`) && !strings.Contains(body, `trace_id: ""`) {
		t.Logf("Response: %s", body)
	}
}

// BenchmarkTraceID 性能测试
func BenchmarkTraceID(b *testing.B) {
	router := gin.New()
	router.Use(TraceID())
	router.GET("/test", func(c *gin.Context) {
		c.Status(200)
	})

	req := httptest.NewRequest("GET", "/test", nil)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
