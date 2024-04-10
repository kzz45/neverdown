package clientgo

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kzz45/neverdown/pkg/websocket"
	"github.com/kzz45/neverdown/pkg/zaplogger"

	gorillawebsocket "github.com/gorilla/websocket"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	rbacv1 "github.com/kzz45/neverdown/pkg/apis/rbac/v1"
	"github.com/kzz45/neverdown/pkg/jingx/aggregator"
	"github.com/kzz45/neverdown/pkg/jingx/proto"
	"github.com/kzz45/neverdown/pkg/jingx/registry"
	"github.com/kzz45/neverdown/pkg/jingx/server"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Option struct {
	Address string
	// authentication
	Username string
	Password string
}

// Tag was the shortage struct of the guldan object
type Tag struct {
	Tag            string
	RepositoryMeta jingxv1.RepositoryMeta
	GitReference   jingxv1.GitReference
	DockerImage    jingxv1.DockerImage
}

type Client struct {
	Option *Option
	Handle func(in *proto.Response) error

	ctx             context.Context
	cancel          context.CancelFunc
	token           string
	expireAt        int64
	signalContext   context.Context
	signalCancel    context.CancelFunc
	websocketOption *websocket.Option

	cache *runtime
}

func New(option *Option) *Client {
	ctx, cancel := context.WithCancel(context.Background())
	c := &Client{
		Option: option,
		ctx:    ctx,
		cancel: cancel,
		cache:  newRuntime(),
	}
	if err := c.auth(); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	return c
}

func (c *Client) auth() error {
	meta := &rbacv1.AccountMeta{
		Username: c.Option.Username,
		Password: c.Option.Password,
	}
	raw, err := meta.Marshal()
	if err != nil {
		return err
	}
	req := &proto.Request{
		GroupVersionKind: rbacv1.GroupVersionKind{},
		Namespace:        "",
		Verb:             "",
		Raw:              raw,
	}
	data, err := req.Marshal()
	if err != nil {
		return err
	}
	u := url.URL{
		Scheme: "https",
		Host:   c.Option.Address,
		Path:   server.RouterAuthority,
	}
	resp, err := http.Post(u.String(), "multipart/form-data", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http StatusCode:%d", resp.StatusCode)
	}
	cot, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	res := &proto.Response{}
	if err = res.Unmarshal(cot); err != nil {
		return err
	}
	if res.Code != 0 {
		return fmt.Errorf("login error code:%d err:%s", res.Code, string(res.Raw))
	}
	content := &proto.Context{}
	if err = content.Unmarshal(res.Raw); err != nil {
		return err
	}
	c.token = content.Token
	c.expireAt = int64(content.ExpireAt)
	return nil
}

func (c *Client) keepAlive() {
	retryChannel := make(chan struct{})
	tick := time.NewTicker(time.Second * 2)
	defer tick.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-retryChannel:
			if c.signalCancel != nil {
				c.signalCancel()
			}
			if c.expireAt != 0 {
				c.expireAt = 0
				c.signalContext, c.signalCancel = nil, nil
			}
		case <-tick.C:
			if c.signalContext != nil {
				if c.expireAt-time.Now().Unix() >= 10 {
					continue
				}
			}
			if c.signalCancel != nil {
				c.signalCancel()
				<-retryChannel
			}
			if err := c.auth(); err != nil {
				zaplogger.Sugar().Error(err)
				continue
			}
			c.signalContext, c.signalCancel = context.WithCancel(c.ctx)
			c.retryConnect(retryChannel)
			go c.retryHandle()
		}
	}
}

func (c *Client) retryConnect(retryChannel chan<- struct{}) {
	req := &proto.Request{
		Verb: proto.VerbPing,
	}
	data, err := req.Marshal()
	if err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	u := url.URL{
		Scheme: "wss",
		Host:   c.Option.Address,
		Path:   strings.Replace(server.RouterApi, ":token", c.token, 1),
	}
	c.websocketOption = websocket.NewOption(c.signalContext, u, data, &gorillawebsocket.Dialer{
		NetDial:           nil,
		NetDialContext:    nil,
		Proxy:             http.ProxyFromEnvironment,
		TLSClientConfig:   nil,
		HandshakeTimeout:  5 * time.Second,
		ReadBufferSize:    0,
		WriteBufferSize:   0,
		WriteBufferPool:   nil,
		Subprotocols:      nil,
		EnableCompression: false,
		Jar:               nil,
	})
	c.websocketOption.RetryDuration = 1000
	c.websocketOption.MaxRetryCount = 1
	go websocket.NewClientWithReconnect(c.websocketOption, retryChannel)
}

