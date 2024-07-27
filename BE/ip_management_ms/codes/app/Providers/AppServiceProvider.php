<?php

namespace App\Providers;

use App\Interfaces\RepositoryInterface;
use App\Repositories\IpAddressRepository;
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

        $this->app->singleton(IpAddressManageService::class, function ($app) {
            return new IpAddressManageService($app->make(IpAddressRepository::class));
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
