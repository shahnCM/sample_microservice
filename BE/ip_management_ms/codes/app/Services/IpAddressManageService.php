<?php

namespace App\Services;

use App\Interfaces\RepositoryInterface;
use App\Models\ActionLog;
use App\Models\IpAddress;
use App\Repositories\ActionLogRepository;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class IpAddressManageService
{
    protected $ipAddressRepository;

    public function __construct(RepositoryInterface $repository)
    {
        $this->ipAddressRepository = $repository;
    }

    public function getAllIpAddresses()
    {
        return $this->ipAddressRepository->getAll();
    }

    public function getIpAddressById($id)
    {
        return $this->ipAddressRepository->getById($id);
    }

    public function createIpAddress(array $data): IpAddress
    {
        $ipAddress = $this->ipAddressRepository->store($data);
        $this->logAction('store', $ipAddress->getChangeLogs(), Auth::user());

        return $ipAddress;
    }

    public function updateIpAddress(array $data, int $id): IpAddress
    {
        $ipAddress = $this->ipAddressRepository->getById($id);
        $this->ipAddressRepository->update($data, $ipAddress);
        $this->logAction('update', $ipAddress->getChangeLogs(), Auth::user());

        return $ipAddress;
    }

    protected function logAction(string $action, string|array $change, $user): void
    {
        ActionLog::create([
            'user_id' => $user->id,
            'session_token_trace_id' => $user->sessionTokenTraceId,
            'action' => $action,
            'change' => $change,
        ]);
    }
}