func (c *Client) retryHandle() {
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-c.signalContext.Done():
			return
		case msg, isClose := <-c.websocketOption.Read():
			if !isClose {
				return
			}
			res := &proto.Response{}
			if err := res.Unmarshal(msg); err != nil {
				zaplogger.Sugar().Error(err)
				continue
			}
			if c.Handle == nil {
				if err := c.handle(res); err != nil {
					zaplogger.Sugar().Error(err)
				}
				continue
			}
			if err := c.Handle(res); err != nil {
				zaplogger.Sugar().Error(err)
			}
		}
	}
}

func (c *Client) handle(in *proto.Response) error {
	if in.Code != 0 {
		return fmt.Errorf("response error code:%d err:%s", in.Code, string(in.Raw))
	}
	var err error
	switch in.GroupVersionKind {
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Event"}:
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}:
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}:
		switch in.Verb {
		case proto.VerbList:
			obj := &jingxv1.RepositoryList{}
			if err = obj.Unmarshal(in.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			zaplogger.Sugar().Debugf("list repository:%#v", obj)
			for _, v := range obj.Items {
				c.cache.repository(proto.EventAdded, v.DeepCopy())
			}
		case proto.VerbWatch:
			event := &proto.WatchEvent{}
			if err = event.Unmarshal(in.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			obj := &jingxv1.Repository{}
			if err = obj.Unmarshal(event.Raw); err != nil {
				zaplogger.Sugar().Error(err)
				return err
			}
			c.cache.repository(proto.EventType(event.Type), obj)
		case proto.VerbCreate:
		case proto.VerbDelete:
		case proto.VerbUpdate:
		default:
			err = fmt.Errorf("invaild repository handle verb:%s", in.Verb)
		}
	case rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}:
		switch in.Verb {
		case proto.VerbList:
			obj := &jingxv1.TagList{}
			if err = obj.Unmarshal(in.Raw); err != nil {
				return err
			}
			zaplogger.Sugar().Debugf("list tag:%#v", obj)
			for _, v := range obj.Items {
				c.cache.tag(proto.EventAdded, v.DeepCopy())
			}
		case proto.VerbWatch:
			event := &proto.WatchEvent{}
			if err = event.Unmarshal(in.Raw); err != nil {
				return err
			}
			obj := &jingxv1.Tag{}
			if err = obj.Unmarshal(event.Raw); err != nil {
				return err
			}
			c.cache.tag(proto.EventType(event.Type), obj)
		case proto.VerbCreate:
		case proto.VerbDelete:
		case proto.VerbUpdate:
		default:
			err = fmt.Errorf("invaild tag handle verb:%s", in.Verb)
		}
	default:
		err = fmt.Errorf(aggregator.ErrGVKNotExist, in.GroupVersionKind)
	}
	return nil
}

func (c *Client) DryRun() {
	go c.keepAlive()
}

func (c *Client) Lister() {
	<-time.After(time.Second * 30)
	zaplogger.Sugar().Infof("listerTimer")
	if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}, proto.VerbList, []byte("")); err != nil {
		zaplogger.Sugar().Error(err)
	}
	<-time.After(time.Second * 1)
	if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}, proto.VerbList, []byte("")); err != nil {
		zaplogger.Sugar().Error(err)
	}
	<-time.After(time.Second * 1)
	if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}, proto.VerbList, []byte("")); err != nil {
		zaplogger.Sugar().Error(err)
	}
	//go c.listerTimer()
}

func (c *Client) listerTimer() {
	tick := time.NewTicker(time.Second * 60)
	defer tick.Stop()
	for {
		select {
		case <-c.ctx.Done():
			return
		case <-tick.C:
			zaplogger.Sugar().Infof("listerTimer")
			if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}, proto.VerbList, []byte("")); err != nil {
				zaplogger.Sugar().Error(err)
			}
			<-time.After(time.Second * 1)
			if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}, proto.VerbList, []byte("")); err != nil {
				zaplogger.Sugar().Error(err)

			}
			<-time.After(time.Second * 1)
			if err := c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}, proto.VerbList, []byte("")); err != nil {
				zaplogger.Sugar().Error(err)
			}
		}
	}
}

func (c *Client) Shutdown() {
	c.cancel()
}

