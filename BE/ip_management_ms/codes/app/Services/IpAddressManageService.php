<?php

namespace App\Services;

use App\Interfaces\RepositoryInterface;
use App\Models\IpAddress;
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
        $this->logAction('store', $ipAddress, Auth::user());
        return $ipAddress;
    }

    public function updateIpAddress(array $data, int $id): IpAddress
    {
        $ipAddress = $this->ipAddressRepository->update($data, $id);
        $this->logAction('update', $ipAddress, Auth::user());
        return $ipAddress;
    }

    protected function logAction(string $action, IpAddress $ipAddress, $user): void
    {
        $ipAddress->logs()->create([
            'user_id' => $user->id,
            'username' => $user->name,
            'action' => $action,
            'change' => json_encode($ipAddress->getChanges()),
            'logged_at' => now(),
        ]);
    }
}
