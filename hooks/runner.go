package hooks

import (
	"fmt"

	"github.com/marcinhlybin/docker-env/app"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type HookType string

const (
	PreStartHook  HookType = "pre-start"
	PostStartHook HookType = "post-start"
	PostStopHook  HookType = "post-stop"
)

func RunHooks(hookType HookType, paths []string, p *project.Project) error {
	logger.Info("Running %s hooks", hookType)
	for _, path := range paths {
		hook := NewHook(string(hookType), path, p.Name, p.ServiceName)
		if err := hook.Run(); err != nil {
			return fmt.Errorf("%s hook failed: %w", hookType, err)
		}
	}
	return nil
}

func RunPreStartHooks(p *project.Project, ctx *app.AppContext) error {
	if err := RunHooks(PreStartHook, ctx.Config.PreStartHooks, p); err != nil {
		return err
	}
	return nil
}

func RunPostStartHooks(p *project.Project, ctx *app.AppContext) error {
	if err := RunHooks(PostStartHook, ctx.Config.PostStartHooks, p); err != nil {
		return err
	}
	return nil
}

func RunPostStopHooks(p *project.Project, ctx *app.AppContext) error {
	if err := RunHooks(PostStopHook, ctx.Config.PostStopHooks, p); err != nil {
		return err
	}
	return nil
}
