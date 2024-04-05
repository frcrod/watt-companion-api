-- name: GetAppliances :one
SELECT * FROM appliances
  WHERE id = $1 LIMIT 1;

-- name: GetAppliancesOfUser :many
SELECT * FROM appliances
  WHERE user_id = $1;

-- name: GetGroupsOfUser :many
SELECT * FROM groups
  WHERE user_id = $1;

-- name: InsertApplianceAndReturnId :one
INSERT INTO appliances ("name", wattage, user_id, group_id)
  VALUES ($1, $2, $3, $4)
  RETURNING id; 

-- name: UpdateAppliance :exec
UPDATE appliances 
  SET "name" = $1, wattage = $2
  WHERE id = $3;

-- name: UpdateApplianceGroup :exec
UPDATE appliances 
  SET group_id = $1
  WHERE id = $2;

-- name: UpdateApplianceGroupID :exec
UPDATE appliances 
  SET group_id = $1
  WHERE id = $2;

-- name: CreateUserAndReturnId :one
INSERT INTO users ("email", "nickname")
  VALUES ($1, $2)
  RETURNING id; 

-- name: CheckUserExists :one
SELECT id FROM users
  WHERE email = $1
  LIMIT 1;
