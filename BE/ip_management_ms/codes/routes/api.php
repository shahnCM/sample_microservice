<?php

use App\Http\Controllers\IpAddressController;
use Illuminate\Support\Facades\Route;

Route::group(['middleware' => 'ApiUserJwtAuth'], function () {
    Route::apiResource('v1/ip_addresses', IpAddressController::class);
    Route::get('v1/action-logs', [IpAddressController::class, 'actionLogs']);
});