func (c *Client) send(gvk rbacv1.GroupVersionKind, verb proto.Verb, raw []byte) error {
	req := &proto.Request{
		GroupVersionKind: gvk,
		Namespace:        registry.DefaultNamespace,
		Verb:             verb,
		Raw:              raw,
	}
	data, err := req.Marshal()
	if err != nil {
		return err
	}
	if c.websocketOption == nil {
		return fmt.Errorf("nil websocketOption")
	}
	return c.websocketOption.Send(data)
}

func (c *Client) ListProjects() error {
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}, proto.VerbList, []byte(""))
}

func (c *Client) ListRepositories() error {
	if c.cache.repositories.Count() > 0 {
		return nil
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}, proto.VerbList, []byte(""))
}

func (c *Client) ListTags() error {
	if c.cache.tags.Count() > 0 {
		return nil
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}, proto.VerbList, []byte(""))
}

func (c *Client) UploadTag(tag *Tag) error {
	if err := c.ListTags(); err != nil {
		return err
	}
	if err := c.ListRepositories(); err != nil {
		return err
	}
	<-time.After(time.Millisecond * 1000)
	// check repository
	if bo := c.cache.checkRepositoryExist(
		&jingxv1.Repository{
			ObjectMeta: metav1.ObjectMeta{
				Name: registry.GenRepositoryFullName(tag.RepositoryMeta),
			},
		},
	); !bo {
		if err := c.CreateRepository(tag.RepositoryMeta); err != nil {
			return err
		}
		<-time.After(time.Millisecond * 2000)
	}
	obj := &jingxv1.Tag{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "",
			Namespace: registry.DefaultNamespace,
		},
		Spec: jingxv1.TagSpec{
			RepositoryMeta: tag.RepositoryMeta,
			Tag:            tag.Tag,
			GitReference:   tag.GitReference,
			DockerImage:    tag.DockerImage,
		},
	}
	// todo maybe it could be improved here?
	<-time.After(time.Millisecond * 1000)
	verb := proto.VerbCreate
	if ori, bo := c.cache.checkTagExist(obj); bo {
		verb = proto.VerbUpdate
		obj.Name = ori.Name
	}
	raw, err := obj.Marshal()
	if err != nil {
		return err
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Tag"}, verb, raw)
}

func (c *Client) CreateProject(projectName string) error {
	p := &jingxv1.Project{
		ObjectMeta: metav1.ObjectMeta{
			Name:      projectName,
			Namespace: registry.DefaultNamespace,
		},
		Spec: jingxv1.ProjectSpec{
			GenerateId: "",
			Domains:    []string{"127.0.0.1", "192.168.10.10"},
		},
	}
	raw, err := p.Marshal()
	if err != nil {
		return err
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Project"}, proto.VerbCreate, raw)
}

func (c *Client) CreateRepository(meta jingxv1.RepositoryMeta) error {
	rep := &jingxv1.Repository{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: registry.DefaultNamespace,
		},
		Spec: jingxv1.RepositorySpec{
			RepositoryMeta: jingxv1.RepositoryMeta{
				ProjectName:    meta.ProjectName,
				RepositoryName: meta.RepositoryName,
			},
		},
	}
	raw, err := rep.Marshal()
	if err != nil {
		return err
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}, proto.VerbCreate, raw)
}

func (c *Client) DeleteRepository(meta jingxv1.RepositoryMeta) error {
	rep := &jingxv1.Repository{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      registry.GenRepositoryFullName(meta),
			Namespace: registry.DefaultNamespace,
		},
		Spec: jingxv1.RepositorySpec{
			RepositoryMeta: jingxv1.RepositoryMeta{
				ProjectName:    meta.ProjectName,
				RepositoryName: meta.RepositoryName,
			},
		},
	}
	raw, err := rep.Marshal()
	if err != nil {
		return err
	}
	return c.send(rbacv1.GroupVersionKind{Group: jingxv1.GroupName, Version: jingxv1.SchemeGroupVersion.Version, Kind: "Repository"}, proto.VerbDelete, raw)
}

func (c *Client) init() {
	if err := c.CreateProject("lunara-common"); err != nil {
		zaplogger.Sugar().Error(err)
	}
	if err := c.CreateProject("helix2"); err != nil {
		zaplogger.Sugar().Error(err)
	}
	if err := c.CreateProject("hamster"); err != nil {
		zaplogger.Sugar().Error(err)
	}
	<-time.After(time.Second * 2)
}
