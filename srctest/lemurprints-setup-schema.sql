--
-- This is a test repository stricture for EPrints based on
-- the EPrint repositories at Caltech Library.
-- 
-- ------------------------------------------------------
-- Server version	8.0.18

--
-- Table structure for table `document`
--

DROP TABLE IF EXISTS `document`;
CREATE TABLE `document` (
  `docid` int(11) NOT NULL DEFAULT '0',
  `rev_number` int(11) DEFAULT NULL,
  `eprintid` int(11) DEFAULT NULL,
  `pos` int(11) DEFAULT NULL,
  `format` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `formatdesc` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `language` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `security` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `license` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `main` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `date_embargo_year` smallint(6) DEFAULT NULL,
  `date_embargo_month` smallint(6) DEFAULT NULL,
  `date_embargo_day` smallint(6) DEFAULT NULL,
  `content` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `placement` int(11) DEFAULT NULL,
  `mime_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_duration` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_audio_codec` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_video_codec` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_width` int(11) DEFAULT NULL,
  `media_height` int(11) DEFAULT NULL,
  `media_aspect_ratio` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_sample_start` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `media_sample_stop` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`docid`),
  KEY `document_rev_number` (`rev_number`),
  KEY `document_eprintid` (`eprintid`),
  KEY `document_pos` (`pos`),
  KEY `document_format` (`format`),
  KEY `document_language` (`language`),
  KEY `document_security` (`security`),
  KEY `document_license` (`license`),
  KEY `document_main` (`main`),
  KEY `document_date_eate_embargo_day` (`date_embargo_year`,`date_embargo_month`,`date_embargo_day`),
  KEY `document_content` (`content`),
  KEY `document_placement_1` (`placement`),
  KEY `document_mime_type_1` (`mime_type`)
);

--
-- Table structure for table `document_permission_group`
--

DROP TABLE IF EXISTS `document_permission_group`;
CREATE TABLE `document_permission_group` (
  `docid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `permission_group` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`docid`,`pos`),
  KEY `document_permission_group_pos` (`pos`),
  KEY `document_permisermission_group` (`permission_group`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Table structure for table `document_relation_type`
--

DROP TABLE IF EXISTS `document_relation_type`;
CREATE TABLE `document_relation_type` (
  `docid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `relation_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`docid`,`pos`),
  KEY `document_relation_type_pos` (`pos`),
  KEY `document_relatie_relation_type` (`relation_type`)
);

--
-- Table structure for table `document_relation_uri`
--

DROP TABLE IF EXISTS `document_relation_uri`;
CREATE TABLE `document_relation_uri` (
  `docid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `relation_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`docid`,`pos`),
  KEY `document_relation_uri_pos` (`pos`),
  KEY `document_relation_uri_relation_uri_1` (`relation_uri`)
);

--
-- Table structure for table `eprint`
--

DROP TABLE IF EXISTS `eprint`;
CREATE TABLE `eprint` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `rev_number` int(11) DEFAULT NULL,
  `eprint_status` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `importid` int(11) DEFAULT NULL,
  `source` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `dir` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `datestamp_year` smallint(6) DEFAULT NULL,
  `datestamp_month` smallint(6) DEFAULT NULL,
  `datestamp_day` smallint(6) DEFAULT NULL,
  `datestamp_hour` smallint(6) DEFAULT NULL,
  `datestamp_minute` smallint(6) DEFAULT NULL,
  `datestamp_second` smallint(6) DEFAULT NULL,
  `lastmod_year` smallint(6) DEFAULT NULL,
  `lastmod_month` smallint(6) DEFAULT NULL,
  `lastmod_day` smallint(6) DEFAULT NULL,
  `lastmod_hour` smallint(6) DEFAULT NULL,
  `lastmod_minute` smallint(6) DEFAULT NULL,
  `lastmod_second` smallint(6) DEFAULT NULL,
  `status_changed_year` smallint(6) DEFAULT NULL,
  `status_changed_month` smallint(6) DEFAULT NULL,
  `status_changed_day` smallint(6) DEFAULT NULL,
  `status_changed_hour` smallint(6) DEFAULT NULL,
  `status_changed_minute` smallint(6) DEFAULT NULL,
  `status_changed_second` smallint(6) DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `succeeds` int(11) DEFAULT NULL,
  `commentary` int(11) DEFAULT NULL,
  `replacedby` int(11) DEFAULT NULL,
  `metadata_visibility` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `contact_email` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `fileinfo` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `latitude` float DEFAULT NULL,
  `longitude` float DEFAULT NULL,
  `item_issues_count` int(11) DEFAULT NULL,
  `sword_depositor` int(11) DEFAULT NULL,
  `sword_slug` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `title` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `ispublished` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `full_text_status` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `monograph_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pres_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `keywords` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `note` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `suggestions` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `abstract` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `date_year` smallint(6) DEFAULT NULL,
  `date_month` smallint(6) DEFAULT NULL,
  `date_day` smallint(6) DEFAULT NULL,
  `date_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `series` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `publication` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `volume` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `number` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `publisher` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `place_of_pub` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `edition` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pagerange` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pages` int(11) DEFAULT NULL,
  `event_title` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `event_location` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `event_dates` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `event_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `id_number` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `patent_applicant` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `institution` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `department` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `refereed` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `isbn` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `issn` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `book_title` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `official_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `output_media` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `num_pieces` int(11) DEFAULT NULL,
  `composition_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `data_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pedagogic_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `completion_time` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `task_purpose` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `learning_level` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `rights` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `official_cit` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `doi` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pmc_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `parent_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `alt_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `collection` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `toc` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `interviewer` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `interviewdate` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `errata` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `thesis_degree` VARCHAR(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_degree_grantor` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_submitted_date_year`  smallint(6) DEFAULT NULL,
  `thesis_submitted_date_month` smallint(6) DEFAULT  NULL,
  `thesis_submitted_date_day`   smallint(6) DEFAULT NULL,
  `thesis_defense_date_year` smallint(6) DEFAULT NULL,
  `thesis_defense_date_month` smallint(6) DEFAULT NULL,
  `thesis_defense_date_day` smallint(6) DEFAULT NULL,
  `thesis_degree_date_year` smallint(6) DEFAULT NULL,
  `thesis_degree_date_month` smallint(6) DEFAULT NULL,
  `thesis_degree_date_day` smallint(6) DEFAULT NULL,
  `thesis_approved_date_year` smallint(6) DEFAULT NULL,
  `thesis_approved_date_month` smallint(6) DEFAULT NULL,
  `thesis_approved_date_day` smallint(6) DEFAULT NULL,
  `thesis_public_date_year` smallint(6) DEFAULT NULL,
  `thesis_public_date_month` smallint(6) DEFAULT NULL,
  `thesis_public_date_day` smallint(6) DEFAULT NULL,
  `thesis_awards` LONGTEXT DEFAULT NULL,
  `copyright_statement` LONGTEXT DEFAULT NULL,
  `thesis_author_email` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `gradofc_approval_date_year` smallint(6) DEFAULT NULL,
  `gradofc_approval_date_month` smallint(6) DEFAULT NULL,
  `gradofc_approval_date_day` smallint(6) DEFAULT NULL,
  `hide_thesis_author_email` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `nonsubj_keywords` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `reviewer` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `season` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `classification_code` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `edit_lock_user` int(11) DEFAULT NULL,
  `edit_lock_since` int(11) DEFAULT NULL,
  `edit_lock_until` int(11) DEFAULT NULL,
  `patent_number` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `patent_classification` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `pmid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`),
  KEY `eprint_rev_number` (`rev_number`),
  KEY `eprint_eprint_status` (`eprint_status`),
  KEY `eprint_userid` (`userid`),
  KEY `eprint_importid` (`importid`),
  KEY `eprint_datestamatestamp_second` (`datestamp_year`,`datestamp_month`,`datestamp_day`,`datestamp_hour`,`datestamp_minute`,`datestamp_second`),
  KEY `eprint_lastmod__lastmod_second` (`lastmod_year`,`lastmod_month`,`lastmod_day`,`lastmod_hour`,`lastmod_minute`,`lastmod_second`),
  KEY `eprint_status_c_changed_second` (`status_changed_year`,`status_changed_month`,`status_changed_day`,`status_changed_hour`,`status_changed_minute`,`status_changed_second`),
  KEY `eprint_type` (`type`),
  KEY `eprint_succeeds` (`succeeds`),
  KEY `eprint_replacedby` (`replacedby`),
  KEY `eprint_metadata_visibility` (`metadata_visibility`),
  KEY `eprint_latitude` (`latitude`),
  KEY `eprint_longitude` (`longitude`),
  KEY `eprint_item_issues_count` (`item_issues_count`),
  KEY `eprint_sword_depositor` (`sword_depositor`),
  KEY `eprint_ispublished` (`ispublished`),
  KEY `eprint_full_text_status` (`full_text_status`),
  KEY `eprint_monograph_type` (`monograph_type`),
  KEY `eprint_pres_type` (`pres_type`),
  KEY `eprint_date_yea_month_date_day` (`date_year`,`date_month`,`date_day`),
  KEY `eprint_date_type` (`date_type`),
  KEY `eprint_pagerange` (`pagerange`),
  KEY `eprint_event_type` (`event_type`),
  KEY `eprint_thesis_type` (`thesis_type`),
  KEY `eprint_refereed` (`refereed`),
  KEY `eprint_num_pieces` (`num_pieces`),
  KEY `eprint_pedagogic_type` (`pedagogic_type`),
  KEY `eprint_thesis_defense_date_day` (`thesis_defense_date_year`,`thesis_defense_date_month`,`thesis_defense_date_day`),
  KEY `eprint_thesis_ddegree_date_day` (`thesis_degree_date_year`,`thesis_degree_date_month`,`thesis_degree_date_day`),
  KEY `eprint_thesis_aproved_date_day` (`thesis_approved_date_year`,`thesis_approved_date_month`,`thesis_approved_date_day`),
  KEY `eprint_thesis_ppublic_date_day` (`thesis_public_date_year`,`thesis_public_date_month`,`thesis_public_date_day`),
  KEY `eprint_contact_email_1` (`contact_email`)
);

--
-- Table structure for table `eprint_accompaniment`
--

DROP TABLE IF EXISTS `eprint_accompaniment`;
CREATE TABLE `eprint_accompaniment` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `accompaniment` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_accompaniment_pos` (`pos`)
);

