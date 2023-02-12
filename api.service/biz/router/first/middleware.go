// Code generated by hertz generator.

package First

import (
	"api.service/biz/model/api/douyin/core"
	utils2 "api.service/biz/utils"
	"context"
	"errors"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	hertzSentinel "github.com/hertz-contrib/opensergo/sentinel/adapter"
)

func rootMw() []app.HandlerFunc {
	// your code...
	return []app.HandlerFunc{
		hertzSentinel.SentinelServerMiddleware(
			hertzSentinel.WithServerResourceExtractor(func(ctx context.Context, c *app.RequestContext) string {
				return "interaction_service"
			}),
			hertzSentinel.WithServerBlockFallback(func(ctx context.Context, c *app.RequestContext) {
				c.AbortWithStatusJSON(400, utils.H{
					"status_code": 10222,
					"status_msg":  "too many request; the quota used up",
				})
			}),
		),
		func(ctx context.Context, c *app.RequestContext) {
			hlog.Infof("entry %s", c.FullPath())
		},
		//鉴权token，存储userId
		func(ctx context.Context, c *app.RequestContext) {
			//token不填的接口
			///douyin/user/register/、/douyin/user/login/
			//token选填的接口
			///douyin/feed/
			//token在的位置时form
			///douyin/publish/action
			path := string(c.Request.URI().Path())
			if path == "/douyin/user/register/" || path == "/douyin/user/login/" || path == "/douyin/feed/" {
				return
			}
			var tokenString string
			if path == "/douyin/publish/action/" {
				form, err := c.MultipartForm()
				if err != nil {
					c.AbortWithStatusJSON(400, utils.H{
						"status_code": 10333,
						"status_msg":  err.Error(),
					})
					hlog.Infof("finished %s err:%v", path, err)
					return
				}
				tokenString = form.Value["token"][0]
			} else {
				var douyinToken core.DouyinToken
				err := c.BindAndValidate(&douyinToken)
				if err != nil {
					c.AbortWithStatusJSON(400, utils.H{
						"status_code": 10333,
						"status_msg":  err.Error(),
					})
					hlog.Infof("finished %s err:%v", path, err)
					return
				}
				tokenString = douyinToken.Token
			}
			if tokenString == "" {
				err := errors.New("failed find token")
				c.AbortWithStatusJSON(400, utils.H{
					"status_code": 10333,
					"status_msg":  err.Error(),
				})
				hlog.Infof("finished %s err:%v", path, err)
				return
			}
			claims, err := utils2.ParseToken(tokenString)
			if err != nil {
				c.AbortWithStatusJSON(400, utils.H{
					"status_code": 10333,
					"status_msg":  err.Error(),
				})
				hlog.Infof("finished %s err:%v", path, err)
				return
			}
			c.Set("myId", claims.UserId)
		},
	}
}

func _douyinMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _comment_ctionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _commentlistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoriteMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favorite_ctionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _favoritelistMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _actionMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _listMw() []app.HandlerFunc {
	// your code...
	return nil
}

func _action0Mw() []app.HandlerFunc {
	// your code...
	return nil
}

func _list0Mw() []app.HandlerFunc {
	// your code...
	return nil
}
