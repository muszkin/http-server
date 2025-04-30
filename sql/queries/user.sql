-- name: CreateUser :one
insert into users(id, created_at, updated_at, email)
      values ($1,$2,$3,$4)
          RETURNING *;
-- name: RemoveUsers :exec
delete from users;
