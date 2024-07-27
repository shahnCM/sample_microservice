<?php

namespace App\Models;

use App\Repositories\LogRepository;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class IpAddress extends Model
{
    use HasFactory;

    protected $fillable = ['ip', 'label'];

    public function delete()
    {
        throw new \Exception("Deleting IP addresses is not allowed.");
    }

    /*
    protected static function booted()
    {
        parent::boot();

        $logRepo = new LogRepository;

        static::created(function ($ipAddress) use ($logRepo) {
            $userId = auth()->id();
            $username = auth()->user()->username;
            $change = json_encode([
                'previous' => [],
                'current' => $ipAddress->getAttributes()
            ]);
            $logRepo->createLog($userId, $username, 'store', $change);
        });

        static::updated(function ($ipAddress) use ($logRepo) {
            $userId = auth()->id();
            $username = auth()->user()->username;
            $change = json_encode([
                'previous' => $ipAddress->getOriginal(),
                'current' => $ipAddress->getAttributes(),
            ]);
            $logRepo->createLog($userId, $username, 'update', $change);
        });

        static::deleted(function ($ipAddress) use ($logRepo) {
            $userId = auth()->id();
            $username = auth()->user()->username;
            $change = json_encode([
                'previous' => $ipAddress->getOriginal(),
                'current' => []
            ]);
            $logRepo->createLog($userId, $username, 'delete', $change);
        });
    }
    */
}
