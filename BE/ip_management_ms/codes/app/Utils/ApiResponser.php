<?php

namespace App\Utils;

use Illuminate\Support\Facades\DB;
use Illuminate\Http\Exceptions\HttpResponseException;
use Illuminate\Support\Facades\Log;

class ApiResponser
{
    public function sendSuccessResponse($result, $message, $code = 200)
    {
        $response = [
            'status' => "Success",
            'status_code' => $code,
            'message' => $message,
            'data' => $result
        ];
        if (!empty($message)) {
            $response['message'] = $message;
        }
        return response()->json($response, $code);
    }

    public function sendErrorResponse($message, $code = 200)
    {
        $response = [
            'status' => "Error",
            'status_code' => $code,
            'message' => "Ip Address Service Error",
        ];
        if (!empty($message)) {
            $response['message'] = $message;
        }
        return response()->json($response, $code);
    }

}