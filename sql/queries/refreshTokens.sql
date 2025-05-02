-- name: CreateRefreshToken :one
insert into refresh_tokens(token, created_at, updated_at, user_id, expires_at, revoked_at)
values($1,now(),now(),$2,$3,null)
returning *;

-- name: GetUserFromRefreshToken :one
select users.* from users
join public.refresh_tokens rt on users.id = rt.user_id
where rt.token = $1 and revoked_at is null and expires_at > now();

-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = now(), updated_at = now()
where token = $1;