--
-- Table structure for table `eprint_alt_title`
--

DROP TABLE IF EXISTS `eprint_alt_title`;
CREATE TABLE `eprint_alt_title` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `alt_title` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_alt_title_pos` (`pos`)
);

--
-- Table structure for table `eprint_conductors_id`
--

DROP TABLE IF EXISTS `eprint_conductors_id`;
CREATE TABLE `eprint_conductors_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conductors_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conductors_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_conductors_name`
--

DROP TABLE IF EXISTS `eprint_conductors_name`;
CREATE TABLE `eprint_conductors_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conductors_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `conductors_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `conductors_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `conductors_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conductors_name_pos` (`pos`),
  KEY `eprint_conductors_name_conductors_name_family_1` (`conductors_name_family`)
);

--
-- Table structure for table `eprint_conductors_uri`
--

DROP TABLE IF EXISTS `eprint_conductors_uri`;
CREATE TABLE `eprint_conductors_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conductors_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conductors_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_conductors_orcid`
--

DROP TABLE IF EXISTS `eprint_conductors_orcid`;
CREATE TABLE `eprint_conductors_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conductors_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conductors_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_conf_creators_id`
--

DROP TABLE IF EXISTS `eprint_conf_creators_id`;
CREATE TABLE `eprint_conf_creators_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conf_creators_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conf_creators_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_conf_creators_orcid`
--

DROP TABLE IF EXISTS `eprint_conf_creators_orcid`;
CREATE TABLE `eprint_conf_creators_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conf_creators_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conf_creators_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_conf_creators_name`
--

DROP TABLE IF EXISTS `eprint_conf_creators_name`;
CREATE TABLE `eprint_conf_creators_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conf_creators_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conf_creators_name_pos` (`pos`)
);

--
-- Table structure for table `eprint_conf_creators_uri`
--

DROP TABLE IF EXISTS `eprint_conf_creators_uri`;
CREATE TABLE `eprint_conf_creators_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `conf_creators_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_conf_creators_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_contributors_id`
--

DROP TABLE IF EXISTS `eprint_contributors_id`;
CREATE TABLE `eprint_contributors_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `contributors_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_contributors_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_contributors_orcid`
--

DROP TABLE IF EXISTS `eprint_contributors_orcid`;
CREATE TABLE `eprint_contributors_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `contributors_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_contributors_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_contributors_name`
--

DROP TABLE IF EXISTS `eprint_contributors_name`;
CREATE TABLE `eprint_contributors_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `contributors_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `contributors_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `contributors_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `contributors_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_contributors_name_pos` (`pos`),
  KEY `eprint_contributors_name_contributors_name_family_1` (`contributors_name_family`)
);

--
-- Table structure for table `eprint_contributors_type`
--

DROP TABLE IF EXISTS `eprint_contributors_type`;
CREATE TABLE `eprint_contributors_type` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `contributors_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_contributors_type_pos` (`pos`),
  KEY `eprint_contribuntributors_type` (`contributors_type`)
);

--
-- Table structure for table `eprint_contributors_uri`
--

DROP TABLE IF EXISTS `eprint_contributors_uri`;
CREATE TABLE `eprint_contributors_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `contributors_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_contributors_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_copyright_holders`
--

DROP TABLE IF EXISTS `eprint_copyright_holders`;
CREATE TABLE `eprint_copyright_holders` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `copyright_holders` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_copyright_holders_pos` (`pos`)
);

--
-- Table structure for table `eprint_corp_creators_id`
--

DROP TABLE IF EXISTS `eprint_corp_creators_id`;
CREATE TABLE `eprint_corp_creators_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `corp_creators_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_corp_creators_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_corp_creators_ror`
--

