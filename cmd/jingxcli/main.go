package main

import (
	"flag"
	"time"

	"github.com/kzz45/neverdown/pkg/zaplogger"

	jingxv1 "github.com/kzz45/neverdown/pkg/apis/jingx/v1"
	client "github.com/kzz45/neverdown/pkg/jingx/client-go"
)

type args struct {
	address    *string
	username   *string
	password   *string
	project    *string
	repository *string
	tag        *string
	git        *string
	branch     *string
	commitHash *string
	sha256     *string
}

func (a args) validate() {
	if *a.address == "" {
		zaplogger.Sugar().Fatalf("empty address")
	}
	if *a.username == "" {
		zaplogger.Sugar().Fatalf("empty username")
	}
	if *a.password == "" {
		zaplogger.Sugar().Fatalf("empty password")
	}
	if *a.project == "" {
		zaplogger.Sugar().Fatalf("empty project")
	}
	if *a.repository == "" {
		zaplogger.Sugar().Fatalf("empty repository")
	}
	if *a.tag == "" {
		zaplogger.Sugar().Fatalf("empty tag")
	}
	if *a.git == "" {
		zaplogger.Sugar().Fatalf("empty git")
	}
	if *a.branch == "" {
		zaplogger.Sugar().Fatalf("empty branch")
	}
	if *a.commitHash == "" {
		zaplogger.Sugar().Fatalf("empty commitHash")
	}
	if *a.sha256 == "" {
		zaplogger.Sugar().Fatalf("empty sha256")
	}
}

func main() {
	args := args{
		address:    flag.String("address", "", "address"),
		username:   flag.String("username", "", "username"),
		password:   flag.String("password", "", "password"),
		project:    flag.String("project", "neverdown", "project"),
		repository: flag.String("repository", "openx", "repository"),
		tag:        flag.String("tag", "", "tag"),
		git:        flag.String("git", "", "git"),
		branch:     flag.String("branch", "", "branch"),
		commitHash: flag.String("commitHash", "", "commitHash"),
		sha256:     flag.String("sha256", "", "sha256"),
	}
	flag.Parse()

	args.validate()
	opt := &client.Option{
		Address:  *args.address,
		Username: *args.username,
		Password: *args.password,
	}
	c := client.New(opt)
	c.DryRun()
	zaplogger.Sugar().Info("jingx-cli is preparing in 5s")
	<-time.After(time.Second * 5)
	zaplogger.Sugar().Info("jingx-cli is ready")

	ct := &client.Tag{
		Tag: *args.tag,
		RepositoryMeta: jingxv1.RepositoryMeta{
			ProjectName:    *args.project,
			RepositoryName: *args.repository,
		},
		GitReference: jingxv1.GitReference{
			Git:        *args.git,
			Branch:     *args.branch,
			CommitHash: *args.commitHash,
		},
		DockerImage: jingxv1.DockerImage{
			Sha256:           *args.sha256,
			Author:           *args.username,
			LastModifiedTime: time.Now().Unix(),
		},
	}
	zaplogger.Sugar().Infof("Tag info:%#v", ct)
	if err := c.UploadTag(ct); err != nil {
		zaplogger.Sugar().Fatal(err)
	}
	<-time.After(time.Second * 2)
	c.Shutdown()
	zaplogger.Sugar().Info("jingx-cli is preparing is stopping")
	<-time.After(time.Second * 2)
	zaplogger.Sugar().Info("jingx-cli is preparing is shutdown")
}
