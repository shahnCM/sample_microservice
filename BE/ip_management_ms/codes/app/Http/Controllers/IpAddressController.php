<?php

namespace App\Http\Controllers;

use App\Http\Requests\IpAddressRequest;
use App\Utils\ApiResponser;
use App\Http\Resources\IpAddressResource;
use App\Services\IpAddressManageService;
use Illuminate\Support\Facades\DB;
class IpAddressController extends Controller
{

    public function index()
    {
        $data = app(IpAddressManageService::class)->getAllIpAddresses();
        return app(ApiResponser::class)->sendSuccessResponse(IpAddressResource::collection($data),'',200);
    }

    public function show($id)
    {
        $ipAddress = app(IpAddressManageService::class)->getIpAddressById($id);
        return app(ApiResponser::class)->sendSuccessResponse(new IpAddressResource($ipAddress),'',200);
    }

    public function store(IpAddressRequest $request)
    {
        DB::beginTransaction();
        try{
             $ipAddress = app(IpAddressManageService::class)->createIpAddress($request->toArray());
             DB::commit();
        }catch(\Exception $e){
            DB::rollback();
            throw new \Exception("Save Ip Address Failed", 422);
        }
        return app(ApiResponser::class)->sendSuccessResponse(new IpAddressResource($ipAddress),'Ip Address Create Successful',201);
    }

    public function update(IpAddressRequest $request, $id)
    {
        DB::beginTransaction();
        try{
             app(IpAddressManageService::class)->updateIpAddress($request->toArray(),$id);
             DB::commit();
        }catch(\Exception $e){
            DB::rollback();
            throw new \Exception("Update Ip Address Failed", 422);
        }
        return app(ApiResponser::class)->sendSuccessResponse('Ip Address Update Successful','',204);
    }
}