DROP TABLE IF EXISTS `eprint_corp_creators_ror`;
CREATE TABLE `eprint_corp_creators_ror` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `corp_creators_ror` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_corp_creators_ror_pos` (`pos`)
);

--
-- Table structure for table `eprint_corp_creators_name`
--

DROP TABLE IF EXISTS `eprint_corp_creators_name`;
CREATE TABLE `eprint_corp_creators_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `corp_creators_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_corp_creators_name_pos` (`pos`)
);

--
-- Table structure for table `eprint_corp_creators_uri`
--

DROP TABLE IF EXISTS `eprint_corp_creators_uri`;
CREATE TABLE `eprint_corp_creators_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `corp_creators_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_corp_creators_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_creators_id`
--

DROP TABLE IF EXISTS `eprint_creators_id`;
CREATE TABLE `eprint_creators_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `creators_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_creators_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_creators_name`
--

DROP TABLE IF EXISTS `eprint_creators_name`;
CREATE TABLE `eprint_creators_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `creators_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `creators_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `creators_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `creators_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_creators_name_pos` (`pos`),
  KEY `eprint_creators_name_creators_name_family_1` (`creators_name_family`)
);

--
-- Table structure for table `eprint_creators_orcid`
--

DROP TABLE IF EXISTS `eprint_creators_orcid`;
CREATE TABLE `eprint_creators_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `creators_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`)
);

--
-- Table structure for table `eprint_creators_uri`
--

DROP TABLE IF EXISTS `eprint_creators_uri`;
CREATE TABLE `eprint_creators_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `creators_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_creators_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_divisions`
--

DROP TABLE IF EXISTS `eprint_divisions`;
CREATE TABLE `eprint_divisions` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `divisions` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_divisions_pos` (`pos`),
  KEY `eprint_divisions_divisions` (`divisions`)
);

--
-- Table structure for table `eprint_editors_id`
--

DROP TABLE IF EXISTS `eprint_editors_id`;
CREATE TABLE `eprint_editors_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editors_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_editors_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_editors_orcid`
--

DROP TABLE IF EXISTS `eprint_editors_orcid`;
CREATE TABLE `eprint_editors_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editors_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_editors_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_editors_name`
--

DROP TABLE IF EXISTS `eprint_editors_name`;
CREATE TABLE `eprint_editors_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editors_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `editors_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `editors_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `editors_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_editors_name_pos` (`pos`),
  KEY `eprint_editors_name_editors_name_family_1` (`editors_name_family`)
);

--
-- Table structure for table `eprint_editors_orcid`
--

DROP TABLE IF EXISTS `eprint_editors_orcid`;
CREATE TABLE `eprint_editors_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editors_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`)
);

--
-- Table structure for table `eprint_editors_uri`
--

DROP TABLE IF EXISTS `eprint_editors_uri`;
CREATE TABLE `eprint_editors_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editors_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_editors_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_exhibitors_id`
--

DROP TABLE IF EXISTS `eprint_exhibitors_id`;
CREATE TABLE `eprint_exhibitors_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `exhibitors_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_exhibitors_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_exhibitors_orcid`
--

DROP TABLE IF EXISTS `eprint_exhibitors_orcid`;
CREATE TABLE `eprint_exhibitors_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `exhibitors_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_exhibitors_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_exhibitors_name`
--

DROP TABLE IF EXISTS `eprint_exhibitors_name`;
CREATE TABLE `eprint_exhibitors_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `exhibitors_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `exhibitors_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `exhibitors_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `exhibitors_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_exhibitors_name_pos` (`pos`),
  KEY `eprint_exhibitors_name_exhibitors_name_family_1` (`exhibitors_name_family`)
);

--
-- Table structure for table `eprint_exhibitors_uri`
--

DROP TABLE IF EXISTS `eprint_exhibitors_uri`;
CREATE TABLE `eprint_exhibitors_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `exhibitors_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_exhibitors_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_funders_agency`
--

DROP TABLE IF EXISTS `eprint_funders_agency`;
CREATE TABLE `eprint_funders_agency` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `funders_agency` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_funders_agency_pos` (`pos`)
);

--
-- Table structure for table `eprint_funders_ror`
--

DROP TABLE IF EXISTS `eprint_funders_ror`;
CREATE TABLE `eprint_funders_ror` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `funders_ror` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_funders_ror_pos` (`pos`)
);

--
-- Table structure for table `eprint_funders_grant_number`
--

DROP TABLE IF EXISTS `eprint_funders_grant_number`;
CREATE TABLE `eprint_funders_grant_number` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `funders_grant_number` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_funders_rant_number_pos` (`pos`)
);

--
-- Table structure for table `eprint_item_issues_comment`
--

DROP TABLE IF EXISTS `eprint_item_issues_comment`;
CREATE TABLE `eprint_item_issues_comment` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_comment` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issues_comment_pos` (`pos`)
);

--
-- Table structure for table `eprint_item_issues_description`
--

DROP TABLE IF EXISTS `eprint_item_issues_description`;
CREATE TABLE `eprint_item_issues_description` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_description` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issdescription_pos` (`pos`)
);

--
-- Table structure for table `eprint_item_issues_id`
--

DROP TABLE IF EXISTS `eprint_item_issues_id`;
CREATE TABLE `eprint_item_issues_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issues_id_pos` (`pos`),
  KEY `eprint_item_issues_id_item_issues_id_1` (`item_issues_id`)
);

--
-- Table structure for table `eprint_item_issues_reported_by`
--

DROP TABLE IF EXISTS `eprint_item_issues_reported_by`;
CREATE TABLE `eprint_item_issues_reported_by` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_reported_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issreported_by_pos` (`pos`),
  KEY `eprint_item_issues_reported_by` (`item_issues_reported_by`)
);

--
-- Table structure for table `eprint_item_issues_resolved_by`
--

DROP TABLE IF EXISTS `eprint_item_issues_resolved_by`;
CREATE TABLE `eprint_item_issues_resolved_by` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_resolved_by` int(11) DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issresolved_by_pos` (`pos`),
  KEY `eprint_item_issues_resolved_by` (`item_issues_resolved_by`)
);

--
-- Table structure for table `eprint_item_issues_status`
--

DROP TABLE IF EXISTS `eprint_item_issues_status`;
CREATE TABLE `eprint_item_issues_status` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_status` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issues_status_pos` (`pos`),
  KEY `eprint_item_issm_issues_status` (`item_issues_status`)
);

--
-- Table structure for table `eprint_item_issues_timestamp`
--

DROP TABLE IF EXISTS `eprint_item_issues_timestamp`;
CREATE TABLE `eprint_item_issues_timestamp` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_timestamp_year` smallint(6) DEFAULT NULL,
  `item_issues_timestamp_month` smallint(6) DEFAULT NULL,
  `item_issues_timestamp_day` smallint(6) DEFAULT NULL,
  `item_issues_timestamp_hour` smallint(6) DEFAULT NULL,
  `item_issues_timestamp_minute` smallint(6) DEFAULT NULL,
  `item_issues_timestamp_second` smallint(6) DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_isss_timestamp_pos` (`pos`),
  KEY `eprint_item_issimestamp_second` (`item_issues_timestamp_year`,`item_issues_timestamp_month`,`item_issues_timestamp_day`,`item_issues_timestamp_hour`,`item_issues_timestamp_minute`,`item_issues_timestamp_second`)
);

--
-- Table structure for table `eprint_item_issues_type`
--

DROP TABLE IF EXISTS `eprint_item_issues_type`;
CREATE TABLE `eprint_item_issues_type` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `item_issues_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_item_issues_type_pos` (`pos`),
  KEY `eprint_item_issues_type_item_issues_type_1` (`item_issues_type`)
);

--
-- Table structure for table `eprint_keyword`
--

DROP TABLE IF EXISTS `eprint_keyword`;
CREATE TABLE `eprint_keyword` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `keyword` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_keyword_pos` (`pos`)
);

