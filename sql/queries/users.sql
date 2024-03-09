-- name: CreateUser :one
INSERT INTO users (id, email, hashed_password, username, first_name, last_name, user_role)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateUser :one
UPDATE users SET
    username = $2,
    email = $3,
    first_name = $4,
    last_name = $5,
    phone_number = $6,
    date_of_birth = $7,
    gender = $8
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
UPDATE users SET
    account_status = 'deleted'
WHERE id = $1;

-- name: FindUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: PartialFindUsersByUsername :many
SELECT * FROM users
WHERE username LIKE $1
LIMIT $2 OFFSET $3;

-- name: FindUserByUsername :one
SELECT * FROM users
WHERE username = $1;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: FindUserByPassword :one
SELECT * FROM users
WHERE hashed_password = $1;

-- name: GetActiveUsers :many
SELECT * FROM users
WHERE account_status = 'active'
LIMIT $1 OFFSET $2;

-- name: GetInactiveUsers :many
SELECT * FROM users
WHERE account_status = 'inactive'
LIMIT $1 OFFSET $2;

-- name: GetSuspendedUsers :many
SELECT * FROM users
WHERE account_status = 'suspended'
LIMIT $1 OFFSET $2;

-- name: GetDeletedUsers :many
SELECT * FROM users
WHERE account_status = 'deleted'
LIMIT $1 OFFSET $2;

-- name: GetAdminUsers :many
SELECT * FROM users
WHERE user_role = 'admin'
LIMIT $1 OFFSET $2;

-- name: GetSuperAdminUsers :many
SELECT * FROM users
WHERE user_role = 'superadmin'
LIMIT $1 OFFSET $2;

-- name: GetCustomers :many
SELECT * FROM users
WHERE user_role = 'customer'
LIMIT $1 OFFSET $2;

-- name: GetAllUsers :many
SELECT * FROM users
LIMIT $1 OFFSET $2;

-- name: CountAllUsers :one
SELECT COUNT(*) FROM users;

-- name: CountActiveUsers :one
SELECT COUNT(*) FROM users
WHERE account_status = 'active';

-- name: CountInactiveUsers :one
SELECT COUNT(*) FROM users
WHERE account_status = 'inactive';

-- name: CountSuspendedUsers :one
SELECT COUNT(*) FROM users
WHERE account_status = 'suspended';

-- name: CountDeletedUsers :one
SELECT COUNT(*) FROM users
WHERE account_status = 'deleted';

-- name: CountAdminUsers :one
SELECT COUNT(*) FROM users
WHERE user_role = 'admin';

-- name: CountSuperAdminUsers :one
SELECT COUNT(*) FROM users
WHERE user_role = 'superadmin';

-- name: CountCustomers :one
SELECT COUNT(*) FROM users
WHERE user_role = 'customer';

-- name: CountAllUsersByUsername :one
SELECT COUNT(*) FROM users
WHERE username = $1;

-- name: RecoverUser :one
UPDATE users SET
    account_status = 'active'
WHERE id = $1
RETURNING *;

-- name: SuspendUser :one
UPDATE users SET
    account_status = 'suspended'
WHERE id = $1
RETURNING *;

-- name: ActivateUser :one
UPDATE users SET
    account_status = 'active'
WHERE id = $1
RETURNING *;

-- name: DeactivateUser :one
UPDATE users SET
    account_status = 'inactive'
WHERE id = $1
RETURNING *;

-- name: PromoteUserToAdmin :one
UPDATE users SET
    user_role = 'admin'
WHERE id = $1
RETURNING *;

-- name: PromoteUserToSuperAdmin :one
UPDATE users SET
    user_role = 'superadmin'
WHERE id = $1
RETURNING *;

-- name: DemoteSuperAdminToAdmin :one
UPDATE users SET
    user_role = 'admin'
WHERE id = $1
RETURNING *;

-- name: DemoteSuperAdminToUser :one
UPDATE users SET
    user_role = 'user'
WHERE id = $1
RETURNING *;

-- name: DemoteAdminToUser :one
UPDATE users SET
    user_role = 'user'
WHERE id = $1
RETURNING *;

-- name: DemoteUserToCustomer :one
UPDATE users SET
    user_role = 'customer'
WHERE id = $1
RETURNING *;

-- name: UpdateUserPassword :one
UPDATE users SET
    hashed_password = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserLastLogin :exec
UPDATE users SET
    last_login = $2
WHERE id = $1;

-- name: EnableTwoFactorAuth :one
UPDATE users SET
    two_factor_auth = TRUE
WHERE id = $1
RETURNING *;

-- name: DisableTwoFactorAuth :one
UPDATE users SET
    two_factor_auth = FALSE
WHERE id = $1
RETURNING *;

-- name: UpdateUserProfilePicture :one
UPDATE users SET
    profile_picture = $2
WHERE id = $1
RETURNING *;

-- name: StoreRefreshToken :one
INSERT INTO refresh_tokens (id, user_id, token, created_at, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetRefreshTokenByUserID :one
SELECT * FROM refresh_tokens
WHERE user_id = $1;


-- name: GetUserByRefreshToken :one
SELECT * FROM users
WHERE id = (
    SELECT user_id FROM refresh_tokens
    WHERE token = $1
);