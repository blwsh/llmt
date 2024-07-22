<?php

include 'vendor/autoload.php';

use greeter\Gereeter;
use BenWatson\MyPhpProject\auth\Authorizer;

$authorizer = new Authorizer('ben');

if (!$authorizer->authorize('ben')) {
    die('You are authorized!');
}

$greeter = new Gereeter();

echo $greeter->greet();