--
-- Table structure for table `eprint_local_group`
--

DROP TABLE IF EXISTS `eprint_local_group`;
CREATE TABLE `eprint_local_group` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `local_group` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_local_group_eprintid` (`eprintid`),
  KEY `eprint_local_group_pos` (`pos`)
);

--
-- Table structure for table `eprint_lyricists_orcid`
--

DROP TABLE IF EXISTS `eprint_lyricists_orcid`;
CREATE TABLE `eprint_lyricists_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `lyricists_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_lyricists_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_lyricists_id`
--

DROP TABLE IF EXISTS `eprint_lyricists_id`;
CREATE TABLE `eprint_lyricists_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `lyricists_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_lyricists_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_lyricists_name`
--

DROP TABLE IF EXISTS `eprint_lyricists_name`;
CREATE TABLE `eprint_lyricists_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `lyricists_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `lyricists_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `lyricists_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `lyricists_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_lyricists_name_pos` (`pos`),
  KEY `eprint_lyricists_name_lyricists_name_family_1` (`lyricists_name_family`)
);

--
-- Table structure for table `eprint_lyricists_uri`
--

DROP TABLE IF EXISTS `eprint_lyricists_uri`;
CREATE TABLE `eprint_lyricists_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `lyricists_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_lyricists_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_other_numbering_system_id`
--

DROP TABLE IF EXISTS `eprint_other_numbering_system_id`;
CREATE TABLE `eprint_other_numbering_system_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `other_numbering_system_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_other_nug_system_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_other_numbering_system_name`
--

DROP TABLE IF EXISTS `eprint_other_numbering_system_name`;
CREATE TABLE `eprint_other_numbering_system_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `other_numbering_system_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_other_nusystem_name_pos` (`pos`)
);

--
-- Table structure for table `eprint_patent_assignee`
--

DROP TABLE IF EXISTS `eprint_patent_assignee`;
CREATE TABLE `eprint_patent_assignee` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `patent_assignee` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`)
);

--
-- Table structure for table `eprint_producers_orcid`
--

DROP TABLE IF EXISTS `eprint_producers_orcid`;
CREATE TABLE `eprint_producers_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `producers_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_producers_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_producers_id`
--

DROP TABLE IF EXISTS `eprint_producers_id`;
CREATE TABLE `eprint_producers_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `producers_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_producers_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_producers_name`
--

DROP TABLE IF EXISTS `eprint_producers_name`;
CREATE TABLE `eprint_producers_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `producers_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `producers_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `producers_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `producers_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_producers_name_pos` (`pos`),
  KEY `eprint_producers_name_producers_name_family_1` (`producers_name_family`)
);

--
-- Table structure for table `eprint_producers_uri`
--

DROP TABLE IF EXISTS `eprint_producers_uri`;
CREATE TABLE `eprint_producers_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `producers_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_producers_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_projects`
--

DROP TABLE IF EXISTS `eprint_projects`;
CREATE TABLE `eprint_projects` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `projects` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_projects_pos` (`pos`)
);

--
-- Table structure for table `eprint_reference`
--

DROP TABLE IF EXISTS `eprint_reference`;
CREATE TABLE `eprint_reference` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `reference` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_reference_pos` (`pos`)
);

--
-- Table structure for table `eprint_referencetext`
--

DROP TABLE IF EXISTS `eprint_referencetext`;
CREATE TABLE `eprint_referencetext` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `referencetext` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_referencetext_pos` (`pos`)
);

--
-- Table structure for table `eprint_related_patents`
--

DROP TABLE IF EXISTS `eprint_related_patents`;
CREATE TABLE `eprint_related_patents` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `related_patents` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`)
);

--
-- Table structure for table `eprint_related_url_description`
--

DROP TABLE IF EXISTS `eprint_related_url_description`;
CREATE TABLE `eprint_related_url_description` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `related_url_description` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`)
);

--
-- Table structure for table `eprint_related_url_type`
--

DROP TABLE IF EXISTS `eprint_related_url_type`;
CREATE TABLE `eprint_related_url_type` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `related_url_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_related_url_type_pos` (`pos`),
  KEY `eprint_related_elated_url_type` (`related_url_type`)
);

--
-- Table structure for table `eprint_related_url_url`
--

DROP TABLE IF EXISTS `eprint_related_url_url`;
CREATE TABLE `eprint_related_url_url` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `related_url_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_related_url_url_pos` (`pos`)
);

--
-- Table structure for table `eprint_relation_type`
--

DROP TABLE IF EXISTS `eprint_relation_type`;
CREATE TABLE `eprint_relation_type` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `relation_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_relation_type_pos` (`pos`),
  KEY `eprint_relatione_relation_type` (`relation_type`)
);

--
-- Table structure for table `eprint_relation_uri`
--

DROP TABLE IF EXISTS `eprint_relation_uri`;
CREATE TABLE `eprint_relation_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `relation_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_relation_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_skill_areas`
--

DROP TABLE IF EXISTS `eprint_skill_areas`;
CREATE TABLE `eprint_skill_areas` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `skill_areas` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_skill_areas_pos` (`pos`)
);

--
-- Table structure for table `eprint_subjects`
--

DROP TABLE IF EXISTS `eprint_subjects`;
CREATE TABLE `eprint_subjects` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `subjects` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_subjects_pos` (`pos`),
  KEY `eprint_subjects_subjects` (`subjects`)
);

--
-- Table structure for table `eprint_thesis_advisor_email`
--

DROP TABLE IF EXISTS `eprint_thesis_advisor_email`;
CREATE TABLE `eprint_thesis_advisor_email` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_advisor_email` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_avisor_email_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_advisor_orcid`
--

DROP TABLE IF EXISTS `eprint_thesis_advisor_orcid`;
CREATE TABLE `eprint_thesis_advisor_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_advisor_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_advisor_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_advisor_id`
--

DROP TABLE IF EXISTS `eprint_thesis_advisor_id`;
CREATE TABLE `eprint_thesis_advisor_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_advisor_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_advisor_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_advisor_name`
--

DROP TABLE IF EXISTS `eprint_thesis_advisor_name`;
CREATE TABLE `eprint_thesis_advisor_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_advisor_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_advisor_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_advisor_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_advisor_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_advisor_name_pos` (`pos`),
  KEY `eprint_thesis_advisor_name_thesis_advisor_name_family_1` (`thesis_advisor_name_family`)
);

--
-- Table structure for table `eprint_thesis_advisor_uri`
--

