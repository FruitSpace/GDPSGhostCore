CREATE TABLE users
(
    uid            int(11)      NOT NULL AUTO_INCREMENT PRIMARY KEY,
    uname          varchar(16)  NOT NULL,
    passhash       varchar(128) NOT NULL,
    email          varchar(256) NOT NULL,
    role_id        int(4)       NOT NULL DEFAULT 0,

    stars          int(11)      NOT NULL DEFAULT 0,
    diamonds       int(11)      NOT NULL DEFAULT 0,
    coins          int(11)      NOT NULL DEFAULT 0,
    ucoins         int(11)      NOT NULL DEFAULT 0,
    demons         int(11)      NOT NULL DEFAULT 0,
    cpoints        int(11)      NOT NULL DEFAULT 0,
    orbs           int(11)      NOT NULL DEFAULT 0,
    moons           int(11)      NOT NULL DEFAULT 0,

    regDate        DATETIME     NOT NULL,
    accessDate     DATETIME     NOT NULL,
    lastIP         varchar(64)           DEFAULT 'Unknown',
    gameVer        int(4)                DEFAULT 20,
    lvlsCompleted  int(11)               DEFAULT 0,
    special        int(11)      NOT NULL DEFAULT 0,
    protect_meta JSON NOT NULL DEFAULT '{"comm_time":0,"post_time":0,"msg_time":0}',
    protect_levelsToday int(10) NOT NULL DEFAULT 0,
    protect_todayStars int(10) NOT NULL DEFAULT 0,

    isBanned       tinyint(1)   NOT NULL DEFAULT 0,
    blacklist      text         NOT NULL DEFAULT '',
    friends_cnt    int(11)      NOT NULL DEFAULT 0,
    friendship_ids TEXT         NOT NULL DEFAULT '',

    iconType       TINYINT      NOT NULL DEFAULT 0,
    vessels        JSON         NOT NULL DEFAULT '{"clr_primary":0,"clr_secondary":0,"cube":0,"ship":0,"ball":0,"ufo":0,"wave":0,"robot":0,"spider":0,"swing":0,"jetpack":0,"trace":0,"death":0}',
    chests         JSON         NOT NULL DEFAULT '{"small_count":0,"big_count":0,"small_time":0,"big_time":0}',
    settings       JSON         NOT NULL DEFAULT '{"frS":0,"cS":0,"mS":0,"youtube":"","twitch":"","twitter":""}'
);

CREATE TABLE levels
(
    id                   int(11)          NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name                 varchar(32)      NOT NULL DEFAULT 'Unnamed',
    description          varchar(256)     NOT NULL DEFAULT '',
    uid                  int(11)          NOT NULL,
    password             varchar(8)       NOT NULL,
    version              tinyint          NOT NULL DEFAULT 1,

    length               tinyint(1)       NOT NULL DEFAULT 0,
    difficulty           tinyint(2)       NOT NULL DEFAULT 0,
    demonDifficulty      tinyint(2)       NOT NULL DEFAULT -1,
    suggestDifficulty    float(3, 1)      NOT NULL DEFAULT 0,
    suggestDifficultyCnt int(11)          NOT NULL DEFAULT 0,

    track_id             mediumint(7)     NOT NULL DEFAULT 0,
    song_id              mediumint(7)     NOT NULL DEFAULT 0,
    versionGame          tinyint(3)       NOT NULL,
    versionBinary        tinyint(3)       NOT NULL,
    stringExtra          mediumtext       NOT NULL,
    stringSettings       mediumtext       NOT NULL,
    stringLevel          longtext         NOT NULL,
    stringLevelInfo      mediumtext       NOT NULL,
    original_id          int(11)          NOT NULL DEFAULT 0,

    objects              int(11) UNSIGNED NOT NULL,
    starsRequested       tinyint(2)       NOT NULL,
    starsGot             tinyint(2)       NOT NULL DEFAULT 0,
    ucoins               tinyint(1)       NOT NULL,
    coins                tinyint(1)       NOT NULL DEFAULT 0,
    downloads            int(11) UNSIGNED NOT NULL DEFAULT 0,
    likes                int(11)          NOT NULL DEFAULT 0,
    reports              int(11) UNSIGNED NOT NULL DEFAULT 0,
    collab               TEXT             NOT NULL DEFAULT '',

    is2p                 tinyint(1)       NOT NULL DEFAULT 0,
    isVerified           tinyint(1)       NOT NULL DEFAULT 0,
    isFeatured           tinyint(1)       NOT NULL DEFAULT 0,
    isHall               tinyint(1)       NOT NULL DEFAULT 0,
    isEpic               tinyint(1)       NOT NULL DEFAULT 0,
    isUnlisted           tinyint(1)       NOT NULL DEFAULT 0,
    isLDM                tinyint(1)       NOT NULL DEFAULT 0,

    uploadDate           DATETIME         NOT NULL,
    updateDate           DATETIME         NOT NULL
);

CREATE TABLE levelpacks
(
    id             int(11)      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    packType       tinyint(1)   NOT NULL,
    packName       varchar(256) NOT NULL,
    levels         varchar(512) NOT NULL,

    packStars      tinyint(3)   NOT NULL DEFAULT 0,
    packCoins      tinyint(2)   NOT NULL DEFAULT 0,
    packDifficulty tinyint(2)   NOT NULL,
    packColor      varchar(11)  NOT NULL
);

