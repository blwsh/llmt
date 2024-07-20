# Authorizer.php Documentation

## Overview
The `Authorizer` class is part of the `BenWatson\MyPhpProject\auth` namespace. It is designed to manage user authorization by determining if a provided username matches a predefined name set during the instantiation of the class.

## Logic
- The class is initialized with a `name` which represents the authorized user.
- The `authorize` method accepts a `user` string and checks if it matches the `name`.
- The method returns `true` if the user matches the name; otherwise, it returns `false`.