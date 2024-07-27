<?php

namespace App\Repositories;

use App\Models\Log;

class LogRepository
{
    public function createLog($userId, $username, $action, $change)
    {
        return Log::create([
            'user_id' => $userId,
            'username' => $username,
            'action' => $action,
            'change' => $change,
            'logged_at' => now(),
        ]);
    }
}
