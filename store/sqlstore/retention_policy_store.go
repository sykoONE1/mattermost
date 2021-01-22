package sqlstore

import (
	"errors"

	"github.com/mattermost/mattermost-server/v5/einterfaces"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/store"
)

type SqlRetentionPolicyStore struct {
	*SqlStore
	metrics einterfaces.MetricsInterface
}

type retentionPolicyRow struct {
	Id           string
	DisplayName  string
	PostDuration int64
	ChannelId    *string
	TeamId       *string
}

func newSqlRetentionPolicyStore(sqlStore *SqlStore, metrics einterfaces.MetricsInterface) store.RetentionPolicyStore {
	s := &SqlRetentionPolicyStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.RetentionPolicy{}, "RetentionPolicies")
		table.SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("DisplayName").SetMaxSize(26)

		tableC := db.AddTableWithName(model.RetentionPolicyChannel{}, "RetentionPoliciesChannels")
		tableC.SetKeys(false, "ChannelId")
		tableC.ColMap("PolicyId").SetMaxSize(26)
		tableC.ColMap("ChannelId").SetMaxSize(26)

		tableP := db.AddTableWithName(model.RetentionPolicyTeam{}, "RetentionPoliciesTeams")
		tableP.SetKeys(false, "TeamId")
		tableP.ColMap("PolicyId").SetMaxSize(26)
		tableP.ColMap("TeamId").SetMaxSize(26)
	}

	return s
}

func (s *SqlRetentionPolicyStore) createIndexesIfNotExists() {
	s.CreateForeignKeyIfNotExists("RetentionPoliciesChannels", "PolicyId", "RetentionPolicies", "Id", true)
	s.CreateForeignKeyIfNotExists("RetentionPoliciesChannels", "ChannelId", "Channels", "Id", true)
	s.CreateForeignKeyIfNotExists("RetentionPoliciesTeams", "PolicyId", "RetentionPolicies", "Id", true)
	s.CreateForeignKeyIfNotExists("RetentionPoliciesTeams", "TeamId", "Teams", "Id", true)
}

// TODO: check whether the raw queries work with MySQL and SQLite (only tested with PostgreSQL)

func (s *SqlRetentionPolicyStore) Save(policy *model.RetentionPolicy) (*model.RetentionPolicy, error) {
	policy.Id = model.NewId()
	builder := s.getQueryBuilder().
		Insert("RetentionPolicies").
		Columns("Id", "DisplayName", "PostDuration").
		Values(policy.Id, policy.DisplayName, policy.PostDuration)
	_, err := builder.RunWith(s.GetMaster()).Exec()
	setEmptyTeamsAndChannels(policy)
	return policy, err
}

func (s *SqlRetentionPolicyStore) Update(policy *model.RetentionPolicy) (*model.RetentionPolicy, error) {
	builder := s.getQueryBuilder().Update("RetentionPolicies")
	if policy.DisplayName != "" {
		builder = builder.Set("DisplayName", policy.DisplayName)
	}
	if policy.PostDuration != 0 {
		builder = builder.Set("PostDuration", policy.PostDuration)
	}
	builder = builder.
		Where("Id = ?", policy.Id).
		Suffix("RETURNING Id, DisplayName, PostDuration")
	sQuery, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	rows := make([]*model.RetentionPolicy, 0, 1)
	_, err = s.GetMaster().Select(&rows, sQuery, args...)
	if err != nil {
		return nil, err
	} else if len(rows) == 0 {
		return nil, errors.New("policy not found")
	}
	policy = rows[0]
	setEmptyTeamsAndChannels(policy)
	return policy, nil
}

func (s *SqlRetentionPolicyStore) Get(id string) (*model.RetentionPolicy, error) {
	const sQuery = `
	SELECT A.Id, A.DisplayName, A.PostDuration, B.ChannelId, NULL AS TeamId
	FROM RetentionPolicies AS A
	LEFT JOIN RetentionPoliciesChannels AS B ON A.Id = B.PolicyId
	WHERE A.Id = :policyId
	UNION
	SELECT A.Id, A.DisplayName, A.PostDuration, NULL AS ChannelId, C.TeamId
	FROM RetentionPolicies AS A
	LEFT JOIN RetentionPoliciesTeams AS C ON A.Id = C.PolicyId
	WHERE A.Id = :policyId`
	var rows []retentionPolicyRow
	_, err := s.GetReplica().Select(&rows, sQuery, map[string]string{"policyId": id})
	if err != nil {
		return nil, err
	} else if len(rows) == 0 {
		return nil, errors.New("policy not found")
	}
	policy := &model.RetentionPolicy{
		Id:           rows[0].Id,
		DisplayName:  rows[0].DisplayName,
		PostDuration: rows[0].PostDuration,
		ChannelIds:   make([]string, 0),
		TeamIds:      make([]string, 0),
	}
	for _, row := range rows {
		if row.ChannelId != nil {
			policy.ChannelIds = append(policy.ChannelIds, *row.ChannelId)
		} else if row.TeamId != nil {
			policy.TeamIds = append(policy.TeamIds, *row.TeamId)
		}
	}
	return policy, nil
}

