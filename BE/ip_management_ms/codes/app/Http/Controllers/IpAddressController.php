<?php

namespace App\Http\Controllers;

use App\Http\Requests\IpAddressRequest;
use App\Http\Resources\ActionLogCollection;
use App\Http\Resources\ActionLogResource;
use App\Http\Resources\IpAddressCollection;
use App\Services\ActionLogSearchService;
use Illuminate\Http\Request;
use App\Utils\ApiResponser;
use App\Http\Resources\IpAddressResource;
use App\Services\IpAddressManageService;
use Illuminate\Support\Facades\DB;
use Symfony\Component\HttpKernel\Exception\UnauthorizedHttpException;
use Symfony\Component\HttpKernel\Exception\UnprocessableEntityHttpException;

class IpAddressController extends Controller
{

    public function index()
    {
        return new IpAddressCollection(app(IpAddressManageService::class)->getAllIpAddresses());
    }

    public function show($id)
    {
        $ipAddress = app(IpAddressManageService::class)->getIpAddressById($id);
        return app(ApiResponser::class)->sendSuccessResponse(new IpAddressResource($ipAddress), '', 200);
    }

    public function store(IpAddressRequest $request)
    {
        DB::beginTransaction();
        try {
            $ipAddress = app(IpAddressManageService::class)->createIpAddress($request->validated());
            DB::commit();
            return app(ApiResponser::class)->sendSuccessResponse(new IpAddressResource($ipAddress), 'Ip Address Create Successful', 201);
        } catch (\Exception $e) {
            \Log::error("Critical Error Store: " . $e->__tostring());
            DB::rollback();
            throw new UnprocessableEntityHttpException("Save Ip Address Failed");
        }
    }

    public function update(IpAddressRequest $request, $id)
    {
        DB::beginTransaction();
        try {
            app(IpAddressManageService::class)->updateIpAddress($request->validated(), $id);
            DB::commit();
            return app(ApiResponser::class)->sendSuccessResponse(null, 'Ip Address Update Successful', 201);
        } catch (\Exception $e) {
            \Log::error("Critical Error Update: " . $e->__tostring());
            DB::rollback();
            throw new UnprocessableEntityHttpException("Update Ip Address Failed");
        }
    }

    public function destroy($id)
    {
        throw new UnauthorizedHttpException('Invalid JWT format');
    }

    public function actionLogs(Request $request)
    {
        return  new ActionLogCollection(app(ActionLogSearchService::class)->getActionLogs($request));
    }
}