<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class ActionLog extends Model
{
    use HasFactory;

    protected $fillable = ['user_id', 'username', 'session_token_trace_id', 'change'];

    protected $casts = [
        'change' => 'array',
    ];
}
