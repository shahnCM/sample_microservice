<?php

namespace App\Repositories;

use App\Models\ActionLog;
use App\Interfaces\RepositoryInterface;

class ActionLogRepository implements RepositoryInterface
{
    public function getAll()
    {
        return ActionLog::paginate(10);
    }

    public function getById($id)
    {
        return ActionLog::findOrFail($id);
    }

    public function store(array $data)
    {
        return ActionLog::create($data);
    }

    public function update(array $data, $model)
    {
        return tap($model)->update($data);
    }

    public function search($userId = null, $before = null, $after = null, $start = null, $end = null)
    {
        $query = ActionLog::query();

        // Filter by user_id
        if (!is_null($userId)) {
            $query->where('user_id', $userId);
        }

        // Filter by created_at before a certain date
        if (!is_null($before)) {
            $query->where('created_at', '<', $before);
        }

        // Filter by created_at after a certain date
        if (!is_null($after)) {
            $query->where('created_at', '>', $after);
        }

        // Filter by created_at between two dates
        if (!is_null($start) && !is_null($end)) {
            $query->whereBetween('created_at', [$start, $end]);
        }

        // Get the results
        return $query->paginate(15);

    }
}