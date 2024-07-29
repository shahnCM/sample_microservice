<?php

namespace App\Http\Resources;

use Illuminate\Http\Resources\Json\ResourceCollection;

class ActionLogCollection extends ResourceCollection
{
    public function toArray($request)
    {
        return [
            'status' => 'Success',
            'status_code' => 200,
            'message' => 'Action Logs',
            'data' => $this->collection
        ];
    }
}