<?php

namespace App\Http\Resources;

use Carbon\Carbon;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class ActionLogResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            'user_id' => $this->user_id,
            'username' => $this->username,
            'session_token_trace_id' => $this->session_token_trace_id,
            'action' => $this->action,
            'change' => $this->change,
            'logged_at' => $this->created_at->format('Y-m-d\TH:i:s.vP')
        ];
    }
}
