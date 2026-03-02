package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"xiaoweishu/internal/integration/startup"
	"xiaoweishu/internal/repository/dao"
	ijwt "xiaoweishu/internal/web/jwt"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// 预期输入
type Article struct {
	Id      int64  `json:"id"`
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
	db     *gorm.DB
}

func (s *ArticleTestSuite) SetupSuite() {
	// 在所有测试执行之前，初始化一些内容
	s.server = gin.Default()
	s.server.Use(func(ctx *gin.Context) {
		ctx.Set("claims", &ijwt.UserClaims{
			Uid: 123,
		})
	})
	s.db = startup.IniTesttDB()
	artHdl := startup.InitArticleHandler()
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TearDownTest() {
	// 清空所有数据，并且自增主键恢复到1
	s.db.Exec("TRUNCATE TABLE articles")
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
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 1).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Ctime > 0)
				assert.True(t, art.Utime > 0)
				art.Ctime = 0
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       1,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 123,
					Ctime:    0,
					Utime:    0,
				}, art)
			},
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
		{
			name: "更新帖子-更新成功",
			before: func(t *testing.T) {
				// 提前准备数据
				err := s.db.Create(&dao.Article{
					Id:       2,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 123,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 2).First(&art).Error
				assert.NoError(t, err)
				assert.True(t, art.Utime > 234)
				art.Utime = 0
				assert.Equal(t, dao.Article{
					Id:       2,
					Title:    "新的标题",
					Content:  "新的内容",
					AuthorId: 123,
					Ctime:    123,
					Utime:    0,
				}, art)
			},
			article: Article{
				Id:      2,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Msg:  "OK",
				Data: 2,
			},
		},
		{
			name: "修改他人帖子-修改失败",
			before: func(t *testing.T) {
				// 提前准备数据
				err := s.db.Create(&dao.Article{
					Id:      3,
					Title:   "标题",
					Content: "内容",
					// 测试模拟的用户是 123， 此处帖子的用户是 789
					// 意味着修改他人的数据
					AuthorId: 789,
					Ctime:    123,
					Utime:    234,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				// 验证数据库
				var art dao.Article
				err := s.db.Where("id=?", 3).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, dao.Article{
					Id:       3,
					Title:    "标题",
					Content:  "内容",
					AuthorId: 789,
					Ctime:    123,
					Utime:    234,
				}, art)
			},
			article: Article{
				Id:      3,
				Title:   "新的标题",
				Content: "新的内容",
			},
			wantCode: http.StatusOK,
			wantRes: Result[int64]{
				Code: 5,
				Msg:  "系统错误",
				Data: 0,
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
