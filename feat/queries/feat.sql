-- name: ListUserFeats :many
SELECT
    feat.*,
    user_feats.user_id,
    user_feats.created_at AS achieved_at
FROM
    feat
    LEFT JOIN user_feats ON user_feats.feat_id = feat.id
WHERE
    user_id = ?;