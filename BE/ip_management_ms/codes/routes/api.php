<?php

use App\Http\Controllers\IpAddressController;
use Illuminate\Support\Facades\Route;

Route::group(['middleware' => 'ApiUserJwtAuth'], function () {
    Route::apiResource('ip_addresses', IpAddressController::class)->except(['destroy']);
});