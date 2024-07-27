<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Http\Exceptions\HttpResponseException;
use Illuminate\Validation\Validator;

class IpAddressRequest extends FormRequest
{
    public function authorize()
    {
        return true;
    }

    public function rules()
    {
        $rules = [
            'ip' => 'required|ip|unique:ip_addresses',
            'label' => 'nullable|string|max:255',
        ];

        return $rules;
    }

    public function failedValidation(Validator $validator)
    {
        throw new HttpResponseException(response()->json([
            'status' => "Error",
            'status_code' => 422,
            'message' => 'Validation Error',
            'errors' => $validator->errors()
        ], 422));
    }
}
