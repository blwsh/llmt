# Authorizer.php Documentation

## Overview
The `Authorizer` class is part of the `BenWatson\MyPhpProject\auth` namespace. It is designed to facilitate user authorization by comparing a provided username with a stored name. The class uses a simple equality check to determine if the user is authorized.

## Logic
- The class is instantiated with a `name` parameter, which is stored as a private property.
- It contains an `authorize` method that compares the provided `user` string with the stored `name`.
- The `authorize` method returns `true` if the `user` matches the `name`, and `false` otherwise. matches the `$name`.
  - Returns `false` otherwise. `false` otherwise.r the user is authorized (true) or not (false).