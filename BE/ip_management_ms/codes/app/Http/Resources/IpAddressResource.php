<?php

namespace App\Http\Resources;

use Carbon\Carbon;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class IpAddressResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            'ip_address' => $this->ip,
            'label' => $this->label,
            'created_at' => $this->created_at->format('Y-m-d\TH:i:s.vP')
        ];
    }
}
