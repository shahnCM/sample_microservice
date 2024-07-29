<?php

use App\Http\Controllers\IpAddressController;
use Illuminate\Support\Facades\Route;

Route::group(['middleware' => 'ApiUserJwtAuth'], function () {
    Route::apiResource('ip_addresses', IpAddressController::class);
    Route::get('action-logs', [IpAddressController::class, 'actionLogs']);
});