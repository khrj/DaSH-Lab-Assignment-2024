-- Q1
SELECT DISTINCT (category)
FROM crew
ORDER BY category 
LIMIT 10;

-- Q2
SELECT premiered, CONCAT(primary_title, ' (', original_title, ')') AS title FROM titles 
WHERE 
	premiered IS NOT NULL 
	AND genres LIKE "%Action%"
	AND type = "movie"
	AND original_title != primary_title 
ORDER BY premiered DESC, primary_title 
LIMIT 10;

-- Q3
SELECT primary_title, (IFNULL(ended, 2024) - premiered) AS years_running FROM titles 
WHERE 
	premiered IS NOT NULL
	AND type = "tvSeries"
ORDER BY years_running DESC, primary_title
LIMIT 20;

-- Q4
SELECT 
	CONCAT(CAST(10 * ROUND(born / 10) as INTEGER), 's') AS decade, 
	COUNT(DISTINCT people.person_id) as num_directors
FROM crew
JOIN people ON crew.person_id = people.person_id 
WHERE 
	category = "director" 
	AND born IS NOT NULL 
	AND (10 * ROUND(born / 10)) >= 1900
GROUP BY decade
ORDER BY decade;

-- Q5
SELECT 
	titles.type AS title_type, 
	ROUND(AVG(rating),2) AS avg_rating,
	MIN(rating) AS min_rating,
	MAX(rating) AS max_rating
FROM akas 
JOIN titles ON akas.title_id = titles.title_id
JOIN ratings ON akas.title_id = ratings.title_id 
WHERE "language" = "de" AND (types = "imdbDisplay" OR types = "original")
GROUP BY titles.type;

-- Q6
SELECT people.person_id, people.name, ROUND(AVG(rating), 2) AS avg_rating FROM crew
JOIN people ON people.person_id = crew.person_id
JOIN ratings ON crew.title_id = ratings.title_id
WHERE people.person_id IN (
	SELECT people.person_id FROM crew
	JOIN people ON people.person_id = crew.person_id
	WHERE crew.category = "actor" AND crew.characters LIKE "%Batman%"
)
GROUP BY people.person_id
ORDER BY avg_rating DESC
LIMIT 10;

-- Q7
SELECT COUNT(DISTINCT person_id) from people WHERE born = (
	SELECT premiered FROM titles 
	WHERE primary_title = "The Prestige" AND type = "movie"
);

-- Q8
-- Error in the question: 'Rose%' includes names such as Rosemary, the actual query should be for 'Rose %'
-- However, to match the sanity check in the question, 'Rose%' is needed.
SELECT DISTINCT name FROM crew
JOIN people ON people.person_id = crew.person_id
WHERE title_id IN (
	SELECT DISTINCT title_id FROM crew 
	JOIN people ON people.person_id = crew.person_id
	WHERE name LIKE "Rose%"
)
AND category = "director"
ORDER BY name;


-- Q9
SELECT category, name, died, primary_title, runtime_minutes, dead_num AS category_rank FROM (
	SELECT 
		ROW_NUMBER() OVER (
			PARTITION BY category
			ORDER BY died
		) as dead_num,
		*
	FROM (
		SELECT 
			ROW_NUMBER() OVER (
				PARTITION BY person_id
				ORDER BY runtime_minutes DESC, titles.title_id
			) AS work_num,
			*
		FROM crew
		JOIN titles ON titles.title_id = crew.title_id
		WHERE runtime_minutes IS NOT NULL
	) AS long_works
	JOIN people ON people.person_id = long_works.person_id 
	WHERE work_num = 1 AND died IS NOT NULL
)
WHERE dead_num <= 5
ORDER BY category, died, name;

-- Q10
SELECT GROUP_CONCAT(atom) FROM (
	SELECT DISTINCT atom FROM crew, json_each(crew.characters)
	JOIN people ON crew.person_id = people.person_id
	WHERE name = "Leonardo DiCaprio" AND atom NOT LIKE "%Self%"
);
