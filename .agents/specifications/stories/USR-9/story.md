# USR-9: Profile Management

## Overview
Allow users to manage their own profile data (such as name, email, and changing their password) while logged into the application.

## Requirements
1. The user profile needs to have a `name` field (or `first_name` and `last_name` fields), as the `User` model currently only has `email` and `password`.
2. A new route (e.g. `/profile`) to view and edit profile data.
3. A form allowing the user to change their `name` and `email`.
4. A form allowing the user to change their password securely (requiring their current password to change it).