DROP TABLE IF EXISTS `eprint_thesis_advisor_uri`;
CREATE TABLE `eprint_thesis_advisor_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_advisor_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_advisor_uri_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_committee_email`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_email`;
CREATE TABLE `eprint_thesis_committee_email` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_email` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_cittee_email_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_committee_orcid`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_orcid`;
CREATE TABLE `eprint_thesis_committee_orcid` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_orcid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_committee_orcid_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_committee_id`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_id`;
CREATE TABLE `eprint_thesis_committee_id` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_id` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_committee_id_pos` (`pos`)
);

--
-- Table structure for table `eprint_thesis_committee_name`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_name`;
CREATE TABLE `eprint_thesis_committee_name` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_committee_name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_committee_name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `thesis_committee_name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_cmittee_name_pos` (`pos`),
  KEY `eprint_thesis_committee_name_thesis_committee_name_family_1` (`thesis_committee_name_family`)
);

--
-- Table structure for table `eprint_thesis_committee_role`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_role`;
CREATE TABLE `eprint_thesis_committee_role` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_role` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_cmittee_role_pos` (`pos`),
  KEY `eprint_thesis_c_committee_role` (`thesis_committee_role`)
);

--
-- Table structure for table `eprint_thesis_committee_uri`
--

DROP TABLE IF EXISTS `eprint_thesis_committee_uri`;
CREATE TABLE `eprint_thesis_committee_uri` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `thesis_committee_uri` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_thesis_cmmittee_uri_pos` (`pos`)
);

--
-- Table structure for table `event_queue`
--

DROP TABLE IF EXISTS `event_queue`;
CREATE TABLE `event_queue` (
  `eventqueueid` varchar(45) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `datestamp_year` smallint(6) DEFAULT NULL,
  `datestamp_month` smallint(6) DEFAULT NULL,
  `datestamp_day` smallint(6) DEFAULT NULL,
  `datestamp_hour` smallint(6) DEFAULT NULL,
  `datestamp_minute` smallint(6) DEFAULT NULL,
  `datestamp_second` smallint(6) DEFAULT NULL,
  `hash` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `unique` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `oneshot` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `priority` int(11) DEFAULT NULL,
  `start_time_year` smallint(6) DEFAULT NULL,
  `start_time_month` smallint(6) DEFAULT NULL,
  `start_time_day` smallint(6) DEFAULT NULL,
  `start_time_hour` smallint(6) DEFAULT NULL,
  `start_time_minute` smallint(6) DEFAULT NULL,
  `start_time_second` smallint(6) DEFAULT NULL,
  `end_time_year` smallint(6) DEFAULT NULL,
  `end_time_month` smallint(6) DEFAULT NULL,
  `end_time_day` smallint(6) DEFAULT NULL,
  `end_time_hour` smallint(6) DEFAULT NULL,
  `end_time_minute` smallint(6) DEFAULT NULL,
  `end_time_second` smallint(6) DEFAULT NULL,
  `due_time_year` smallint(6) DEFAULT NULL,
  `due_time_month` smallint(6) DEFAULT NULL,
  `due_time_day` smallint(6) DEFAULT NULL,
  `due_time_hour` smallint(6) DEFAULT NULL,
  `due_time_minute` smallint(6) DEFAULT NULL,
  `due_time_second` smallint(6) DEFAULT NULL,
  `repetition` int(11) DEFAULT NULL,
  `status` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `description` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `pluginid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `action` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `params` blob,
  `cleanup` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eventqueueid`),
  KEY `event_queue_datestamp_year_6` (`datestamp_year`,`datestamp_month`,`datestamp_day`,`datestamp_hour`,`datestamp_minute`,`datestamp_second`),
  KEY `event_queue_hash_1` (`hash`),
  KEY `event_queue_unique_1` (`unique`),
  KEY `event_queue_oneshot_1` (`oneshot`),
  KEY `event_queue_priority_1` (`priority`),
  KEY `event_queue_start_time_year_6` (`start_time_year`,`start_time_month`,`start_time_day`,`start_time_hour`,`start_time_minute`,`start_time_second`),
  KEY `event_queue_end_time_year_6` (`end_time_year`,`end_time_month`,`end_time_day`,`end_time_hour`,`end_time_minute`,`end_time_second`),
  KEY `event_queue_due_time_year_6` (`due_time_year`,`due_time_month`,`due_time_day`,`due_time_hour`,`due_time_minute`,`due_time_second`),
  KEY `event_queue_status_1` (`status`),
  KEY `event_queue_userid_1` (`userid`),
  KEY `event_queue_pluginid_1` (`pluginid`),
  KEY `event_queue_action_1` (`action`),
  KEY `event_queue_cleanup_1` (`cleanup`)
);

--
-- Table structure for table `file`
--

DROP TABLE IF EXISTS `file`;
CREATE TABLE `file` (
  `fileid` int(11) NOT NULL,
  `datasetid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `objectid` int(11) DEFAULT NULL,
  `filename` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `mime_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `hash` varchar(64) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `hash_type` varchar(32) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `filesize` bigint(20) DEFAULT NULL,
  `mtime_year` smallint(6) DEFAULT NULL,
  `mtime_month` smallint(6) DEFAULT NULL,
  `mtime_day` smallint(6) DEFAULT NULL,
  `mtime_hour` smallint(6) DEFAULT NULL,
  `mtime_minute` smallint(6) DEFAULT NULL,
  `mtime_second` smallint(6) DEFAULT NULL,
  PRIMARY KEY (`fileid`),
  KEY `file_datasetid_1` (`datasetid`),
  KEY `file_objectid_1` (`objectid`),
  KEY `file_filename_1` (`filename`),
  KEY `file_hash_1` (`hash`),
  KEY `file_hash_type_1` (`hash_type`),
  KEY `file_mtime_year_6` (`mtime_year`,`mtime_month`,`mtime_day`,`mtime_hour`,`mtime_minute`,`mtime_second`)
);

--
-- Table structure for table `file_copies_pluginid`
--

DROP TABLE IF EXISTS `file_copies_pluginid`;
CREATE TABLE `file_copies_pluginid` (
  `fileid` int(11) NOT NULL,
  `pos` int(11) NOT NULL,
  `copies_pluginid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`fileid`,`pos`),
  KEY `file_copies_pluginid_copies_pluginid_1` (`copies_pluginid`)
);

--
-- Table structure for table `file_copies_sourceid`
--

DROP TABLE IF EXISTS `file_copies_sourceid`;
CREATE TABLE `file_copies_sourceid` (
  `fileid` int(11) NOT NULL,
  `pos` int(11) NOT NULL,
  `copies_sourceid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`fileid`,`pos`),
  KEY `file_copies_sourceid_copies_sourceid_1` (`copies_sourceid`)
);

--
-- Table structure for table `history`
--

DROP TABLE IF EXISTS `history`;
CREATE TABLE `history` (
  `historyid` int(11) NOT NULL DEFAULT '0',
  `userid` int(11) DEFAULT NULL,
  `actor` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `datasetid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `objectid` int(11) DEFAULT NULL,
  `revision` int(11) DEFAULT NULL,
  `timestamp_year` smallint(6) DEFAULT NULL,
  `timestamp_month` smallint(6) DEFAULT NULL,
  `timestamp_day` smallint(6) DEFAULT NULL,
  `timestamp_hour` smallint(6) DEFAULT NULL,
  `timestamp_minute` smallint(6) DEFAULT NULL,
  `timestamp_second` smallint(6) DEFAULT NULL,
  `action` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `details` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`historyid`),
  KEY `history_userid` (`userid`),
  KEY `history_objectid` (`objectid`),
  KEY `history_revision` (`revision`),
  KEY `history_timestaimestamp_second` (`timestamp_year`,`timestamp_month`,`timestamp_day`,`timestamp_hour`,`timestamp_minute`,`timestamp_second`),
  KEY `history_action` (`action`),
  KEY `history_datasetid_1` (`datasetid`)
);

