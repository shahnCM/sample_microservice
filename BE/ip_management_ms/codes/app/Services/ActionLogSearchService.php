<?php

namespace App\Services;

use App\Interfaces\RepositoryInterface;
use App\Models\ActionLog;
use App\Models\IpAddress;
use App\Repositories\ActionLogRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class ActionLogSearchService
{
    protected $actionLogRepository;

    public function __construct(RepositoryInterface $repository)
    {
        $this->actionLogRepository = $repository;
    }

    public function getActionLogs(Request $request) {

        $userId = null;
        $before = null;
        $after = null;
        $start = null;
        $end = null;

        // Filter by user_id
        if ($request->has('user_id')) {
            $userId = $request->input('user_id');
        }

        // Filter by created_at before a certain date
        if ($request->has('before')) {
            $before = $request->input('before');
        }

        // Filter by created_at after a certain date
        if ($request->has('after')) {
            $after = $request->input('after');
        }

        // Filter by created_at between two dates
        if ($request->has(['start', 'end'])) {
            $start = $request->input('start');
            $end = $request->input('end');
        }

        return app(ActionLogRepository::class)->search($userId, $before, $after, $start, $end);
    }
}
