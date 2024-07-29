<?php

namespace App\Providers;

use App\Repositories\ActionLogRepository;
use App\Repositories\IpAddressRepository;
use App\Services\ActionLogSearchService;
use App\Services\IpAddressManageService;
use App\Utils\ApiResponser;
use Illuminate\Support\ServiceProvider;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register any application services.
     */
    public function register(): void
    {
        $this->app->singleton(ApiResponser::class, ApiResponser::class);

        $this->app->singleton(IpAddressRepository::class, IpAddressRepository::class);

        $this->app->singleton(ActionLogRepository::class, ActionLogRepository::class);

        $this->app->singleton(IpAddressManageService::class, function ($app) {
            return new IpAddressManageService($app->make(IpAddressRepository::class));
        });

        $this->app->singleton(ActionLogSearchService::class, function ($app) {
            return new ActionLogSearchService($app->make(ActionLogRepository::class));
        });
    }

    /**
     * Bootstrap any application services.
     */
    public function boot(): void
    {
        //
    }
}
