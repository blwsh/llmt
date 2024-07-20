# Authorizer.php Documentation

## Overview
The `Authorizer` class is responsible for user authorization within the `BenWatson\MyPhpProject\auth` namespace. It checks if a given user matches the initialized name, thus determining if the user has the right to access specific functionalities.

## Logic
- The class has a private property `$name` which stores the name associated with the authorization.
- The constructor accepts a string `$name` to initialize the `$name` property.
- The `authorize` method takes a string `$user` as a parameter and returns a boolean:
  - Returns `true` if the provided `$user` matches the initialized `$name`.
  - Returns `false` otherwise.r the user is authorized (true) or not (false).