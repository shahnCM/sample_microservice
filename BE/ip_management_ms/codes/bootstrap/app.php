<?php

use App\Enums\StatusCode;
use App\Utils\ApiResponser;
use Illuminate\Foundation\Application;
use Symfony\Component\HttpFoundation\Response;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        api: __DIR__ . '/../routes/api.php',
        apiPrefix: 'ip_management/api',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware) {
        $middleware->alias([
            'ApiUserJwtAuth' => \App\Http\Middleware\ApiUserJwtAuth::class,
        ]);
    })
    ->withExceptions(function (Exceptions $exceptions) {
        return $exceptions->respond(function (Response $response) {
            return app(ApiResponser::class)->sendErrorResponse(StatusCode::fromCode($response->getStatusCode())->message(), $response->getStatusCode());
        });

    })->create();
