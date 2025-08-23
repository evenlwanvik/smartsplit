package workout

import (
	"context"
	"database/sql"
)

// Repository provides access to workout domain store.
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new Workout repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SelectMuscle(ctx context.Context, id int) (*Muscle, error) {
	const query = `
SELECT id, name, muscle_group, description
FROM workout.muscles
WHERE id = $1;
`
	var muscle Muscle
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&muscle.ID,
		&muscle.Name,
		&muscle.Group,
		&muscle.Description,
	)
	return &muscle, err
}

// SelectMuscles returns slice of muscles.
func (r *Repository) SelectMuscles(ctx context.Context) ([]*Muscle, error) {
	// TODO: Make muscles user specific. Maybe in a later version.
	const query = `
SELECT id, name, muscle_group, description
FROM workout.muscles;
`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var muscles []*Muscle
	for rows.Next() {
		m := new(Muscle)
		if err := rows.Scan(&m.ID, &m.Name, &m.Group, &m.Description); err != nil {
			return nil, err
		}
		muscles = append(muscles, m)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return muscles, nil
}

// InsertMuscle inserts a new muscle and returns its ID.
func (r *Repository) InsertMuscle(ctx context.Context, input *MuscleInput) (*Muscle, error) {
	const query = `
INSERT INTO workout.muscles (name, muscle_group, description)
VALUES ($1, $2, $3)
RETURNING id, name, muscle_group, description;
`
	var muscle Muscle
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.Name,
		input.MuscleGroup,
		input.Description,
	).Scan(
		&muscle.ID,
		&muscle.Name,
		&muscle.Group,
		&muscle.Description,
	)
	return &muscle, err
}

// SelectRanks returns a slice of muscle ranks.
func (r *Repository) SelectRanks(ctx context.Context, filters Filters) ([]*MuscleRank, error) {
	const query = `
SELECT id, user_id, muscle_id, rank, updated_at
FROM workout.muscles_ranks
WHERE (user_id = :user_id OR :user_id IS NULL);
`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("user_id", filters.UserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ranks []*MuscleRank
	for rows.Next() {
		mr := new(MuscleRank)
		if err := rows.Scan(&mr.ID, &mr.UserID, &mr.MuscleID, &mr.Rank, &mr.UpdatedAt); err != nil {
			return nil, err
		}
		ranks = append(ranks, mr)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ranks, nil
}

// UpsertRank creates or updates a muscle rank for a user.
func (r *Repository) UpsertRank(ctx context.Context, input *MuscleRank) (*MuscleRank, error) {
	const query = `
INSERT INTO workout.muscles_ranks (user_id, muscle_id, rank)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, muscle_id)
  DO UPDATE SET rank = EXCLUDED.rank, updated_at = now()
RETURNING id, user_id, muscle_id, rank, updated_at;
`
	var mr MuscleRank
	err := r.db.QueryRowContext(ctx, query, input.UserID, input.MuscleID, input.Rank).
		Scan(&mr.ID, &mr.UserID, &mr.MuscleID, &mr.Rank, &mr.UpdatedAt)
	return &mr, err
}

// SelectPlans returns a slice of workout plans.
func (r *Repository) SelectPlans(ctx context.Context, filters Filters) ([]*Plan, error) {
	const query = `
SELECT id, user_id, date, created_at, notes
FROM workout.plans
WHERE (user_id = :user_id OR :user_id IS NULL);
`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("user_id", filters.UserID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*Plan
	for rows.Next() {
		p := new(Plan)
		if err := rows.Scan(&p.ID, &p.UserID, &p.Date, &p.CreatedAt, &p.Notes); err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return plans, nil
}

// InsertPlan inserts a new workout plan and returns its ID.
func (r *Repository) InsertPlan(ctx context.Context, input PlanInput) (*Plan, error) {
	const query = `
INSERT INTO workout.plans (user_id, date, notes)
VALUES ($1, NOW(), $2)
RETURNING id, user_id, date, notes;
`
	var plan Plan
	err := r.db.QueryRowContext(ctx, query, input.UserID, input.Notes).Scan(
		&plan.ID, &plan.UserID, &plan.Date, &plan.Notes,
	)
	return &plan, err
}

// DeletePlan deletes a workout plan by ID; returns deleted plan.
func (r *Repository) DeletePlan(ctx context.Context, id int) (*Plan, error) {
	const query = `
DELETE FROM workout.plans
WHERE id = $1
RETURNING id, user_id, date, created_at, notes;
`
	var plan Plan
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&plan.ID, &plan.UserID, &plan.Date, &plan.CreatedAt, &plan.Notes,
	)
	return &plan, err
}

// SelectPlanEntries returns a slice of plan entries.
func (r *Repository) SelectPlanEntries(ctx context.Context, filters Filters) ([]*PlanEntry, error) {
	const query = `
SELECT id, plan_id, muscle_id, sets, created_at
FROM workout.plan_entries
WHERE (user_id = $1 OR $1 IS NULL)
AND (plan_id = $2 OR $2 IS NULL);
`
	rows, err := r.db.QueryContext(
		ctx,
		query,
		filters.UserID,
		filters.PlanID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*PlanEntry
	for rows.Next() {
		e := new(PlanEntry)
		if err := rows.Scan(&e.ID, &e.PlanID, &e.MuscleID, &e.Sets, &e.CreatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// DeleteManyPlanEntries deletes plan entries by filters; returns number of deleted entries.
func (r *Repository) DeleteManyPlanEntries(ctx context.Context, filters Filters) (int64, error) {
	const query = `
DELETE FROM workout.plan_entries
WHERE (plan_id = $1 OR $1 IS NULL)
AND (muscle_id = $2 OR $2 IS NULL);
`
	result, err := r.db.ExecContext(
		ctx,
		query,
		filters.PlanID,
		filters.MuscleID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InsertPlanEntry creates a new plan entry; returns error if duplicate.
func (r *Repository) InsertPlanEntry(ctx context.Context, input PlanEntry) (*PlanEntry, error) {
	const query = `
INSERT INTO workout.plan_entries (plan_id, muscle_id, sets)
VALUES ($1, $2, $3)
RETURNING id, created_at, plan_id, muscle_id, sets;
`
	var pe PlanEntry
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.PlanID,
		input.MuscleID,
		input.Sets,
	).Scan(&pe.ID, &pe.CreatedAt, &pe.PlanID, &pe.MuscleID, &pe.Sets)
	return &pe, err
}

func (r *Repository) PatchPlanEntry(ctx context.Context, input PlanEntryPatch) (*PlanEntry, error) {
	const query = `
UPDATE workout.plan_entries
SET sets = $2
WHERE id = $1
RETURNING id, created_at, plan_id, muscle_id, sets;
`
	var pe PlanEntry
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.ID,
		input.Sets,
	).Scan(&pe.ID, &pe.CreatedAt, &pe.PlanID, &pe.MuscleID, &pe.Sets)
	return &pe, err
}