--
-- Table structure for table `import`
--

DROP TABLE IF EXISTS `import`;
CREATE TABLE `import` (
  `importid` int(11) NOT NULL DEFAULT '0',
  `datestamp_year` smallint(6) DEFAULT NULL,
  `datestamp_month` smallint(6) DEFAULT NULL,
  `datestamp_day` smallint(6) DEFAULT NULL,
  `datestamp_hour` smallint(6) DEFAULT NULL,
  `datestamp_minute` smallint(6) DEFAULT NULL,
  `datestamp_second` smallint(6) DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `source_repository` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `url` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `description` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `last_run_year` smallint(6) DEFAULT NULL,
  `last_run_month` smallint(6) DEFAULT NULL,
  `last_run_day` smallint(6) DEFAULT NULL,
  `last_run_hour` smallint(6) DEFAULT NULL,
  `last_run_minute` smallint(6) DEFAULT NULL,
  `last_run_second` smallint(6) DEFAULT NULL,
  `last_success_year` smallint(6) DEFAULT NULL,
  `last_success_month` smallint(6) DEFAULT NULL,
  `last_success_day` smallint(6) DEFAULT NULL,
  `last_success_hour` smallint(6) DEFAULT NULL,
  `last_success_minute` smallint(6) DEFAULT NULL,
  `last_success_second` smallint(6) DEFAULT NULL,
  PRIMARY KEY (`importid`),
  KEY `import_datestamatestamp_second` (`datestamp_year`,`datestamp_month`,`datestamp_day`,`datestamp_hour`,`datestamp_minute`,`datestamp_second`),
  KEY `import_userid` (`userid`),
  KEY `import_last_run_year_6` (`last_run_year`,`last_run_month`,`last_run_day`,`last_run_hour`,`last_run_minute`,`last_run_second`),
  KEY `import_last_success_year_6` (`last_success_year`,`last_success_month`,`last_success_day`,`last_success_hour`,`last_success_minute`,`last_success_second`)
);

--
-- Table structure for table `loginticket`
--

DROP TABLE IF EXISTS `loginticket`;
CREATE TABLE `loginticket` (
  `code` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `userid` int(11) DEFAULT NULL,
  `ip` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `expires` int(11) DEFAULT NULL,
  `securecode` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`code`),
  KEY `loginticket_userid` (`userid`),
  KEY `loginticket_expires` (`expires`),
  KEY `loginticket_securecode_1` (`securecode`),
  KEY `loginticket_ip_1` (`ip`),
  KEY `loginticket_time_1` (`time`)
);

--
-- Table structure for table `message`
--

DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `messageid` int(11) NOT NULL DEFAULT '0',
  `datestamp_year` smallint(6) DEFAULT NULL,
  `datestamp_month` smallint(6) DEFAULT NULL,
  `datestamp_day` smallint(6) DEFAULT NULL,
  `datestamp_hour` smallint(6) DEFAULT NULL,
  `datestamp_minute` smallint(6) DEFAULT NULL,
  `datestamp_second` smallint(6) DEFAULT NULL,
  `userid` int(11) DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `message` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`messageid`),
  KEY `message_datestaatestamp_second` (`datestamp_year`,`datestamp_month`,`datestamp_day`,`datestamp_hour`,`datestamp_minute`,`datestamp_second`),
  KEY `message_userid` (`userid`),
  KEY `message_type` (`type`)
);

--
-- Table structure for table `mf`
--

DROP TABLE IF EXISTS `mf`;
CREATE TABLE `mf` (
  `metafieldid` int(11) NOT NULL,
  `mfdatasetid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `parent` int(11) DEFAULT NULL,
  `name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `provenance` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `required` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `multiple` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `allow_null` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `export_as_xml` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `volatile` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `min_resolution` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `sql_index` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `render_input` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `render_value` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `input_ordered` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `maxlength` int(11) DEFAULT NULL,
  `browse_link` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `top` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `datasetid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `set_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `options` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `render_order` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `hide_honourific` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `hide_lineage` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `family_first` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `input_style` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `input_rows` int(11) DEFAULT NULL,
  `input_cols` int(11) DEFAULT NULL,
  `input_boxes` int(11) DEFAULT NULL,
  `sql_counter` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `default_value` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`),
  KEY `mf_mfdatasetid_1` (`mfdatasetid`),
  KEY `mf_parent_1` (`parent`),
  KEY `mf_type_1` (`type`),
  KEY `mf_provenance_1` (`provenance`),
  KEY `mf_required_1` (`required`),
  KEY `mf_multiple_1` (`multiple`),
  KEY `mf_allow_null_1` (`allow_null`),
  KEY `mf_export_as_xml_1` (`export_as_xml`),
  KEY `mf_volatile_1` (`volatile`),
  KEY `mf_min_resolution_1` (`min_resolution`),
  KEY `mf_sql_index_1` (`sql_index`),
  KEY `mf_input_ordered_1` (`input_ordered`),
  KEY `mf_maxlength_1` (`maxlength`),
  KEY `mf_render_order_1` (`render_order`),
  KEY `mf_hide_honourific_1` (`hide_honourific`),
  KEY `mf_hide_lineage_1` (`hide_lineage`),
  KEY `mf_family_first_1` (`family_first`),
  KEY `mf_input_style_1` (`input_style`),
  KEY `mf_input_rows_1` (`input_rows`),
  KEY `mf_input_cols_1` (`input_cols`),
  KEY `mf_input_boxes_1` (`input_boxes`)
);

--
-- Table structure for table `mf_fields_allow_null`
--

DROP TABLE IF EXISTS `mf_fields_allow_null`;
CREATE TABLE `mf_fields_allow_null` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_allow_null` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_allow_null_pos` (`pos`),
  KEY `mf_fields_allowelds_allow_null` (`fields_allow_null`)
);

--
-- Table structure for table `mf_fields_browse_link`
--

DROP TABLE IF EXISTS `mf_fields_browse_link`;
CREATE TABLE `mf_fields_browse_link` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_browse_link` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_browse_link_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_datasetid`
--

DROP TABLE IF EXISTS `mf_fields_datasetid`;
CREATE TABLE `mf_fields_datasetid` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_datasetid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_datasetid_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_default_value`
--

DROP TABLE IF EXISTS `mf_fields_default_value`;
CREATE TABLE `mf_fields_default_value` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_default_value` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`)
);

--
-- Table structure for table `mf_fields_export_as_xml`
--

DROP TABLE IF EXISTS `mf_fields_export_as_xml`;
CREATE TABLE `mf_fields_export_as_xml` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_export_as_xml` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_export_as_xml_pos` (`pos`),
  KEY `mf_fields_expors_export_as_xml` (`fields_export_as_xml`)
);

