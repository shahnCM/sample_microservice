<?php

namespace App\Repositories;

use App\Models\IpAddress;
use App\Interfaces\RepositoryInterface;
use App\Utils\ApiResponser;
use Illuminate\Http\Exceptions\HttpResponseException;

class IpAddressRepository implements RepositoryInterface
{
   public function getAll()
   {
      return IpAddress::paginate(15);
   }

   public function getById($id)
   {
      return IpAddress::findOrFail($id);
   }

   public function store(array $data)
   {
      return IpAddress::create($data);
   }

   public function update(array $data, $model)
   {
      return tap($model)->update($data);
   }

   public function search($userId = null, $before = null, $after = null, $start = null, $end = null)
   {
      throw new HttpResponseException(app(ApiResponser::class)->sendErrorResponse("Please search for action logs", 422));
   }
}