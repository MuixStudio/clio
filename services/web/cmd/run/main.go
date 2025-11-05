package run

import (
	"context"
	"fmt"

	kzap "github.com/go-kratos/kratos/contrib/log/zap/v2"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/muixstudio/clio/services/common/metadata"
	"github.com/muixstudio/clio/services/web/config"
	"github.com/muixstudio/clio/services/web/router"
	"github.com/muixstudio/clio/services/web/svc"
	"github.com/muixstudio/clio/services/web/svc/observability"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	flags = &cmdFlags{}
)

func init() {
	Cmd.PersistentFlags().StringVarP(
		&flags.configPath,
		"config",
		"c",
		"etc/config.yaml",
		"path to configuration file",
	)
}

// cmdFlags 命令行参数
type cmdFlags struct {
	configPath string
}

// App 应用封装
type App struct {
	cfg    *config.Config
	svcCtx *svc.ServiceContext
	logger log.Logger
	flags  *cmdFlags
}

var Cmd = &cobra.Command{

	Use:   "run",
	Short: "Start the web server",
	Long:  `Start the web server with the specified configuration file.`,
	Example: `web run
				  web run --config /path/to/config.yaml
				  web run -c ./custom-config.yaml
				`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApp(cmd.Context(), flags)
	},
}

// runApp 运行应用主逻辑
func runApp(ctx context.Context, flags *cmdFlags) error {
	// 1. 初始化应用
	app, err := initializeApp(ctx, flags)
	if err != nil {
		return fmt.Errorf("failed to initialize app: %w", err)
	}

	// 2. 创建 Kratos 应用
	kratosApp, err := app.createKratosApp()
	if err != nil {
		return fmt.Errorf("failed to create kratos app: %w", err)
	}

	// 3. 运行应用
	if err := kratosApp.Run(); err != nil {
		return fmt.Errorf("app run failed: %w", err)
	}

	return nil
}

// initializeApp 初始化应用
func initializeApp(ctx context.Context, flags *cmdFlags) (*App, error) {
	app := &App{flags: flags}

	// 加载配置（带默认值）
	cfg, err := config.Parse(flags.configPath)
	if err != nil {
		return nil, fmt.Errorf("load config failed: %w", err)
	}
	app.cfg = &cfg

	// 初始化日志
	logger, err := initLogger(&cfg)
	if err != nil {
		return nil, fmt.Errorf("init logger failed: %w", err)
	}
	app.logger = logger

	// 4. 初始化可观测性
	observability.InitMeterProvider()

	// 5. 初始化服务上下文
	svcCtx, err := svc.NewServiceContext(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("init service context failed: %w", err)
	}
	app.svcCtx = svcCtx

	return app, nil
}

// initLogger 初始化日志系统
func initLogger(cfg *config.Config) (log.Logger, error) {
	// 根据环境决定日志模式
	var zapLogger *zap.Logger
	var err error

	switch cfg.Env {
	case "production":
		zapLogger, err = zap.NewProduction()
	case "staging":
		// Staging 使用生产配置但稍微宽松
		zapConfig := zap.NewProductionConfig()
		zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		zapLogger, err = zapConfig.Build()
	default: // development
		zapLogger = zap.NewExample()
	}

	if err != nil {
		return nil, fmt.Errorf("create zap logger: %w", err)
	}

	logger := kzap.NewLogger(zapLogger)

	// 设置为全局默认 logger
	log.SetLogger(logger)

	return logger, nil
}

// createKratosApp 创建 Kratos 应用
func (app *App) createKratosApp() (*kratos.App, error) {
	// 创建 HTTP 服务器
	httpSrv, err := app.createHTTPServer()
	if err != nil {
		return nil, fmt.Errorf("create http server: %w", err)
	}

	meta := metadata.Get()

	// 创建 Kratos 应用（使用编译时注入的元数据）
	kratosApp := kratos.New(
		kratos.Name(meta.Name),
		kratos.Version(meta.Version),
		kratos.Metadata(map[string]string{
			"description": meta.Description,
			"env":         app.cfg.Env,
			"build_time":  meta.BuildTime,
			"git_commit":  meta.GitCommit,
			"git_branch":  meta.GitBranch,
		}),
		kratos.Logger(app.logger),
		kratos.Server(httpSrv),
		kratos.Context(context.Background()),
	)

	return kratosApp, nil
}

// createHTTPServer 创建 HTTP 服务器
func (app *App) createHTTPServer() (*http.Server, error) {
	// 初始化路由
	engine := router.NewEngine(context.Background(), app.svcCtx)

	// 构建监听地址
	addr := fmt.Sprintf("%s:%d", app.cfg.Host, app.cfg.Port)

	// 创建 HTTP 服务器
	httpSrv := http.NewServer(
		http.Address(addr),
		// 可以从配置中读取更多选项
		// http.Timeout(app.cfg.Server.HTTP.Timeout),
	)

	// 注册路由
	httpSrv.HandlePrefix("/", engine)

	return httpSrv, nil
}

func tip() {
	fmt.Println(`
       .__  .__        
  ____ |  | |__| ____  
_/ ___\|  | |  |/  _ \ 
\  \___|  |_|  (  <_> )
 \___  >____/__|\____/ 
     \/                
	`)
}
