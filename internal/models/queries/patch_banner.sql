UPDATE banners SET content = $1, is_active = $2 WHERE id = $3;
DELETE FROM banner_tag_feature WHERE banner_id = $1;
INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES