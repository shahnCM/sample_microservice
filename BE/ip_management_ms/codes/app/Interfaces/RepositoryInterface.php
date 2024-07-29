<?php

namespace App\Interfaces;
use Illuminate\Database\Eloquent\Model;

interface RepositoryInterface
{
    public function getAll();
    public function getById($id);
    public function search($userId, $before, $after, $start, $end);
    public function store(array $data);
    public function update(array $data, $model);
    // public function delete($id);
}
