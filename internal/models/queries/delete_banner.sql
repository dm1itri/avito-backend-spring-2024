DELETE FROM banners WHERE id = $1;
DELETE FROM banner_tag_feature WHERE banner_id = $1
