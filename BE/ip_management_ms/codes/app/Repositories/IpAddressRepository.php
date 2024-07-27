<?php

namespace App\Repositories;
use App\Models\IpAddress;
use App\Interfaces\RepositoryInterface;
class IpAddressRepository implements RepositoryInterface
{
    public function getAll(){
        return IpAddress::paginate(10);
    }

    public function getById($id){
       return IpAddress::findOrFail($id);
    }

    public function store(array $data){
       return IpAddress::create($data);
    }

    public function update(array $data,$id){
       return IpAddress::whereId($id)->update($data);
    }
}