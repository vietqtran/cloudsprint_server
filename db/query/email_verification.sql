-- name: CreateEmailOTP :one
INSERT INTO email_otps (
    id,
    email,
    otp,
    expires_at
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetEmailOTPByCode :one
SELECT * FROM email_otps
WHERE otp = $1 AND email = $2 AND expires_at > NOW() AND used = false
LIMIT 1;

-- name: MarkEmailOTPUsed :exec
UPDATE email_otps
SET used = true
WHERE id = $1;

-- name: UpdateAccountEmailVerified :one
UPDATE accounts
SET email_verified = true,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetEmailVerificationStatus :one
SELECT email_verified 
FROM accounts
WHERE id = $1
LIMIT 1; 