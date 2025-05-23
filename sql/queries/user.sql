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

-- name: UpdateUserEmailAndPassword :one
update users
set email = $1, hashed_password = $2, updated_at = now()
where id = $3
returning *;

-- name: UpdateUserIsChirpyRed :one
update users
set is_chirpy_red = $1, updated_at = now()
where id = $2
returning *;