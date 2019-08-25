package command

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/sambaiz/cdkbot/functions/operation/constant"

	"github.com/golang/mock/gomock"
	cdkMock "github.com/sambaiz/cdkbot/functions/operation/cdk/mock"
	"github.com/sambaiz/cdkbot/functions/operation/config"
	configMock "github.com/sambaiz/cdkbot/functions/operation/config/mock"
	gitMock "github.com/sambaiz/cdkbot/functions/operation/git/mock"
	platformMock "github.com/sambaiz/cdkbot/functions/operation/platform/mock"
	"github.com/stretchr/testify/assert"
)

func TestRunner_updateStatus(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	platformClient := platformMock.NewMockClienter(ctrl)
	resultState := constant.StateMergeReady
	statusDescription := "description"
	platformClient.EXPECT().SetStatus(ctx, constant.StateRunning, "").Return(nil)
	platformClient.EXPECT().AddLabel(ctx, constant.LabelRunning).Return(nil)
	platformClient.EXPECT().SetStatus(ctx, resultState, statusDescription).Return(nil)
	platformClient.EXPECT().RemoveLabel(ctx, constant.LabelRunning).Return(nil)
	runner := Runner{
		platform: platformClient,
	}
	assert.Nil(t, runner.updateStatus(
		ctx,
		func() (constant.State, string, error) {
			return resultState, statusDescription, nil
		},
	))
}

func TestRunner_setup(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	platformClient := platformMock.NewMockClienter(ctrl)
	gitClient := gitMock.NewMockClienter(ctrl)
	configClient := configMock.NewMockReaderer(ctrl)
	cdkClient := cdkMock.NewMockClienter(ctrl)

	baseBranch := "develop"
	cfg := config.Config{
		CDKRoot: ".",
		Targets: map[string]config.Target{
			baseBranch: {},
		},
	}
	constructSetupMock(
		ctx,
		platformClient,
		gitClient,
		configClient,
		cdkClient,
		cfg,
		baseBranch,
	)
	runner := &Runner{
		platform: platformClient,
		git:      gitClient,
		config:   configClient,
		cdk:      cdkClient,
	}
	cdkPath, retCfg, retTarget, err := runner.setup(ctx)
	assert.Equal(t, fmt.Sprintf("%s/%s", clonePath, cfg.CDKRoot), cdkPath)
	assert.Equal(t, *retCfg, cfg)
	assert.Equal(t, *retTarget, cfg.Targets[baseBranch])
	assert.Nil(t, err)
}

func constructSetupMock(
	ctx context.Context,
	platformClient *platformMock.MockClienter,
	gitClient *gitMock.MockClienter,
	configClient *configMock.MockReaderer,
	cdkClient *cdkMock.MockClienter,
	cfg config.Config,
	baseBranch string,
) {
	hash := "hash"
	platformClient.EXPECT().GetPullRequestLatestCommitHash(ctx).Return(hash, nil)
	platformClient.EXPECT().GetPullRequestBaseBranch(ctx).Return(baseBranch, nil)
	gitClient.EXPECT().Clone(clonePath, &hash).Return(nil)
	gitClient.EXPECT().Merge(clonePath, fmt.Sprintf("remotes/origin/%s", baseBranch)).Return(nil)
	configClient.EXPECT().Read(fmt.Sprintf("%s/cdkbot.yml", clonePath)).Return(&cfg, nil)
	_, ok := cfg.Targets[baseBranch]
	if !ok {
		return
	}

	cdkPath := fmt.Sprintf("%s/%s", clonePath, cfg.CDKRoot)
	cdkClient.EXPECT().Setup(cdkPath).Return(nil)

	return
}

func TestValidateStackName(t *testing.T) {
	tests := []struct {
		title              string
		in                 string
		isError            bool
	}{
		{
			title: "valid",
			in: "Stack-1",
			isError: false,
		},
		{
			title: "invalid character",
			in: "Sta`ck1",
			isError: true,
		},
		{
			title: "too long",
			in: strings.Repeat("A", 129),
			isError: true,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.isError, validateStackName(test.in) != nil)
	}
}

