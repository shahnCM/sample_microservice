<?php

namespace App\Http\Resources;

use Illuminate\Http\Resources\Json\ResourceCollection;

class IpAddressCollection extends ResourceCollection
{
    public function toArray($request)
    {
        return [
            'status' => 'Success',
            'status_code' => 200,
            'message' => 'Ip Address',
            'data' => $this->collection
        ];
    }
}