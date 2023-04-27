package auth

import (
	"context"
	"encoding/json"
	"fmt"
	log "github.com/Heqiaomu/glog"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"os"
	"testing"
)

func TestParseFromGrpcCtx(t *testing.T) {
	log.Logger()
	p1 := gomonkey.ApplyFunc(metadata.FromIncomingContext, func(ctx context.Context) (md metadata.MD, ok bool) {
		m := make(map[string][]string)
		auth := make([]string, 0)
		auth = append(auth, "{\"userId\":\"mock-user-id\"}")
		m["x-forwarded-auth-user"] = auth
		return m, true
	})
	ParseFromGrpcCtx(context.Background())
	p1.Reset()
	p2 := gomonkey.ApplyFunc(metadata.FromIncomingContext, func(ctx context.Context) (md metadata.MD, ok bool) {
		m := make(map[string][]string)
		auth := make([]string, 0)
		auth = append(auth, "{\"userId\":\"mock-user-id\"}")
		m["x-forwarded-auth-user"] = auth
		return m, false
	})
	ParseFromGrpcCtx(context.Background())
	p2.Reset()
	p3 := gomonkey.ApplyFunc(metadata.FromIncomingContext, func(ctx context.Context) (md metadata.MD, ok bool) {
		return make(map[string][]string), true
	})
	ParseFromGrpcCtx(context.Background())
	p3.Reset()
	p4 := gomonkey.ApplyFunc(metadata.FromIncomingContext, func(ctx context.Context) (md metadata.MD, ok bool) {
		m := make(map[string][]string)
		auth := make([]string, 0)
		auth = append(auth, "12345")
		m["x-forwarded-auth-user"] = auth
		return m, true
	})
	ParseFromGrpcCtx(context.Background())
	p4.Reset()
}

func TestToAuthParams(t *testing.T) {
	au := Auth{}
	ToAuthParams(&au)
	ToAuthParams(nil)
}

func TestAuth(t *testing.T) {
	marshal, _ := json.Marshal(Auth{
		UserID:  "USER-0b9525fe8b374e34abcd1b7a785cc75d",
		GroupID: "ORG-3d80d6bfacdd41c1998044515e4b2525",
		Viewer:  "ORGANIZATION",
	})
	fmt.Printf("%s", marshal)
}

func TestNewInnerAuth(t *testing.T) {
	t.Run("new inner auth", func(t *testing.T) {
		got := NewInnerAuth()
		assert.Equal(t, &Auth{
			From: Invoker,
		}, got)
	})

}

func TestAuth_Empty(t *testing.T) {
	type fields struct {
		UserID  string
		GroupID string
		Viewer  string
		Admin   bool
		From    string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "nil auth",
			want: false,
		},
		{
			name: "inner invoker",
			fields: fields{
				UserID:  "a",
				GroupID: "b",
				Viewer:  "ORGANIZATION",
				From:    Invoker,
			},
			want: true,
		},
		{
			name: "empty viewer",
			fields: fields{
				UserID:  "a",
				GroupID: "b",
				Viewer:  "",
				From:    "not invoker",
			},
			want: false,
		},
		{
			name: "empty viewer",
			fields: fields{
				UserID:  "a",
				GroupID: "b",
				Viewer:  "ORGANIZATION",
				From:    "not invoker",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				UserID:  tt.fields.UserID,
				GroupID: tt.fields.GroupID,
				Viewer:  tt.fields.Viewer,
				Admin:   tt.fields.Admin,
				From:    tt.fields.From,
			}
			assert.Equalf(t, tt.want, a.Empty(), "Empty()")
		})
	}
}

func TestGetOutGoingContext(t *testing.T) {
	log.Logger()
	defer os.Remove("log")
	headers := metadata.New(map[string]string{"x-forwarded-auth-user": "", "abc": "", "xyz": ""})
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "empty metadata",
			args: args{
				ctx: context.Background(),
			},
			want: nil,
		},
		{
			name: "context with metadata",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), headers),
			},
			want: GetOutGoingContext(metadata.NewIncomingContext(context.Background(), headers)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, GetOutGoingContext(tt.args.ctx), "GetOutGoingContext(%v)", tt.args.ctx)
		})
	}
}

func TestCopyContext(t *testing.T) {
	log.Logger()
	defer os.Remove("log")
	t.Run("empty metadata", func(t *testing.T) {
		got, got1, err := CopyContext(context.Background())
		assert.NotNil(t, err)
		assert.Nil(t, got)
		assert.Nil(t, got1)
	})

	t.Run("not empty metadata", func(t *testing.T) {
		headers := metadata.New(map[string]string{"x-forwarded-auth-user": "", "abc": "", "xyz": ""})
		got, got1, err := CopyContext(metadata.NewIncomingContext(context.Background(), headers))
		assert.NotNil(t, got)
		assert.NotNil(t, got1)
		assert.Nil(t, err)
	})

}

func TestAuth_GenChainsAuthSql(t *testing.T) {
	type fields struct {
		UserID  string
		GroupID string
		Viewer  string
		Admin   bool
		From    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "USER viewer",
			fields: fields{
				UserID: "a",
				Viewer: UserView,
			},
			want: "invitation_rec.invitees_user_id = 'a'",
		},
		{
			name: "ORGANIZATION viewer",
			fields: fields{
				GroupID: "org1",
				Viewer:  GroupView,
			},
			want: "invitation_rec.invitees_group_id = 'org1'",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Auth{
				UserID:  tt.fields.UserID,
				GroupID: tt.fields.GroupID,
				Viewer:  tt.fields.Viewer,
				Admin:   tt.fields.Admin,
				From:    tt.fields.From,
			}
			assert.Equalf(t, tt.want, a.GenChainsAuthSql(), "GenChainsAuthSql()")
		})
	}
}

func TestParseFromGrpcCtxWithAdmin(t *testing.T) {
	log.Logger()
	defer os.Remove("log")
	t.Run("parse from context", func(t *testing.T) {
		headers := metadata.New(map[string]string{"x-forwarded-auth-user": "{\"UserID\":\"a\",\"GroupID\":\"b\",\"Viewer\":\"c\",\"Admin\":false,\"From\":\"e\"}", "abc": "", "xyz": ""})
		got, err := ParseFromGrpcCtxWithAdmin(metadata.NewIncomingContext(context.Background(), headers))
		assert.NotNil(t, got)
		assert.Nil(t, err)
	})

}

func TestMustParseFromGrpcCtx(t *testing.T) {
	log.Logger()
	defer os.Remove("log")
	headers := metadata.New(map[string]string{"x-forwarded-auth-user": "{\"UserID\":\"a\",\"GroupID\":\"b\",\"Viewer\":\"c\",\"Admin\":false,\"From\":\"e\"}", "abc": "", "xyz": ""})

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *Auth
	}{
		{
			name: "right auth",
			args: args{
				ctx: metadata.NewIncomingContext(context.Background(), headers),
			},
			want: &Auth{
				UserID:  "a",
				GroupID: "b",
				Viewer:  "c",
				Admin:   false,
				From:    "e",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MustParseFromGrpcCtx(tt.args.ctx), "MustParseFromGrpcCtx(%v)", tt.args.ctx)
		})
	}
}
