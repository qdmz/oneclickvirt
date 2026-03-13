USE oneclickvirt;
SELECT id, `key`, value FROM system_configs WHERE `key`='enable_captcha';
UPDATE system_configs SET value='false' WHERE `key`='enable_captcha';
SELECT id, `key`, value FROM system_configs WHERE `key`='enable_captcha';
