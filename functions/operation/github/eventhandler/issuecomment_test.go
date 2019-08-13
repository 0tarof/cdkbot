package eventhandler

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sambaiz/cdkbot/functions/operation/github/client"
	"github.com/sambaiz/cdkbot/lib/config"
	"github.com/stretchr/testify/assert"
)

func TestEventHandlerIssueCommentCreated(t *testing.T) {
	tests := []struct {
		title                string
		in                   issueCommentEvent
		cfg                  config.Config
		baseBranch           string
		resultHasDiff        bool
		outState             client.State
		outStatusDescription string
		isError              bool
	}{
		{
			title: "no targets are matched",
			in: issueCommentEvent{
				ownerName:   "owner",
				repoName:    "repo",
				issueNumber: 1,
				commentBody: "/deploy TestStack",
				cloneURL:    "http://github.com/sambaiz/cdkbot",
			},
			cfg: config.Config{
				CDKRoot: ".",
				Targets: map[string]config.Target{
					"master": {},
				},
			},
			baseBranch:           "develop",
			outState:             client.StateSuccess,
			outStatusDescription: "No targets are matched",
		},
		{
			title: "comment diff and has diffs",
			in: issueCommentEvent{
				ownerName:   "owner",
				repoName:    "repo",
				issueNumber: 1,
				commentBody: "/diff",
				cloneURL:    "http://github.com/sambaiz/cdkbot",
			},
			cfg: config.Config{
				CDKRoot: ".",
				Targets: map[string]config.Target{
					"develop": {
						Contexts: map[string]string{
							"env": "stg",
						},
					},
				},
			},
			baseBranch:           "develop",
			resultHasDiff:        true,
			outState:             client.StateFailure,
			outStatusDescription: "Diffs still remain",
		},
		{
			title: "comment deploy and has no diffs",
			in: issueCommentEvent{
				ownerName:   "owner",
				repoName:    "repo",
				issueNumber: 1,
				commentBody: "/deploy TestStack",
				cloneURL:    "http://github.com/sambaiz/cdkbot",
			},
			cfg: config.Config{
				CDKRoot: ".",
				Targets: map[string]config.Target{
					"develop": {
						Contexts: map[string]string{
							"env": "stg",
						},
					},
				},
			},
			baseBranch:           "develop",
			resultHasDiff:        false,
			outState:             client.StateSuccess,
			outStatusDescription: "There are no diffs. Let's merge!",
		},
	}

	constructEventHandlerWithMock := func(
		ctx context.Context,
		ctrl *gomock.Controller,
		event issueCommentEvent,
		cfg config.Config,
		baseBranch string,
		resultHasDiff bool,
	) *EventHandler {
		githubClient, gitClient, configClient, cdkClient := constructSetupMocks(
			ctx,
			ctrl,
			event.ownerName,
			event.repoName,
			event.issueNumber,
			event.cloneURL,
			cfg,
			baseBranch,
		)

		if _, ok := cfg.Targets[baseBranch]; !ok {
			return &EventHandler{
				cli:    githubClient,
				git:    gitClient,
				config: configClient,
				cdk:    cdkClient,
			}
		}

		target := cfg.Targets[baseBranch]
		cdkPath := fmt.Sprintf("%s/%s", clonePath, cfg.CDKRoot)
		cmd := parseCommand(event.commentBody)
		if cmd.action == actionDiff {
			// doActionDiff()
			result := "result"
			cdkClient.EXPECT().Diff(cdkPath, cmd.args, target.Contexts).Return(result, resultHasDiff)
			githubClient.EXPECT().CreateComment(
				ctx,
				event.ownerName,
				event.repoName,
				event.issueNumber,
				fmt.Sprintf("### cdk diff %s\n```%s```", cmd.args, result),
			).Return(nil)
		} else if cmd.action == actionDeploy {
			// doActionDeploy()
			result := "result"
			cdkClient.EXPECT().Deploy(cdkPath, cmd.args, target.Contexts).Return(result, nil)
			cdkClient.EXPECT().Diff(cdkPath, "", target.Contexts).Return("", resultHasDiff)
			githubClient.EXPECT().CreateComment(
				ctx,
				event.ownerName,
				event.repoName,
				event.issueNumber,
				fmt.Sprintf("### cdk deploy %s\n```%s```\n%s", cmd.args, result, "All stacks have been deployed :tada:"),
			).Return(nil)
		}

		return &EventHandler{
			cli:    githubClient,
			git:    gitClient,
			config: configClient,
			cdk:    cdkClient,
		}
	}

	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			eventHandler := constructEventHandlerWithMock(ctx, ctrl, test.in, test.cfg, test.baseBranch, test.resultHasDiff)
			cmd := parseCommand(test.in.commentBody)
			state, statusDescription, err := eventHandler.issueCommentCreated(ctx, test.in, cmd)
			assert.Equal(t, test.isError, err != nil)
			assert.Equal(t, test.outState, state)
			assert.Equal(t, test.outStatusDescription, statusDescription)
		})
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		title   string
		in      string
		out     *command
		isError bool
	}{
		{
			title: "diff",
			in:    "/diff aaa bbb",
			out: &command{
				action: actionDiff,
				args:   "aaa bbb",
			},
		},
		{
			title: "deploy",
			in:    "/deploy aaa bbb",
			out: &command{
				action: actionDeploy,
				args:   "aaa bbb",
			},
		},
		{
			title: "unknown",
			in:    "/unknown aaa bbb",
			out:   nil,
		},
	}
	for _, test := range tests {
		t.Run(test.title, func(t *testing.T) {
			cmd := parseCommand(test.in)
			assert.Equal(t, test.out, cmd)
		})
	}
}
