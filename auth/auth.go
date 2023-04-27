package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

// Auth 对象 也是需要通用
type Auth struct {
	UserID        string `json:"userId"`
	GroupID       string `json:"groupId"`
	Viewer        string `json:"viewer"`
	Admin         bool   `json:"admin"`
	From          string `json:"from"`
	SenderIsAdmin bool   `json:"senderIsAdmin"`
}

func NewInnerAuth() *Auth {
	return &Auth{
		From: Invoker,
	}
}

func (a *Auth) Empty() bool {
	if a == nil {
		return false
	}
	if a.InnerInvoke() {
		// 内部调用的不用管视图
		return true
	}
	if a.Viewer == "" {
		return false
	}
	return true
}

const (
	UserView  = "USER"
	GroupView = "ORGANIZATION"
	Invoker   = "blocface"
)

const HeaderKeyAuth = "x-forwarded-auth-user"

// GetOutGoingContext logger 中需要使用
func GetOutGoingContext(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		var m = make(map[string]string)
		m[HeaderKeyAuth] = md[HeaderKeyAuth][0]
		return metadata.NewOutgoingContext(ctx, metadata.New(m))
	} else {
		log.Errorf("cannot parse metadata from grpc ctx")
		return nil
	}
}

// CopyContext logger 中需要使用
func CopyContext(ctx context.Context) (context.Context, context.CancelFunc, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		cancel, cancelFunc := context.WithCancel(metadata.NewOutgoingContext(context.TODO(), md))
		return cancel, cancelFunc, nil
	} else {
		log.Errorf("cannot parse metadata from grpc ctx")
		return nil, nil, errors.New("context no metadata")
	}
}

// InnerAuthContext 内部服务调用的context生成
func (a *Auth) InnerAuthContext() (context.Context, context.CancelFunc) {
	marshal, _ := json.Marshal(a)
	md := metadata.New(map[string]string{HeaderKeyAuth: string(marshal)})
	return context.WithCancel(metadata.NewOutgoingContext(context.TODO(), md))
}

// GenChainsAuthSql logger 中需要使用
func (a *Auth) GenChainsAuthSql() string {
	var buf bytes.Buffer
	if a.Viewer == "USER" {
		buf.WriteString(fmt.Sprintf("invitation_rec.invitees_user_id = '%s", a.UserID))
	} else {
		buf.WriteString(fmt.Sprintf("invitation_rec.invitees_group_id = '%s", a.GroupID))
	}
	buf.WriteString("'")
	return buf.String()
}

func ParseFromGrpcCtxWithAdmin(ctx context.Context) (*Auth, error) {
	grpcCtx, err := ParseFromGrpcCtx(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "parse auth")
	}
	if grpcCtx.Admin {
		grpcCtx.Viewer = "ADMIN"
	}
	return grpcCtx, nil
}

// MustParseFromGrpcCtx  必须解析出 auth
func MustParseFromGrpcCtx(ctx context.Context) *Auth {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Debugf("start to parse auth form grpc header=%v with key=%s", md, HeaderKeyAuth)
		authJSONs := md[HeaderKeyAuth]
		if len(authJSONs) > 0 {
			authJSON := authJSONs[0]
			a := Auth{}
			err := json.Unmarshal([]byte(authJSON), &a)
			if err != nil {
				log.Errorf("fail to unmarshal auth from authJSON=%s, err=%v", authJSON, err)
				return nil
			}
			// 追加auth view 的逻辑
			if a.Empty() {
				return &a
			}
			return nil
		} else {
			//log.Errorf("authJSON is empty from grpc header=%v", md)
			return nil
		}
	} else {
		log.Errorf("cannot parse metadata from grpc ctx")
		return nil
	}
}

func ParseFromGrpcCtx(ctx context.Context) (*Auth, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		log.Debugf("start to parse auth form grpc header=%v with key=%s", md, HeaderKeyAuth)
		authJSONs := md[HeaderKeyAuth]
		if len(authJSONs) > 0 {
			authJSON := authJSONs[0]
			auth := Auth{}
			err := json.Unmarshal([]byte(authJSON), &auth)
			if err != nil {

				log.Errorf("fail to unmarshal auth from authJSON=%s, err=%v", authJSON, err)
				return nil, fmt.Errorf("fail to unmarshal auth from authJSON=%s, err=%v", authJSON, err)
			}
			//
			return &auth, nil
		} else {
			//log.Errorf("authJSON is empty from grpc header=%v", md)
			return nil, nil
		}
	} else {
		log.Errorf("cannot parse metadata from grpc ctx")
		return nil, nil
	}
}

func ToAuthParams(at *Auth) []string {
	if at != nil {
		return []string{at.UserID, at.GroupID, at.Viewer}
	} else {
		return nil
	}
}
func (a *Auth) GenAuthSql(table ...string) string {
	prefix := ""
	if table != nil {
		prefix = table[0] + "."
	}
	var buf bytes.Buffer
	if a.InnerInvoke() {
		// 是内部服务调用，目的是获取所有资源，不受视图限制
		return buf.String()
	}
	sql := a.genSearchSql(prefix)
	buf.WriteString(sql)
	if sql != "" {
		buf.WriteString(fmt.Sprintf(" and %sviewer = '", prefix))
	} else {
		buf.WriteString(fmt.Sprintf("%sviewer = '", prefix))
	}
	buf.WriteString(a.Viewer)
	buf.WriteString("'")
	return buf.String()
}
func (a *Auth) genSearchSql(prefix string) string {
	var buf bytes.Buffer
	if a.Admin {
		return buf.String()
	}
	switch a.Viewer {
	case "ORGANIZATION":
		if a.GroupID != "" {
			buf.WriteString(fmt.Sprintf("%sgroup_id = '", prefix))
			buf.WriteString(a.GroupID)
			buf.WriteString("'")

		}
		break
	case "USER":
		if a.UserID != "" {
			buf.WriteString(fmt.Sprintf("%suser_id = '", prefix))
			buf.WriteString(a.UserID)
			buf.WriteString("'")
		}
		break
	default:
		// TODO 是否需要添加肯定错误的sql 防止直接查询
	}
	return buf.String()
}

func (a *Auth) InUser() bool {
	return a.Viewer == UserView
}

func (a *Auth) InGroup() bool {
	return a.Viewer == GroupView
}

func (a *Auth) InnerInvoke() bool {
	return a.From == Invoker
}

func (a *Auth) GenGitOrgId() string {
	if a.Viewer == "USER" {
		return a.UserID
	} else {
		return a.GroupID
	}
}
