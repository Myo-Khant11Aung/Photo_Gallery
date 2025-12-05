--
-- PostgreSQL database dump (CLEANED + FIXED FOR NEON)
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);

--
-- Data for table: walls
--

INSERT INTO public.walls (id, name, created_at)
VALUES (1, 'Family Wall', '2025-07-07 16:27:17.319187');

--
-- Data for table: albums
--
-- (Your albums table is empty. Leave it empty.)
--

--
-- Data for table: users
-- FIXED: correct column order (id, username, email, password_hash, created_at, wall_id)
--

INSERT INTO public.users (id, username, email, password_hash, created_at, wall_id)
VALUES
(1, 'kevin', 'kevin@example.com', 'hashedpassword', '2025-07-07 16:27:17.312311', NULL),
(2, 'kiki', 'myokhantaung2004@gmail.com', '$2a$10$6NYUG73Z0ntCQT2z6.PmS.ewjus.d9BtbODfhYQjpIeKJxThcEJbC', '2025-08-04 17:59:22.483809', 1);

--
-- Data for table: images
--

INSERT INTO public.images
(id, filename, upload_time, memo, user_id, wall_id, album_date, album_id)
VALUES
(1, 'walls/1/1764797023124090000_1754867493671536000_IMG_4275.PNG', '2025-12-03 13:23:43.480394', '', 2, 1, '2025-12-03', NULL),
(2, 'walls/1/1764797023502874000_1754867493687967000_IMG_4276.PNG', '2025-12-03 13:23:43.684689', 'Hello', 2, 1, '2025-12-03', NULL),
(3, 'walls/1/1764797023687607000_1754867493689576000_IMG_4277.PNG', '2025-12-03 13:23:43.863048', '', 2, 1, '2025-12-03', NULL),
(4, 'walls/1/1764797023865503000_1755052912724381000_IMG_4296.jpg', '2025-12-03 13:23:44.890917', '', 2, 1, '2025-12-03', NULL),
(5, 'walls/1/1764797024894955000_1755052939265883000_IMG_4611.jpg', '2025-12-03 13:23:45.579264', '', 2, 1, '2025-12-03', NULL),
(6, 'walls/1/1764797025581526000_1755052963667693000_IMG_4265.PNG', '2025-12-03 13:23:45.860473', '', 2, 1, '2025-12-03', NULL),
(7, 'walls/1/1764797025863546000_1755052963676084000_IMG_4266.PNG', '2025-12-03 13:23:46.077494', '', 2, 1, '2025-12-03', NULL),
(8, 'walls/1/1764797026078962000_1755052963676646000_IMG_4267.JPEG', '2025-12-03 13:23:46.286376', '', 2, 1, '2025-12-03', NULL),
(9, 'walls/1/1764797026288246000_1755052963677650000_IMG_4268.JPEG', '2025-12-03 13:23:46.517854', '', 2, 1, '2025-12-03', NULL),
(10, 'walls/1/1764797026519901000_1755052963678761000_IMG_4269.PNG', '2025-12-03 13:23:46.762881', '', 2, 1, '2025-12-03', NULL),
(11, 'walls/1/1764797026764662000_1755052963680232000_IMG_4270.JPEG', '2025-12-03 13:23:46.961128', '', 2, 1, '2025-12-03', NULL),
(12, 'walls/1/1764797026962947000_1755052963680944000_IMG_4271.PNG', '2025-12-03 13:23:47.184744', '', 2, 1, '2025-12-03', NULL),
(13, 'walls/1/1764797027186632000_1755052963681386000_IMG_4272.JPEG', '2025-12-03 13:23:47.41719', '', 2, 1, '2025-12-03', NULL),
(14, 'walls/1/1764797027419761000_1755052963682389000_IMG_4273.JPEG', '2025-12-03 13:23:47.617427', '', 2, 1, '2025-12-03', NULL),
(15, 'walls/1/1764797027618745000_1755052963683160000_IMG_4274.PNG', '2025-12-03 13:23:47.842797', '', 2, 1, '2025-12-03', NULL),
(16, 'walls/1/1764797027844295000_1755052963685819000_IMG_4275.PNG', '2025-12-03 13:23:48.0881', '', 2, 1, '2025-12-03', NULL),
(17, 'walls/1/1764797028088915000_1755052963686240000_IMG_4276.PNG', '2025-12-03 13:23:48.300872', '', 2, 1, '2025-12-03', NULL),
(18, 'walls/1/1764797028302271000_LinkedIn_Headshot_Edited.jpg',        '2025-12-03 13:23:48.530975', '', 2, 1, '2025-12-03', NULL);

--
-- Reset sequences
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);
SELECT pg_catalog.setval('public.images_id_seq', 18, true);
SELECT pg_catalog.setval('public.walls_id_seq', 1, true);
SELECT pg_catalog.setval('public.albums_id_seq', 1, false);