--
-- Table structure for table `mf_fields_family_first`
--

DROP TABLE IF EXISTS `mf_fields_family_first`;
CREATE TABLE `mf_fields_family_first` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_family_first` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_family_first_pos` (`pos`),
  KEY `mf_fields_familds_family_first` (`fields_family_first`)
);

--
-- Table structure for table `mf_fields_hide_honourific`
--

DROP TABLE IF EXISTS `mf_fields_hide_honourific`;
CREATE TABLE `mf_fields_hide_honourific` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_hide_honourific` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_hide_honourific_pos` (`pos`),
  KEY `mf_fields_hide_hide_honourific` (`fields_hide_honourific`)
);

--
-- Table structure for table `mf_fields_hide_lineage`
--

DROP TABLE IF EXISTS `mf_fields_hide_lineage`;
CREATE TABLE `mf_fields_hide_lineage` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_hide_lineage` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_hide_lineage_pos` (`pos`),
  KEY `mf_fields_hide_ds_hide_lineage` (`fields_hide_lineage`)
);

--
-- Table structure for table `mf_fields_input_boxes`
--

DROP TABLE IF EXISTS `mf_fields_input_boxes`;
CREATE TABLE `mf_fields_input_boxes` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_input_boxes` int(11) DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_input_boxes_pos` (`pos`),
  KEY `mf_fields_inputlds_input_boxes` (`fields_input_boxes`)
);

--
-- Table structure for table `mf_fields_input_cols`
--

DROP TABLE IF EXISTS `mf_fields_input_cols`;
CREATE TABLE `mf_fields_input_cols` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_input_cols` int(11) DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_input_cols_pos` (`pos`),
  KEY `mf_fields_inputelds_input_cols` (`fields_input_cols`)
);

--
-- Table structure for table `mf_fields_input_ordered`
--

DROP TABLE IF EXISTS `mf_fields_input_ordered`;
CREATE TABLE `mf_fields_input_ordered` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_input_ordered` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_input_ordered_pos` (`pos`),
  KEY `mf_fields_inputs_input_ordered` (`fields_input_ordered`)
);

--
-- Table structure for table `mf_fields_input_rows`
--

DROP TABLE IF EXISTS `mf_fields_input_rows`;
CREATE TABLE `mf_fields_input_rows` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_input_rows` int(11) DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_input_rows_pos` (`pos`),
  KEY `mf_fields_inputelds_input_rows` (`fields_input_rows`)
);

--
-- Table structure for table `mf_fields_input_style`
--

DROP TABLE IF EXISTS `mf_fields_input_style`;
CREATE TABLE `mf_fields_input_style` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_input_style` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_input_style_pos` (`pos`),
  KEY `mf_fields_inputlds_input_style` (`fields_input_style`)
);

--
-- Table structure for table `mf_fields_maxlength`
--

DROP TABLE IF EXISTS `mf_fields_maxlength`;
CREATE TABLE `mf_fields_maxlength` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_maxlength` int(11) DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_maxlength_pos` (`pos`),
  KEY `mf_fields_maxleields_maxlength` (`fields_maxlength`)
);

--
-- Table structure for table `mf_fields_mfremoved`
--

DROP TABLE IF EXISTS `mf_fields_mfremoved`;
CREATE TABLE `mf_fields_mfremoved` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_mfremoved` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_mfremoved_fields_mfremoved_1` (`fields_mfremoved`)
);

--
-- Table structure for table `mf_fields_min_resolution`
--

DROP TABLE IF EXISTS `mf_fields_min_resolution`;
CREATE TABLE `mf_fields_min_resolution` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_min_resolution` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_min_resolution_pos` (`pos`),
  KEY `mf_fields_min_r_min_resolution` (`fields_min_resolution`)
);

--
-- Table structure for table `mf_fields_options`
--

DROP TABLE IF EXISTS `mf_fields_options`;
CREATE TABLE `mf_fields_options` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_options` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_options_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_render_input`
--

DROP TABLE IF EXISTS `mf_fields_render_input`;
CREATE TABLE `mf_fields_render_input` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_render_input` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_render_input_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_render_order`
--

DROP TABLE IF EXISTS `mf_fields_render_order`;
CREATE TABLE `mf_fields_render_order` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_render_order` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_render_order_pos` (`pos`),
  KEY `mf_fields_rendeds_render_order` (`fields_render_order`)
);

--
-- Table structure for table `mf_fields_render_value`
--

DROP TABLE IF EXISTS `mf_fields_render_value`;
CREATE TABLE `mf_fields_render_value` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_render_value` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_render_value_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_required`
--

DROP TABLE IF EXISTS `mf_fields_required`;
CREATE TABLE `mf_fields_required` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_required` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_required_pos` (`pos`),
  KEY `mf_fields_requifields_required` (`fields_required`)
);

--
-- Table structure for table `mf_fields_set_name`
--

DROP TABLE IF EXISTS `mf_fields_set_name`;
CREATE TABLE `mf_fields_set_name` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_set_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_set_name_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_sql_counter`
--

DROP TABLE IF EXISTS `mf_fields_sql_counter`;
CREATE TABLE `mf_fields_sql_counter` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_sql_counter` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`)
);

--
-- Table structure for table `mf_fields_sub_name`
--

DROP TABLE IF EXISTS `mf_fields_sub_name`;
CREATE TABLE `mf_fields_sub_name` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_sub_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_sub_name_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_top`
--

DROP TABLE IF EXISTS `mf_fields_top`;
CREATE TABLE `mf_fields_top` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_top` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_top_pos` (`pos`)
);

--
-- Table structure for table `mf_fields_type`
--

DROP TABLE IF EXISTS `mf_fields_type`;
CREATE TABLE `mf_fields_type` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_type` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_type_pos` (`pos`),
  KEY `mf_fields_type_fields_type` (`fields_type`)
);

--
-- Table structure for table `mf_fields_volatile`
--

DROP TABLE IF EXISTS `mf_fields_volatile`;
CREATE TABLE `mf_fields_volatile` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `fields_volatile` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_fields_volatile_pos` (`pos`),
  KEY `mf_fields_volatfields_volatile` (`fields_volatile`)
);

--
-- Table structure for table `mf_phrase_help_lang`
--

DROP TABLE IF EXISTS `mf_phrase_help_lang`;
CREATE TABLE `mf_phrase_help_lang` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `phrase_help_lang` varchar(16) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_phrase_help_lang_pos` (`pos`),
  KEY `mf_phrase_help_hrase_help_lang` (`phrase_help_lang`)
);

--
-- Table structure for table `mf_phrase_help_text`
--

DROP TABLE IF EXISTS `mf_phrase_help_text`;
CREATE TABLE `mf_phrase_help_text` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `phrase_help_text` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_phrase_help_text_pos` (`pos`)
);

--
-- Table structure for table `mf_phrase_name_lang`
--

DROP TABLE IF EXISTS `mf_phrase_name_lang`;
CREATE TABLE `mf_phrase_name_lang` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `phrase_name_lang` varchar(16) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_phrase_name_lang_pos` (`pos`),
  KEY `mf_phrase_name_hrase_name_lang` (`phrase_name_lang`)
);

