package integration

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
	"xiaoweishu/internal/integration/startup"
)

// 预期输入
type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// 预期输出
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

// 测试套件
type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
}

func (s *ArticleTestSuite) SetupSuite() {
	// 在所有测试执行之前，初始化一些内容
	s.server = gin.Default()
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

func (s *ArticleTestSuite) TestArticleLog() {
	s.T().Log("这是测试套件")
}

func (s *ArticleTestSuite) TestArticleEdit() {
	t := s.T()
	testCases := []struct {
		name string
		// 需要考虑准备数据，验证数据。
		before func(t *testing.T)
		// 需要考虑数据库的数据对不对，redis的数据对不对
		after func(t *testing.T)

		// 预期输入
		article Article

		// HTTP 响应码
		wantCode int
		wantRes  Result[int64]
	}{
		{
			name:   "新建帖子-保存成功",
			before: func(t *testing.T) {},
			after:  func(t *testing.T) {},
			article: Article{
				Title:   "标题",
				Content: "内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 1,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 构造请求
			// 执行
			// 验证结果
			tc.before(t)
			reqBody, err := json.Marshal(tc.article)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != http.StatusOK {
				return
			}
			var res Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&res)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, res)
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestArticlePublish() {

}
