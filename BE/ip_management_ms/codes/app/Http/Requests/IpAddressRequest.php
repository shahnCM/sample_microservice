<?php

namespace App\Http\Requests;

use App\Utils\ApiResponser;
use Illuminate\Contracts\Validation\Validator;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Http\Exceptions\HttpResponseException;

class IpAddressRequest extends FormRequest
{
    public function authorize()
    {
        return true;
    }

    public function rules()
    {
        $rules = [
            'label' => 'nullable|string|max:255',
        ];

        if ($this->isMethod('post')) {
            $rules['ip'] = 'required|ip|unique:ip_addresses,ip';
        } else if ($this->isMethod('put')) {
            $rules['ip'] = 'required|ip|unique:ip_addresses,ip,' . $this->route('ip_address');
        }



        return $rules;
    }

    public function failedValidation(Validator $validator)
    {
        throw new HttpResponseException(app(ApiResponser::class)->sendErrorResponse("Validation Error", 422, $validator->errors()->toArray()));
    }
}