--
-- Table structure for table `mf_phrase_name_text`
--

DROP TABLE IF EXISTS `mf_phrase_name_text`;
CREATE TABLE `mf_phrase_name_text` (
  `metafieldid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `phrase_name_text` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`metafieldid`,`pos`),
  KEY `mf_phrase_name_text_pos` (`pos`)
);

--
-- Table structure for table `subject`
--

DROP TABLE IF EXISTS `subject`;
CREATE TABLE `subject` (
  `subjectid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `rev_number` int(11) DEFAULT NULL,
  `depositable` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`subjectid`),
  KEY `subject_rev_number` (`rev_number`),
  KEY `subject_depositable` (`depositable`)
);

--
-- Table structure for table `subject_ancestors`
--

DROP TABLE IF EXISTS `subject_ancestors`;
CREATE TABLE `subject_ancestors` (
  `subjectid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `ancestors` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`subjectid`,`pos`),
  KEY `subject_ancestors_pos` (`pos`),
  KEY `subject_ancestors_ancestors_1` (`ancestors`)
);

--
-- Table structure for table `subject_name_lang`
--

DROP TABLE IF EXISTS `subject_name_lang`;
CREATE TABLE `subject_name_lang` (
  `subjectid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `name_lang` varchar(16) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`subjectid`,`pos`),
  KEY `subject_name_lang_pos` (`pos`),
  KEY `subject_name_lang_name_lang` (`name_lang`)
);

--
-- Table structure for table `subject_name_name`
--

DROP TABLE IF EXISTS `subject_name_name`;
CREATE TABLE `subject_name_name` (
  `subjectid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `name_name` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`subjectid`,`pos`),
  KEY `subject_name_name_pos` (`pos`)
);

--
-- Table structure for table `subject_parents`
--

DROP TABLE IF EXISTS `subject_parents`;
CREATE TABLE `subject_parents` (
  `subjectid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '',
  `pos` int(11) NOT NULL DEFAULT '0',
  `parents` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`subjectid`,`pos`),
  KEY `subject_parents_pos` (`pos`),
  KEY `subject_parents_parents_1` (`parents`)
);

--
-- Table structure for table `upload_progress`
--

DROP TABLE IF EXISTS `upload_progress`;
CREATE TABLE `upload_progress` (
  `progressid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL,
  `expires` int(11) DEFAULT NULL,
  `size` bigint(20) DEFAULT NULL,
  `received` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`progressid`),
  KEY `upload_progress_expires_1` (`expires`),
  KEY `upload_progress_size_1` (`size`),
  KEY `upload_progress_received_1` (`received`)
);

--
-- Table structure for table `user`
--

DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `rev_number` int(11) DEFAULT NULL,
  `username` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `usertype` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `newemail` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `newpassword` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pin` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `pinsettime` int(11) DEFAULT NULL,
  `joined_year` smallint(6) DEFAULT NULL,
  `joined_month` smallint(6) DEFAULT NULL,
  `joined_day` smallint(6) DEFAULT NULL,
  `joined_hour` smallint(6) DEFAULT NULL,
  `joined_minute` smallint(6) DEFAULT NULL,
  `joined_second` smallint(6) DEFAULT NULL,
  `email` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `lang` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `frequency` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `mailempty` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `latitude` float DEFAULT NULL,
  `longitude` float DEFAULT NULL,
  `name_honourific` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `name_given` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `name_family` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `name_lineage` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `dept` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `org` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `address` longtext CHARACTER SET utf8 COLLATE utf8_bin,
  `country` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `hideemail` varchar(5) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `os` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  `review_filters` text CHARACTER SET utf8 COLLATE utf8_bin,
  `preference` blob,
  PRIMARY KEY (`userid`),
  KEY `user_rev_number` (`rev_number`),
  KEY `user_usertype` (`usertype`),
  KEY `user_pinsettime` (`pinsettime`),
  KEY `user_joined_yeae_joined_second` (`joined_year`,`joined_month`,`joined_day`,`joined_hour`,`joined_minute`,`joined_second`),
  KEY `user_lang` (`lang`),
  KEY `user_frequency` (`frequency`),
  KEY `user_mailempty` (`mailempty`),
  KEY `user_latitude` (`latitude`),
  KEY `user_longitude` (`longitude`),
  KEY `user_hideemail` (`hideemail`),
  KEY `user_os` (`os`),
  KEY `user_username_1` (`username`),
  KEY `user_newemail_1` (`newemail`),
  KEY `user_email_1` (`email`),
  KEY `user_name_family_1` (`name_family`)
);

--
-- Table structure for table `user_editperms`
--

DROP TABLE IF EXISTS `user_editperms`;
CREATE TABLE `user_editperms` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `editperms` mediumtext CHARACTER SET utf8 COLLATE utf8_bin,
  PRIMARY KEY (`userid`,`pos`),
  KEY `user_editperms_pos` (`pos`)
);

--
-- Table structure for table `user_items_fields`
--

DROP TABLE IF EXISTS `user_items_fields`;
CREATE TABLE `user_items_fields` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `items_fields` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`userid`,`pos`),
  KEY `user_items_fields_pos` (`pos`),
  KEY `user_items_fields_items_fields` (`items_fields`)
);

--
-- Table structure for table `user_permission_group`
--

DROP TABLE IF EXISTS `user_permission_group`;
CREATE TABLE `user_permission_group` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `permission_group` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`userid`,`pos`),
  KEY `user_permission_group_pos` (`pos`),
  KEY `user_permissionermission_group` (`permission_group`)
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Table structure for table `user_review_fields`
--

DROP TABLE IF EXISTS `user_review_fields`;
CREATE TABLE `user_review_fields` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `review_fields` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`userid`,`pos`),
  KEY `user_review_fields_pos` (`pos`),
  KEY `user_review_fies_review_fields` (`review_fields`)
);

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
CREATE TABLE `user_roles` (
  `userid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `roles` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`userid`,`pos`),
  KEY `user_roles_pos` (`pos`),
  KEY `roles_0` (`roles`)
);

--
-- Table structure for table `version`
--

DROP TABLE IF EXISTS `version`;
CREATE TABLE `version` (
  `version` varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT NULL
) ENGINE=MyISAM DEFAULT CHARSET=latin1;

--
-- Table structure for table `eprint_option_major`
--

DROP TABLE IF EXISTS `eprint_option_major`;
CREATE TABLE `eprint_option_major` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `option_major` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_option_major_pos` (`pos`),
  KEY `eprint_option_mor_option_major` (`option_major`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

--
-- Table structure for table `eprint_option_minor`
--

DROP TABLE IF EXISTS `eprint_option_minor`;
CREATE TABLE `eprint_option_minor` (
  `eprintid` int(11) NOT NULL DEFAULT '0',
  `pos` int(11) NOT NULL DEFAULT '0',
  `option_minor` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin DEFAULT NULL,
  PRIMARY KEY (`eprintid`,`pos`),
  KEY `eprint_option_minor_pos` (`pos`),
  KEY `eprint_option_mor_option_minor` (`option_minor`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8 COLLATE=utf8_bin;
