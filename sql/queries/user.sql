-- name: CreateUser :one
insert into users(id, created_at, updated_at, email, hashed_password)
      values ($1,$2,$3,$4, $5)
          RETURNING *;
-- name: RemoveUsers :exec
delete from users;

-- name: GetUserByEmail :one
select * from users where email = $1;

-- name: GetUserById :one
select * from users where id = $1;

