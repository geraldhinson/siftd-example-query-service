WITH vars AS (
  SELECT
    '126f81a3-7482-4884-8b58-6dba1b209d9a'::VARCHAR AS id,
    'GUID-fake-member-GUID'::varchar AS ownerId
)
SELECT *
FROM public."Resources" r, vars
WHERE r."Id" = vars.id
  AND r."OwnerId" = vars.ownerId;