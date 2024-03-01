# Go Auth API Server

The Go Auth API Server is a RESTful API server built with Go that provides authentication and user management functionalities. It uses PostgreSQL as the database backend and integrates with sqlc for database access and Goose for database migrations.

## Table of Contents

- [Go Auth API Server](#go-auth-api-server)
  - [Table of Contents](#table-of-contents)
  - [Setup](#setup)
  - [Protected Routes](#protected-routes)
  - [Endpoints](#endpoints)
    - [User Authentication](#user-authentication)
    - [User Management](#user-management)
    - [Password Reset](#password-reset)
    - [Administration](#administration)

<a name="setup"></a>

## Setup

To set up the Go Auth API Server, follow these steps:

1. **Clone the Repository:**

    ```
    git clone https://github.com/StaphoneWizzoh/Go_Auth.git
    cd Go_Auth
    ```

2. **Database Setup:**

    - Set up a PostgreSQL database.
    - Create a `.env` file in the root directory of the project.
    - Add the following configuration parameters to the `.env` file:
        ```
        DB_URL=postgresql://username:password@localhost:5432/database_name
        PORT=8080
        ```
        Replace `username`, `password`, and `database_name` with your PostgreSQL credentials and database name.

3. **Database Migration:**

    - Run database migrations using Goose to create necessary tables in your database:
        ```bash
        goose -dir sql/schema postgres "user=your_db_user dbname=your_db_name sslmode=disable" up
        ```

4. **Build and Run:**
    - Build and run the Go Auth API server:
        ```bash
        cd cmd
        go build -o go-auth .
        ./go-auth
        ```
    - The server will start running on the port specified in the `.env` file (default: 8080).

<a name="endpoints"></a>

## Protected Routes

These are the API routes that require permissions to access.

-   **Authentication Routes**
    These routes require a valid access token on the header of the request for a valid response.
    ```http
    Authorization: Bearer {<ACCESS_TOKEN>}
    ```
    The routes in this category include:
    -   `/api/users/update`
    -   `/api/users/update-profile-picture`
    -   `/api/users/reset-password`
    -   `/api/admin/*`
-   **Administration Routes**
    These routes require not only valid access token on the header of the request for a valid response but also the access token should belong to an administrator or a super administrator.

    ```http
    Authorization: Bearer {<ACCESS_TOKEN>}
    ```

    ```json
    "user_role": "admin"
    ```

    ```json
    "user_role": "superadmin"
    ```

    The routes in this category include:

    -   `/api/admin/*`

## Endpoints

<a name="user-authentication"></a>

### User Authentication

-   **Register User**

    -   **URL:** `/api/users/register`
    -   **Method:** `POST`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com",
            "password": "password",
            "first_name": "John",
            "last_name": "Doe",
            "user_role": "user"
        }
        ```
    -   The _email_ field should be unique for each account entry.
    -   Accepted values for _user_role_ are: `superadmin`, `admin`, `user`

    -   **Expected Response:**

        ```json
        {
            "id": 1,
            "username": "JohnDoe",
            "email": "user@example.com",
            "created_at": "2024-02-29T19:44:14.906388Z",
            "last_login": {
                "time": "0001-01-01T00:00:00Z",
                "valid": false
            },
            "user_role": "user",
            "profile_picture": {
                "String": "",
                "Valid": false
            },
            "two_factor_auth": false
        }
        ```

    -   **Expected Status Code**
        ```bash
        HTTP/1.1 200 OK
        ```

-   **Login User**

    -   **URL:** `/api/users/login`
    -   **Method:** `POST`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com",
            "password": "password"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "access_token": "<JWT_ACCESS_TOKEN>",
            "refresh_token": "<JWT_REFRESH_TOKEN>"
        }
        ```

-   **Refresh Token**
    -   **URL:** `/api/users/refresh`
    -   **Method:** `POST`
    -   **Request Body:**
        ```json
        {
            "refresh_token": "<JWT_REFRESH_TOKEN>"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "access_token": "<NEW_JWT_ACCESS_TOKEN>"
        }
        ```

<a name="user-management"></a>

### User Management

-   **Update User**

    -   **URL:** `/api/users/update`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com",
            "first_name": "John",
            "last_name": "Doe",
            "phone_number": "1234567890",
            "date_of_birth": "1990-01-01",
            "gender": "male"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "id": 1,
            "email": "user@example.com",
            "first_name": "John",
            "last_name": "Doe",
            "phone_number": "1234567890",
            "date_of_birth": "1990-01-01",
            "gender": "male"
        }
        ```

-   **Update Profile Picture**
    -   **URL:** `/api/users/update-profile-picture`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "profile_picture": "base64_encoded_image"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "id": 1,
            "email": "user@example.com",
            "profile_picture": "http://example.com/profile.jpg"
        }
        ```

<a name="password-reset"></a>

### Password Reset

-   **Request Password Reset**

    -   **URL:** `/api/users/reset-password`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "message": "Reset password email sent successfully"
        }
        ```

-   **Reset Password**
    -   **URL:** `/reset-password`
    -   **Method:** `POST`
    -   **Request Body:**
        ```html
        <form>
            <input type="hidden" name="token" value="<RESET_PASSWORD_TOKEN>" />
            <input type="password" name="password" placeholder="New Password" />
            <input
                type="password"
                name="confirm_password"
                placeholder="Confirm Password"
            />
            <button type="submit">Reset Password</button>
        </form>
        ```
    -   **Expected Response:**
        ```json
        {
            "message": "Password reset successfully"
        }
        ```

### Administration

-   **Promote User to an Administrator**

    -   **URL:** `/api/admin/promote-admin`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "username": "JohnDoe",
            "email": "user@example.com",
            "last_login": {
                "time": "2024-03-01 13:11:05.00489",
                "valid": true
            },
            "account_status": "active",
            "user_role": "admin"
        }
        ```

-   **Promote User to a Super Administrator**

    -   **URL:** `/api/admin/promote-super-admin`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "username": "JohnDoe",
            "email": "user@example.com",
            "last_login": {
                "time": "2024-03-01 13:11:05.00489",
                "valid": true
            },
            "account_status": "active",
            "user_role": "superadmin"
        }
        ```

-   **Suspend a User Account**

    -   **URL:** `/api/admin/suspend-user`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "username": "JohnDoe",
            "email": "user@example.com",
            "last_login": {
                "time": "2024-03-01 13:11:05.00489",
                "valid": true
            },
            "account_status": "suspended",
            "user_role": "user"
        }
        ```

-   **Recover a User Account**

    -   **URL:** `/api/admin/recover-user`
    -   **Method:** `PUT`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "username": "JohnDoe",
            "email": "user@example.com",
            "last_login": {
                "time": "2024-03-01 13:11:05.00489",
                "valid": true
            },
            "account_status": "active",
            "user_role": "user"
        }
        ```

-   **Delete a User Account**

    -   **URL:** `/api/admin/delete-user`
    -   **Method:** `DELETE`
    -   **Request Body:**
        ```json
        {
            "email": "user@example.com"
        }
        ```
    -   **Expected Response:**
        ```json
        {
            "message": "Successfully deactivated the user's account"
        }
        ```

-   **Get all Users**

    -   **URL:** `/api/admin/all-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "user",
                "account_status": "active",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Active users**

    -   **URL:** `/api/admin/active-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "user",
                "account_status": "active",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Administrators**

    -   **URL:** `/api/admin/admin-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "admin",
                "account_status": "active",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Super Administrators**

    -   **URL:** `/api/admin/super-admin-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "superadmin",
                "account_status": "active",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Disabled (Deleted) Accounts**

    -   **URL:** `/api/admin/deleted-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "user",
                "account_status": "deleted",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Inactive Accounts**

    -   **URL:** `/api/admin/inactive-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "user",
                "account_status": "inactive",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```

-   **Get all Suspended Accounts**

    -   **URL:** `/api/admin/suspended-users`
    -   **Method:** `GET`
    -   **Request Body:**
        ```json
        {
            "limit": 50,
            "offset": 5
        }
        ```
    -   **Expected Response:**
        ```json
        [
            {
                "id": 1,
                "username": "JohnDoe",
                "email": "user@example.com",
                "first_name": "John",
                "last_name": "Doe",
                "gender": "male",
                "date_of_birth": "2004-02-29",
                "phone_number": "1234567890",
                "created_at": "2024-02-29T19:44:14.906388Z",
                "last_login": {
                    "time": "2024-03-01 13:11:05.00489",
                    "valid": true
                },
                "user_role": "user",
                "account_status": "suspended",
                "profile_picture": {
                    "String": "",
                    "Valid": false
                },
                "two_factor_auth": false
            }
        ]
        ```
