<?php

namespace App\Models;

use App\Repositories\LogRepository;
use App\Utils\ApiResponser;
use Carbon\Carbon;
use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Http\Exceptions\HttpResponseException;

class IpAddress extends Model
{
    use HasFactory;

    protected $changeLogs;

    protected $fillable = ['ip', 'label'];

    public function delete()
    {
        throw new HttpResponseException(app(ApiResponser::class)->sendErrorResponse("Deleting IP addresses is not allowed.", 401));
    }


    protected static function booted()
    {
        parent::boot();

        static::creating(function ($ipAddress) {
            $ipAddress->changeLogs = [
                'previous' => null,
                'present' => ['ip' => $ipAddress->ip, 'label' => $ipAddress->label]
            ];
        });

        static::updating(function ($ipAddress) {
            $ipAddress->changeLogs = [
                'previous' => ['ip' => $ipAddress->getOriginal('ip'), 'label' => $ipAddress->getOriginal('label')],
                'present' => ['ip' => $ipAddress->ip, 'label' => $ipAddress->label]
            ];
        });
    }

    // Public method to get change logs
    public function getChangeLogs()
    {
        return $this->changeLogs;
    }

}