func (s *SqlRetentionPolicyStore) GetAll() ([]*model.RetentionPolicy, error) {
	const sQuery = `
	SELECT A.Id, A.DisplayName, A.PostDuration, B.ChannelId, NULL AS TeamId
	FROM RetentionPolicies AS A
	LEFT JOIN RetentionPoliciesChannels AS B ON A.Id = B.PolicyId
	UNION
	SELECT A.Id, A.DisplayName, A.PostDuration, NULL AS ChannelId, C.TeamId
	FROM RetentionPolicies AS A
	LEFT JOIN RetentionPoliciesTeams AS C ON A.Id = C.PolicyId`
	var rows []*retentionPolicyRow
	_, err := s.GetReplica().Select(&rows, sQuery)
	if err != nil {
		return nil, err
	}
	mPolicies := make(map[string]*model.RetentionPolicy)
	for _, row := range rows {
		policy, ok := mPolicies[row.Id]
		if !ok {
			policy = &model.RetentionPolicy{
				Id:           rows[0].Id,
				DisplayName:  rows[0].DisplayName,
				PostDuration: rows[0].PostDuration,
				ChannelIds:   make([]string, 0),
				TeamIds:      make([]string, 0),
			}
			mPolicies[row.Id] = policy
		}
		if row.ChannelId != nil {
			policy.ChannelIds = append(policy.ChannelIds, *row.ChannelId)
		} else if row.TeamId != nil {
			policy.TeamIds = append(policy.TeamIds, *row.TeamId)
		}
	}
	aPolicies := make([]*model.RetentionPolicy, len(mPolicies))
	i := 0
	for _, policy := range mPolicies {
		aPolicies[i] = policy
		i++
	}
	return aPolicies, nil
}

func (s *SqlRetentionPolicyStore) GetAllWithCounts() ([]*model.RetentionPolicyWithCounts, error) {
	const sQuery = `
	SELECT Id, DisplayName, PostDuration,
		SUM(CASE WHEN ChannelId IS NOT NULL THEN 1 ELSE 0 END) AS ChannelCount,
		SUM(CASE WHEN TeamId IS NOT NULL THEN 1 ELSE 0 END) AS TeamCount
	FROM (
		SELECT A.Id, A.DisplayName, A.PostDuration, B.ChannelId, NULL AS TeamId
		FROM RetentionPolicies AS A
		LEFT JOIN RetentionPoliciesChannels AS B ON A.Id = B.PolicyId
		UNION
		SELECT A.Id, A.DisplayName, A.PostDuration, NULL AS ChannelId, C.TeamId
		FROM RetentionPolicies AS A
		LEFT JOIN RetentionPoliciesTeams AS C ON A.Id = C.PolicyId
	) AS D
	GROUP BY Id, DisplayName, PostDuration`
	var rows []*model.RetentionPolicyWithCounts
	_, err := s.GetReplica().Select(&rows, sQuery)
	return rows, err
}

func (s *SqlRetentionPolicyStore) Delete(id string) error {
	builder := s.getQueryBuilder().
		Delete("RetentionPolicies").
		Where("Id = ?", id)
	result, err := builder.RunWith(s.GetMaster()).Exec()
	if err != nil {
		return err
	}
	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if numRowsAffected == 0 {
		return errors.New("policy not found")
	}
	return nil
}

func (s *SqlRetentionPolicyStore) AddChannels(policyId string, channelIds []string) error {
	builder := s.getQueryBuilder().
		Insert("RetentionPoliciesChannels").
		Columns("policyId", "channelId")
	for _, channelId := range channelIds {
		builder = builder.Values(policyId, channelId)
	}
	_, err := builder.RunWith(s.GetMaster()).Exec()
	return err
}

func (s *SqlRetentionPolicyStore) RemoveChannel(policyId string, channelId string) error {
	builder := s.getQueryBuilder().
		Delete("RetentionPoliciesChannels").
		Where("PolicyId = ? AND ChannelId = ?", policyId, channelId)
	result, err := builder.RunWith(s.GetMaster()).Exec()
	if err != nil {
		return err
	}
	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if numRowsAffected == 0 {
		return errors.New("channel was not part of policy")
	}
	return nil
}

func (s *SqlRetentionPolicyStore) AddTeams(policyId string, teamIds []string) error {
	builder := s.getQueryBuilder().
		Insert("RetentionPoliciesTeams").
		Columns("PolicyId", "TeamId")
	for _, teamId := range teamIds {
		builder = builder.Values(policyId, teamId)
	}
	_, err := builder.RunWith(s.GetMaster()).Exec()
	return err
}

func (s *SqlRetentionPolicyStore) RemoveTeam(policyId, teamId string) error {
	builder := s.getQueryBuilder().
		Delete("RetentionPoliciesTeams").
		Where("PolicyId = ? AND TeamId = ?", policyId, teamId)
	result, err := builder.RunWith(s.GetMaster()).Exec()
	numRowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	} else if numRowsAffected == 0 {
		return errors.New("team was not part of policy")
	}
	return nil
}

func setEmptyTeamsAndChannels(policy *model.RetentionPolicy) {
	policy.ChannelIds = make([]string, 0)
	policy.TeamIds = make([]string, 0)
}
