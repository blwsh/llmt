<?php

namespace BenWatson\MyPhpProject\auth;

class Authorizer
{
    private $name;

    public function __construct(string $name)
    {
        $this->name = $name;
    }

    public function authorize(string $user): bool
    {
        return $user === $this->name;
    }
}