CREATE TABLE roles
(
    id           int(11)     NOT NULL PRIMARY KEY AUTO_INCREMENT,
    roleName     varchar(64) NOT NULL DEFAULT 'Moderator',
    commentColor varchar(11) NOT NULL DEFAULT '0,0,255',
    modLevel     tinyint(1)  NOT NULL DEFAULT 1,
    privs        text        NOT NULL DEFAULT '{"cRate":0,"cFeature":0,"cEpic":0,"cVerCoins":0,"cDaily":0,"cWeekly":0,"cDelete":0,"cLvlAccess":0,"aRateDemon":0,"aRateReq":0,"aRateStars":0,"aReqMod":0,"dashboardMod":0,"dashboardBan":0,"dashboardCreatePack":0}'
);

CREATE TABLE songs
(
    id          int(11)       NOT NULL PRIMARY KEY AUTO_INCREMENT,
    author_id   int(11)       NOT NULL DEFAULT 0,
    name        varchar(128)  NOT NULL DEFAULT 'Unnamed',
    artist varchar(128)  NOT NULL DEFAULT 'Unknown',
    size        float(5,2)    NOT NULL,
    url         varchar(1024) NOT NULL,
    isBanned    tinyint(1)    NOT NULL DEFAULT 0,
    downloads int NOT NULL DEFAULT 0
);

CREATE TABLE friendships
(
    id     int(12)    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid1   int(11)    NOT NULL,
    uid2   int(11)    NOT NULL,
    u1_new tinyint(1) NOT NULL DEFAULT 1,
    u2_new tinyint(1) NOT NULL DEFAULT 1
);

CREATE TABLE friendreqs
(
    id         int(12)      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid_src    int(11)      NOT NULL,
    uid_dest   int(11)      NOT NULL,
    uploadDate DATETIME     NOT NULL,
    comment    varchar(512) NOT NULL DEFAULT '',
    isNew      tinyint(1)   NOT NULL DEFAULT 1
);

CREATE TABLE acccomments
(
    id         int(12)      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid        int(11)      NOT NULL,
    comment    varchar(128) NOT NULL,
    postedTime DATETIME     NOT NULL,
    likes      int(11)      NOT NULL DEFAULT 0,
    isSpam     tinyint(1)   NOT NULL DEFAULT 0
);

CREATE TABLE comments
(
    id         int(12)      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid        int(11)      NOT NULL,
    lvl_id     int(11)      NOT NULL,
    comment    varchar(128) NOT NULL,
    postedTime DATETIME     NOT NULL,
    likes      int(11)      NOT NULL DEFAULT 0,
    isSpam     tinyint(1)   NOT NULL DEFAULT 0,
    percent    tinyint(3)   NOT NULL
);

CREATE TABLE scores
(
    id         int(12)    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid        int(11)    NOT NULL,
    lvl_id     int(11)    NOT NULL,
    postedTime DATETIME   NOT NULL,
    percent    tinyint(3) NOT NULL,
    attempts   int(11)    NOT NULL DEFAULT 0,
    coins      tinyint(1) NOT NULL DEFAULT 0
);

CREATE TABLE messages
(
    id         int(12)       NOT NULL PRIMARY KEY AUTO_INCREMENT,
    uid_src    int(11)       NOT NULL,
    uid_dest   int(11)       NOT NULL,
    subject    varchar(256)  NOT NULL DEFAULT '',
    body       varchar(1024) NOT NULL,
    postedTime DATETIME      NOT NULL,
    isNew      tinyint(1)    NOT NULL DEFAULT 1
);

CREATE TABLE quests
(
    id         int(12)     NOT NULL PRIMARY KEY AUTO_INCREMENT,
    type       tinyint(1)  NOT NULL,
    name       varchar(64) NOT NULL DEFAULT '',
    needed     int(7)      NOT NULL DEFAULT 0,
    reward     int(7)      NOT NULL DEFAULT 0,
    lvl_id     int(11)     NOT NULL DEFAULT 0,
    timeExpire DATETIME    NOT NULL
);

CREATE TABLE actions
(
    id        int(13)    NOT NULL PRIMARY KEY AUTO_INCREMENT,
    date      DATETIME   NOT NULL,
    uid       int(11)    NOT NULL,
    type      tinyint(1) NOT NULL,
    target_id int(11)    NOT NULL,
    isMod     tinyint(1) NOT NULL DEFAULT 0,
    data      JSON       NOT NULL DEFAULT '{}'
);

CREATE TABLE rateQueue
(
    id              int(11)      NOT NULL PRIMARY KEY AUTO_INCREMENT,
    lvl_id          int(11)          NOT NULL,
    name            varchar(32)      NOT NULL DEFAULT 'Unnamed',
    uid             int(11)          NOT NULL,
    mod_uid         int(11)          NOT NULL,
    stars           int(11)          NOT NULL DEFAULT 0,
    isFeatured      tinyint(1)       NOT NULL DEFAULT 0
);