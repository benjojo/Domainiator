-- --------------------------------------------------------
-- Host:                         localhost
-- Server version:               5.5.31-0+wheezy1-log - (Debian)
-- Server OS:                    debian-linux-gnu
-- HeidiSQL version:             7.0.0.4053
-- Date/time:                    2013-11-30 14:25:53
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET FOREIGN_KEY_CHECKS=0 */;

-- Dumping database structure for Domaniator
DROP DATABASE IF EXISTS `Domaniator`;
CREATE DATABASE IF NOT EXISTS `Domaniator` /*!40100 DEFAULT CHARACTER SET latin1 */;
USE `Domaniator`;


-- Dumping structure for table Domaniator.CachedResults
DROP TABLE IF EXISTS `CachedResults`;
CREATE TABLE IF NOT EXISTS `CachedResults` (
  `Day` int(10) NOT NULL AUTO_INCREMENT,
  `RequestCount` bigint(20) DEFAULT '0',
  `FailedCount` int(11) DEFAULT '0',
  `TopHeaders` text,
  `AvgContentSize` float DEFAULT '0',
  PRIMARY KEY (`Day`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.


-- Dumping structure for table Domaniator.Results
DROP TABLE IF EXISTS `Results`;
CREATE TABLE IF NOT EXISTS `Results` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `Timestamp` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `Domain` varchar(65) DEFAULT NULL,
  `Data` text,
  PRIMARY KEY (`id`),
  KEY `Domain` (`Domain`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- Data exporting was unselected.
/*!40014 SET FOREIGN_KEY_CHECKS=1 */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
