INSERT INTO banners (id, content, is_active) VALUES (DEFAULT, $1, $2) RETURNING id;
INSERT INTO banner_tag_feature (banner_id, tag_id, feature_id) VALUES