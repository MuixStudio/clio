/*
 * Copyright 2026 MuixStudio
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package driver is the dependency injection container.
//
// All components are wired here. Flow packages only see interfaces; the Driver
// is the only place that references concrete types.
//
// Configuration is loaded externally (via config.Load) and passed to New.
// The storage implementation is supplied by main via WithStore.
//
//	cfg, _ := config.Load("identity.yaml")
//	db, _ := gorm.Open(postgres.Open(cfg.DSN), &gorm.Config{})
//	store, _ := gormstore.New(db)
//	d := driver.New(cfg, driver.WithStore(store))
package driver

import (
	"context"
	"errors"
	"fmt"

	"github.com/muixstudio/clio/internal/courier"
	alertDomain "github.com/muixstudio/clio/internal/domain/alert"
	alertRepo "github.com/muixstudio/clio/internal/domain/alert/repository"
	channelRepo "github.com/muixstudio/clio/internal/domain/channel/repository"
	connectorDomain "github.com/muixstudio/clio/internal/domain/connector"
	connectorRepo "github.com/muixstudio/clio/internal/domain/connector/repository"
	teamRepo "github.com/muixstudio/clio/internal/domain/team/repository"
	"github.com/muixstudio/clio/internal/driver/config"
	loginflow "github.com/muixstudio/clio/internal/flow/login"
	logoutflow "github.com/muixstudio/clio/internal/flow/logout"
	recflow "github.com/muixstudio/clio/internal/flow/recovery"
	regflow "github.com/muixstudio/clio/internal/flow/registration"
	verflow "github.com/muixstudio/clio/internal/flow/verification"
	"github.com/muixstudio/clio/internal/identity"
	"github.com/muixstudio/clio/internal/infra/bus"
	"github.com/muixstudio/clio/internal/infra/zus/dispatch"
	"github.com/muixstudio/clio/internal/persistence"
	"github.com/muixstudio/clio/internal/session"
	codestrategy "github.com/muixstudio/clio/internal/strategy/code"
	oidcstrategy "github.com/muixstudio/clio/internal/strategy/oidc"
	"github.com/muixstudio/clio/internal/strategy/password"
	"go.uber.org/zap"
)

type Driver struct {
	cfg    *config.Config
	store  persistence.Store
	logger *zap.Logger

	pw   *password.Strategy
	oidc *oidcstrategy.Strategy // nil if oidc not enabled in config
	code *codestrategy.Strategy // nil if code not enabled in config

	courier courier.Courier

	bus         *bus.Bus
	dispatchMgr *dispatch.Manager

	loginHooks        *loginflow.HookExecutor
	registrationHooks *regflow.HookExecutor
	logoutHooks       *logoutflow.HookExecutor
	verificationHooks *verflow.HookExecutor
	recoveryHooks     *recflow.HookExecutor
}

// Option configures the Driver after config is applied.
// Only implementation-level choices belong here — scalar values come from cfg.
type Option func(*Driver)

// WithStore sets the Persistence implementation used by the driver.
func WithStore(store persistence.Store) Option {
	return func(d *Driver) { d.store = store }
}

func WithLogger(logger *zap.Logger) Option {
	return func(d *Driver) { d.logger = logger }
}

func wireCourier(cfg config.CourierConfig) courier.Courier {
	if cfg.Strategy == "smtp" {
		return courier.NewSMTPCourier(courier.SMTPConfig{
			Host:     cfg.SMTP.Host,
			Port:     cfg.SMTP.Port,
			Username: cfg.SMTP.Username,
			Password: cfg.SMTP.Password,
			From:     cfg.SMTP.From,
		})
	}
	return courier.NewConsoleCourier()
}

// New wires up the Driver from a Config.
// Strategies, courier, and lifespans are all derived from cfg.
func New(cfg *config.Config, opts ...Option) *Driver {
	if cfg == nil {
		panic(errors.New("config is nil"))
	}
	d := &Driver{cfg: cfg}

	for _, o := range opts {
		o(d)
	}
	if d.store == nil {
		panic("driver: persistence store is required")
	}

	if cfg.SelfService.Methods.Password.Enabled {
		d.pw = password.New(d)
	}

	if cfg.SelfService.Methods.OIDC.Enabled {
		d.oidc = oidcstrategy.New(d)
		for _, p := range cfg.SelfService.Methods.OIDC.Config.Providers {
			switch p.Provider {
			case "github":
				d.oidc.AddProvider(oidcstrategy.NewProviderGitHub(&oidcstrategy.GitHubConfig{
					ClientID:     p.ClientID,
					ClientSecret: p.ClientSecret,
				}))
			default:
				panic(fmt.Sprintf("driver: unknown OIDC provider %q — add a case to New()", p.Provider))
			}
		}
	}

	if cfg.SelfService.Methods.Code.Enabled {
		d.courier = wireCourier(cfg.Courier)
		d.code = codestrategy.New(d)
	}

	d.loginHooks = loginflow.NewHookExecutor(d)
	d.registrationHooks = regflow.NewHookExecutor(d)
	d.logoutHooks = logoutflow.NewHookExecutor(d)
	d.verificationHooks = verflow.NewHookExecutor(d)
	d.recoveryHooks = recflow.NewHookExecutor(d)

	// In-process alert bus decoupling ingestion (webhook) from dispatch.
	d.bus = bus.New(256)

	return d
}

// AlertPublisher returns the publisher webhook ingestion uses to hand alerts to
// the dispatch pipeline.
func (d *Driver) AlertPublisher() alertDomain.AlertPublisher {
	return dispatch.NewAlertPublisher(d.bus.Pub)
}

// StartAlertDispatch starts the dispatch manager's bus consumers. It must be
// called before the HTTP server begins accepting webhooks, because the
// in-process bus drops messages published with no subscriber. It returns once
// the manager has subscribed; the manager then runs until ctx is cancelled or
// StopAlertDispatch is called.
func (d *Driver) StartAlertDispatch(ctx context.Context) error {
	notifier := dispatch.NewLogNotifier(d.logger)
	mgr, err := dispatch.NewManager(d.bus.Sub, d.store, notifier, d.logger)
	if err != nil {
		return err
	}
	d.dispatchMgr = mgr
	return d.dispatchMgr.Start(ctx)
}

// StopAlertDispatch stops the dispatch manager and closes the bus.
func (d *Driver) StopAlertDispatch() {
	if d.dispatchMgr != nil {
		d.dispatchMgr.Stop()
	}
	if d.bus != nil {
		_ = d.bus.Close()
	}
}

// AlertRouteWriter returns the route CRUD side: the store the HTTP handler
// mutates synchronously before publishing a reload.
func (d *Driver) AlertRouteWriter() dispatch.RouteWriter { return d.store }

// AlertRouteReloader returns the reload side: a publisher that asks the dispatch
// Manager to rebuild a team's Dispatcher from the persisted tree.
func (d *Driver) AlertRouteReloader() dispatch.RouteReloader {
	return dispatch.NewReloadPublisher(d.bus.Pub)
}

// AlertRouteLoader returns the query side for alert route trees.
func (d *Driver) AlertRouteLoader() dispatch.RouteTreeLoader { return d.store }

func (d *Driver) OIDCStrategy() *oidcstrategy.Strategy { return d.oidc }

func (d *Driver) Courier() courier.Courier { return d.courier }

func (d *Driver) ConnectorProviders() map[string]connectorDomain.ConnectorProvider {
	return map[string]connectorDomain.ConnectorProvider{
		"prometheus": connectorDomain.NewPrometheusNormalizer(),
		//"zabbix":     &connector.ZabbixProvider{},
	}
}

func (d *Driver) Config() *config.Config {
	return d.cfg
}

// ---- Persistence — delegate everything to d.store -----------------------

func (d *Driver) TeamPersister() teamRepo.TeamPersister {
	return d.store
}

func (d *Driver) TeamMemberPersister() teamRepo.TeamMemberPersister {
	return d.store
}

func (d *Driver) LoggerProvider() *zap.Logger {
	return d.logger
}

func (d *Driver) ConnectorPersister() connectorRepo.ConnectorPersister {
	return d.store
}

func (d *Driver) AlertPersister() alertRepo.AlertPersister {
	return d.store
}

func (d *Driver) ChannelPersister() channelRepo.ChannelPersister {
	return d.store
}

func (d *Driver) ChannelMemberPersister() channelRepo.ChannelMemberPersister {
	return d.store
}

func (d *Driver) LoginFlowPersister() loginflow.FlowPersister {
	return d.store
}

func (d *Driver) LogoutFlowPersister() logoutflow.FlowPersister {
	return d.store
}

func (d *Driver) RegistrationFlowPersister() regflow.FlowPersister {
	return d.store
}

func (d *Driver) RecoveryFlowPersister() recflow.FlowPersister {
	return d.store
}

func (d *Driver) IdentityPersister() identity.IdentityPersister {
	return d.store
}

func (d *Driver) IdentityVerifiableAddressPersister() identity.IdentityVerifiableAddressPersister {
	return d.store
}

func (d *Driver) VerificationFlowPersister() verflow.FlowPersister {
	return d.store
}

func (d *Driver) VerificationCodePersister() codestrategy.VerificationCodePersister {
	return d.store
}

func (d *Driver) SessionPersister() session.SessionPersister {
	return d.store
}

func (d *Driver) SessionTokenExchangeCodePersister() session.SessionTokenExchangeCodePersister {
	return d.store
}

// ---- Hook executors -----------------------------------------------------

func (d *Driver) RegistrationHookExecutor() *regflow.HookExecutor { return d.registrationHooks }
func (d *Driver) LoginHookExecutor() *loginflow.HookExecutor      { return d.loginHooks }
func (d *Driver) LogoutHookExecutor() *logoutflow.HookExecutor    { return d.logoutHooks }
func (d *Driver) VerificationHookExecutor() *verflow.HookExecutor { return d.verificationHooks }
func (d *Driver) RecoveryHookExecutor() *recflow.HookExecutor     { return d.recoveryHooks }

// ---- Hooks (empty by default) -------------------------------------------

func (d *Driver) PreRegistrationHooks(_ context.Context) []regflow.PreHookExecutor   { return nil }
func (d *Driver) PostRegistrationHooks(_ context.Context) []regflow.PostHookExecutor { return nil }
func (d *Driver) PreLoginHooks(_ context.Context) []loginflow.PreHookExecutor        { return nil }
func (d *Driver) PostLoginHooks(_ context.Context) []loginflow.PostHookExecutor      { return nil }
func (d *Driver) PreLogoutHooks(_ context.Context) []logoutflow.PreHookExecutor      { return nil }
func (d *Driver) PostLogoutHooks(_ context.Context) []logoutflow.PostHookExecutor    { return nil }
func (d *Driver) PreVerificationHooks(_ context.Context) []verflow.PreHookExecutor   { return nil }
func (d *Driver) PostVerificationHooks(_ context.Context) []verflow.PostHookExecutor { return nil }
func (d *Driver) PreRecoveryHooks(_ context.Context) []recflow.PreHookExecutor       { return nil }
func (d *Driver) PostRecoveryHooks(_ context.Context) []recflow.PostHookExecutor     { return nil }

// ---- Strategies ---------------------------------------------------------

func (d *Driver) RegistrationStrategies() []regflow.Strategy {
	var strategies []regflow.Strategy
	// Password registration is gated by registration_enabled.
	if d.pw != nil && d.cfg.SelfService.Methods.Password.Config.RegistrationEnabled {
		strategies = append(strategies, d.pw)
	}
	if d.oidc != nil {
		strategies = append(strategies, d.oidc)
	}
	if d.code != nil {
		strategies = append(strategies, d.code)
	}
	return strategies
}

// LoginStrategies returns all enabled strategies as a map keyed by credential type
// for O(1) lookup in the login handler.
func (d *Driver) LoginStrategies() map[identity.CredentialsType]loginflow.Strategy {
	strategies := make(map[identity.CredentialsType]loginflow.Strategy)
	if d.pw != nil {
		strategies[d.pw.ID()] = d.pw
	}
	if d.oidc != nil {
		strategies[d.oidc.ID()] = d.oidc
	}
	if d.code != nil {
		strategies[d.code.ID()] = d.code
	}
	return strategies
}

// VerificationStrategies returns strategies eligible for the verification flow.
// Only the code strategy is included; OIDC and password don't verify addresses.
func (d *Driver) VerificationStrategies() []verflow.Strategy {
	var strategies []verflow.Strategy
	if d.code != nil {
		strategies = append(strategies, d.code)
	}
	return strategies
}

// RecoveryStrategies returns strategies eligible for the recovery flow.
// Only the code strategy is included — recovery via email code.
func (d *Driver) RecoveryStrategies() []recflow.Strategy {
	var strategies []recflow.Strategy
	if d.code != nil {
		strategies = append(strategies, d.code)
	}
	return strategies
}
