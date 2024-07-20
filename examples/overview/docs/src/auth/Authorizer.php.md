# Authorizer.php Documentation

## Overview
The `Authorizer` class is part of the `BenWatson\MyPhpProject\auth` namespace. It is designed to manage user authorization by verifying if a given username matches the initialized name of the authorizer.

## Logic
- The class is initialized with a string parameter representing a user's name.
- The `authorize` method takes a string parameter for a user.
- It checks if the provided user name matches the internal name stored in the instance.
- Returns `true` if the user matches; otherwise, it returns `false`.ameter and checks if it matches the stored `$name`.
- The `authorize` method returns a boolean value indicating whether the user is authorized (true) or not (false).