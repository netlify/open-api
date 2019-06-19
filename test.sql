SELECT
    t.slug,
    p.created_at,
    t.title,
    p.cooked,
    u.username,
    p.user_id,
    t.id
FROM topics t, posts p, users u
WHERE p.user_id = u.id
    --AND p.created_at  > NOW() - INTERVAL '10 DAY'
    --AND p.created_at  < NOW() - INTERVAL '1 DAY' 
    --AND u.admin = false
    AND t.last_post_user_id = u.id
    AND t.closed = true
    AND t.visible = true
    AND p.deleted_at IS null
    AND t.archetype <> 'private_message'
    AND u.username <> 'system'
    AND u.username <> 'discourse'
    AND u.username <> 'discobot'
    AND t.id != 5
    AND t.id != -1
    AND p.user_id != -1
    AND
    NOT EXISTS (SELECT topic_custom_fields.*
    FROM topic_custom_fields
    WHERE
    t.id = topic_custom_fields.topic_id
        AND
        ((topic_custom_fields.name = 'accepted_answer_post_id' AND value IS NOT NULL)
        OR
        (topic_custom_fields.name = 'answered_state' AND topic_custom_fields.value ='true')))
GROUP BY p.created_at,p.cooked, p.user_id,t.id, t.visible, u.username, t.user_id, t.slug, t.closed, t.title
ORDER BY p.created_at DESC