// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: project_members.sql

package sqlc

import (
	"context"
)

const createProjectMember = `-- name: CreateProjectMember :one
INSERT INTO project_members (
    project_memb_id, id_data_pegawai, project_id, project_role
) VALUES (
    $1, $2, $3, $4::project_role_enum
) RETURNING id, project_memb_id, id_data_pegawai, project_id, project_role
`

type CreateProjectMemberParams struct {
	ProjectMembID string      `json:"project_memb_id"`
	IDDataPegawai string      `json:"id_data_pegawai"`
	ProjectID     string      `json:"project_id"`
	ProjectRole   interface{} `json:"project_role"`
}

func (q *Queries) CreateProjectMember(ctx context.Context, arg CreateProjectMemberParams) (ProjectMember, error) {
	row := q.db.QueryRowContext(ctx, createProjectMember,
		arg.ProjectMembID,
		arg.IDDataPegawai,
		arg.ProjectID,
		arg.ProjectRole,
	)
	var i ProjectMember
	err := row.Scan(
		&i.ID,
		&i.ProjectMembID,
		&i.IDDataPegawai,
		&i.ProjectID,
		&i.ProjectRole,
	)
	return i, err
}

const getProjectMemberByID = `-- name: GetProjectMemberByID :one
SELECT id, project_memb_id, id_data_pegawai, project_id, project_role FROM project_members
WHERE project_memb_id = $1
`

func (q *Queries) GetProjectMemberByID(ctx context.Context, projectMembID string) (ProjectMember, error) {
	row := q.db.QueryRowContext(ctx, getProjectMemberByID, projectMembID)
	var i ProjectMember
	err := row.Scan(
		&i.ID,
		&i.ProjectMembID,
		&i.IDDataPegawai,
		&i.ProjectID,
		&i.ProjectRole,
	)
	return i, err
}

const listProjectMembers = `-- name: ListProjectMembers :many
SELECT id, project_memb_id, id_data_pegawai, project_id, project_role FROM project_members
`

func (q *Queries) ListProjectMembers(ctx context.Context) ([]ProjectMember, error) {
	rows, err := q.db.QueryContext(ctx, listProjectMembers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ProjectMember
	for rows.Next() {
		var i ProjectMember
		if err := rows.Scan(
			&i.ID,
			&i.ProjectMembID,
			&i.IDDataPegawai,
			&i.ProjectID,
			&i.ProjectRole,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateProjectMember = `-- name: UpdateProjectMember :one
UPDATE project_members
SET
    project_role = $1::project_role_enum
WHERE project_memb_id = $2
RETURNING id, project_memb_id, id_data_pegawai, project_id, project_role
`

type UpdateProjectMemberParams struct {
	ProjectRole   interface{} `json:"project_role"`
	ProjectMembID string      `json:"project_memb_id"`
}

func (q *Queries) UpdateProjectMember(ctx context.Context, arg UpdateProjectMemberParams) (ProjectMember, error) {
	row := q.db.QueryRowContext(ctx, updateProjectMember, arg.ProjectRole, arg.ProjectMembID)
	var i ProjectMember
	err := row.Scan(
		&i.ID,
		&i.ProjectMembID,
		&i.IDDataPegawai,
		&i.ProjectID,
		&i.ProjectRole,
	)
	return i, err
}
